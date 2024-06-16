from typing import Any, Callable, Dict, Iterable, List, Optional, Tuple

import numpy as np
import pandas as pd
from attrs import define
from numpy.polynomial.chebyshev import Chebyshev
from sklearn.metrics import r2_score
from sklearn.preprocessing import MinMaxScaler
from scipy.optimize import curve_fit


@define
class FittedTrendModel:
    model: Any
    one_month_step: float
    scaler: MinMaxScaler
    last_date: np.datetime64
    last_value: float
    max_value: float


class Model:
    def __init__(
        self,
        df: pd.DataFrame,
        period_model: Any,
        regular_codes: Optional[Tuple[str]] = None,
        value_column: str = "paid_rub",
        date_column: str = "conclusion_date",
    ) -> None:
        """
        Initialize the Model object.

        Args:
            df (pd.DataFrame): The input data containing the data to be used for prediction.
                The DataFrame should have columns "depth3_code_kpgz", "paid_rub", and "conclusion_date".
            period_model (Any): The period model object used for prediction.
            regular_codes (Optional[Tuple[str]], optional): The regular codes to be used for filtering.
                Defaults to None.
            value_column (str, optional): The column name for the value column. Defaults to "paid_rub".
            date_column (str, optional): The column name for the date column. Defaults to "conclusion_date".
        """
        self._value_column = value_column
        self._date_column = date_column
        self._df = df[["depth3_code_kpgz", value_column, date_column]].copy()

        segments_counts = self._df["depth3_code_kpgz"].value_counts()
        self._uniq_segments = self.filter_regular_codes(segments_counts, regular_codes)

        self._period_model = period_model

        self._trend_models_by_segment, self._periods_df = self.fit()

    @property
    def filtered_segments(self) -> Tuple[str]:
        """
        Get the filtered segments.

        Returns:
            Tuple[str]: The filtered segments.
        """
        return tuple(self._trend_models_by_segment.keys())

    def filter_regular_codes(
        self,
        segments: pd.Series,
        regular_codes: Optional[Tuple[str]] = None,
        count_threshold: int = 3,
        starts_with: Optional[str] = "01.",
    ) -> Tuple[str]:
        """
        Filter the segments based on given criteria.

        Args:
            segments (pd.Series): The segments to be filtered.
            regular_codes (Optional[Tuple[str]], optional): The regular codes to be used for filtering.
                Defaults to None.
            count_threshold (int, optional): The count threshold for filtering. Defaults to 3.
            starts_with (Optional[str], optional): The prefix for filtering. Defaults to "01.".

        Returns:
            Tuple[str]: The filtered segments.
        """
        segments = segments[segments > count_threshold]
        segments_np = segments.index.to_numpy()

        if starts_with is not None:
            segments_mask = segments.index.str.startswith(starts_with)
            segments_np = segments.index.to_numpy()[segments_mask]
        segments = segments_np

        if regular_codes is not None:
            segments = np.intersect1d(segments, regular_codes)

        return segments

    def fit_trend(self, x_data: np.ndarray, y_data: np.ndarray) -> Callable[[float], float]:
        """
        Fit a trend model to the given data.

        Args:
            x_data (np.ndarray): The input data.
            y_data (np.ndarray): The target data.

        Returns:
            Callable[[float], float]: A function that takes a float as input and returns the predicted value.
        """

        def app_func(x: float, *params: float) -> float:
            log_func = params[0] * np.log(params[1] * x + params[2])
            poly_func = params[3] * x + params[4]
            sin_func = params[5] * (np.sin(params[6] * (x + params[7])) + params[6] * x)
            return (
                (np.clip(log_func, 0, None) + poly_func) * params[8]
                + params[9]
                + sin_func
            )

        bounds = (
            [0.01, 0, 0.001, 0, 0, 0.001, 0, 0, 0.0001, 0],
            [
                np.inf,
                np.inf,
                np.inf,
                np.inf,
                np.inf,
                0.1,
                np.inf,
                np.pi / 2,
                np.inf,
                np.inf,
            ],
        )
        try:
            popt_func, pcov_func = curve_fit(
                app_func,
                x_data,
                y_data,
                p0=[0.05, 1, 0.5, 1, 0, 0.001, 1.0, 0, 1, 0],
                maxfev=5000,
                bounds=bounds,
            )

            def output_func(x):
                return app_func(x, *popt_func)

        except Exception as ex:
            print(ex)
            model = Chebyshev.fit(x_data, y_data, 1)

            def output_func(x):
                return model(x)

        return output_func

    def train_trend_models(
        self,
        segments: Iterable[str],
        min_dates_records: int = 3,
        min_r2: float = 0.7,
    ) -> dict[str, FittedTrendModel]:
        """
        Train trend models for each segment.

        Args:
            segments (Iterable[str]): The segments to train models for.
            min_dates_records (int, optional): Minimum number of non-zero records
                required for a segment. Defaults to 3.
            min_r2 (float, optional): Minimum r-squared value for a model.
                Defaults to 0.7.

        Returns:
            dict[str, FittedTrendModel]: A dictionary mapping segment names to fitted
                trend models.
        """
        trend_models_by_segment: dict[str, FittedTrendModel] = {}

        for segment in segments:
            filtered_df = self._df[self._df["depth3_code_kpgz"] == segment]
            filtered_df = (
                filtered_df.resample("ME", on=self._date_column)[[self._value_column]]
                .sum()
                .reset_index()
            )

            if (filtered_df[self._value_column] > 0).sum() < min_dates_records:
                continue

            filtered_df["init_value"] = filtered_df[self._value_column]
            filtered_df[self._value_column] = filtered_df[self._value_column].cumsum()
            filtered_df["date_norm"] = (
                filtered_df[self._date_column].astype(np.int64)
                - filtered_df[self._date_column].astype(np.int64).min()
            ) / (
                filtered_df[self._date_column].astype(np.int64).max()
                - filtered_df[self._date_column].astype(np.int64).min()
            )

            scaler = MinMaxScaler()
            X = filtered_df["date_norm"]
            y = filtered_df[self._value_column]
            y = scaler.fit_transform(y.to_numpy().reshape(-1, 1)).reshape(-1)

            fitted_model = self.fit_trend(X, y)

            if r2_score(y, fitted_model(X)) < min_r2:
                continue

            fitted_model = FittedTrendModel(
                model=fitted_model,
                scaler=scaler,
                one_month_step=filtered_df["date_norm"].iloc[1]
                - filtered_df["date_norm"].iloc[0],
                last_date=filtered_df[self._date_column].iloc[-1],
                last_value=filtered_df[self._value_column].iloc[-1],
                max_value=filtered_df["init_value"].max(),
            )

            trend_models_by_segment[segment] = fitted_model

        return trend_models_by_segment

    def get_periods(self, segments: Iterable[str]) -> pd.DataFrame:
        """
        Get periods for each segment.

        Args:
            segments (Iterable[str]): The segments to get periods for.

        Returns:
            pd.DataFrame: A DataFrame containing the last date, period, and segment.
        """
        all_periods = []
        all_timestamps = []
        all_segments = []

        for segment in segments:
            filtered_df = self._df[self._df["depth3_code_kpgz"] == segment]

            period, timestamp = self._period_model.predict(
                filtered_df,
                dates_column=self._date_column,
                values_column=self._value_column,
            )
            all_periods.append(period)
            all_timestamps.append(timestamp)
            all_segments.append(segment)

        periods_df = pd.DataFrame(
            {
                "last_date": all_timestamps,
                "period": all_periods,
                "depth3_code_kpgz": segments,
            }
        )

        return periods_df

    def fit(self) -> Tuple[Dict[str, FittedTrendModel], pd.DataFrame]:
        """
        Train the trend models and get the periods for each segment.

        Returns:
            Tuple[Dict[str, FittedTrendModel], pd.DataFrame]: A tuple containing the trained trend models
            and the periods DataFrame.
        """
        trend_models_by_segment = self.train_trend_models(self._uniq_segments)
        periods_df = self.get_periods(self._uniq_segments)

        return trend_models_by_segment, periods_df

    def predict(
        self,
        start_date: np.datetime64,
        num_months: int,
        segment: str,
    ) -> List[Dict[np.datetime64, float]]:
        """
        Predict the future values for a given segment.

        Args:
            start_date (np.datetime64): The start date for the prediction.
            num_months (int): The number of months to predict.
            segment (str): The segment to predict for.

        Returns:
            List[Dict[np.datetime64, float]]: A list of dictionaries containing the predicted values for each
            future date. Each dictionary contains a single key-value pair, where the key is the predicted date
            and the value is the predicted value.
        """
        segment_period_info = self._periods_df[
            self._periods_df["depth3_code_kpgz"] == segment
        ]
        period = round(segment_period_info["period"].iloc[0])
        last_date = segment_period_info["last_date"].iloc[0]

        fitted_trend_model = self._trend_models_by_segment[segment]
        scaler = fitted_trend_model.scaler
        last_fitted_date = fitted_trend_model.last_date
        last_value = fitted_trend_model.last_value
        one_month_step = fitted_trend_model.one_month_step
        max_value = fitted_trend_model.max_value

        dates = pd.date_range(
            start=last_fitted_date,
            end=start_date + pd.Timedelta(days=int(30.5 * num_months)),
            freq="ME",
        )
        dates_norm = 1 + (np.arange(len(dates)) * one_month_step)
        dates_norm = dates_norm[1:]

        forecasted_values = fitted_trend_model.model(dates_norm)
        forecasted_values = scaler.inverse_transform(
            forecasted_values.reshape(-1, 1)
        ).reshape(-1)

        first_index = period - (last_fitted_date - last_date).days // 30 - 1

        periods_values = []
        sum_values = [
            last_value,
        ]

        i = first_index
        while i < len(forecasted_values):
            value = forecasted_values[i] - sum_values[-1]
            value = np.clip(value, -1, max_value * 1.5)
            if value > 0:
                periods_values.append({dates[i + 1]: value})
                sum_values.append(forecasted_values[i])
            i += period
        return periods_values
