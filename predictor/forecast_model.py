from typing import Any, Iterable, Optional

import numpy as np
import pandas as pd
from attrs import define
from numpy.polynomial.chebyshev import Chebyshev
from sklearn.metrics import r2_score


@define
class FittedTrendModel:
    model: Any
    one_month_step: float
    last_date: np.datetime64
    last_value: float


class Model:
    def __init__(
        self,
        df: pd.DataFrame,
        period_model: Any,
        trend_model: Optional[Any] = None,
        regular_codes: Optional[tuple] = None,
        value_column: str = "paid_rub",
        date_column: str = "conclusion_date",
    ):
        self._value_column = value_column
        self._date_column = date_column
        self._df = df.copy()

        segments_counts = self._df["depth3_code_kpgz"].value_counts()
        self._uniq_segments = self.filter_regular_codes(segments_counts, regular_codes)

        if trend_model is None:
            trend_model = Chebyshev

        self._trend_model = trend_model
        self._period_model = period_model

        self._trend_models_by_segment, self._periods_df = self.fit()

    @property
    def filtered_segments(self):
        return tuple(self._trend_models_by_segment.keys())

    def filter_regular_codes(
        self, segments: pd.Series, regular_codes: tuple, count_threshold: int = 3
    ):

        segments = segments[segments > count_threshold]
        # segments_mask = segments.index.str.startswith('01.')
        segments = segments.index.to_numpy()
        # segments = segments[segments_mask]

        if regular_codes is not None:
            segments = np.intersect1d(segments, regular_codes)

        return segments

    def train_trend_models(
        self, segments: Iterable[str], min_dates_records: int = 3, min_r2: float = 0.7
    ) -> dict[str, Any]:
        trend_models_by_segment = {}

        for segment in segments:
            filtered_df = self._df[self._df["depth3_code_kpgz"] == segment]
            filtered_df = (
                filtered_df.resample("ME", on=self._date_column)[[self._value_column]]
                .sum()
                .reset_index()
            )

            if segment == '03.07.07':
                pass
            
            if (filtered_df[self._value_column] > 0).sum() < min_dates_records:
                continue

            filtered_df[self._value_column] = filtered_df[self._value_column].cumsum()
            filtered_df["date_norm"] = (
                filtered_df[self._date_column].astype(np.int64)
                - filtered_df[self._date_column].astype(np.int64).min()
            ) / (
                filtered_df[self._date_column].astype(np.int64).max()
                - filtered_df[self._date_column].astype(np.int64).min()
            )

            X = filtered_df["date_norm"]
            y = filtered_df[self._value_column]

            fitted_model = self._trend_model.fit(X, y, 1)

            if r2_score(y, fitted_model(X)) < min_r2:
                continue

            fitted_model = FittedTrendModel(
                model=fitted_model,
                one_month_step=filtered_df["date_norm"].iloc[1]
                - filtered_df["date_norm"].iloc[0],
                last_date=filtered_df[self._date_column].iloc[-1],
                last_value=filtered_df[self._value_column].iloc[-1],
            )

            trend_models_by_segment[segment] = fitted_model

        return trend_models_by_segment

    def get_periods(self, segments: Iterable[str]):
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

    def fit(self):
        trend_models_by_segment = self.train_trend_models(self._uniq_segments)
        periods_df = self.get_periods(self._uniq_segments)

        return trend_models_by_segment, periods_df

    def predict(self, start_date: np.datetime64, num_months: int, segment: str):
        segment_period_info = self._periods_df[
            self._periods_df["depth3_code_kpgz"] == segment
        ]
        period = round(segment_period_info["period"].iloc[0])
        last_date = segment_period_info["last_date"].iloc[0]

        fitted_trend_model = self._trend_models_by_segment[segment]
        last_fitted_date = fitted_trend_model.last_date
        last_value = fitted_trend_model.last_value
        one_month_step = fitted_trend_model.one_month_step

        dates = pd.date_range(
            start=last_fitted_date,
            end=start_date + pd.Timedelta(days=int(30.5 * num_months)),
            freq="ME",
        )
        dates_norm = 1 + (np.arange(len(dates)) * one_month_step)
        dates_norm = dates_norm[1:]

        forecasted_values = fitted_trend_model.model(dates_norm)

        first_index = period - (last_fitted_date - last_date).days // 30 - 1

        periods_values = []
        sum_values = [
            last_value,
        ]

        i = first_index
        while i < len(forecasted_values):
            value = forecasted_values[i] - sum_values[-1]
            if value > 0:
                periods_values.append({dates[i + 1]: value})
                sum_values.append(forecasted_values[i])
            i += period
        return periods_values
