import json
from pathlib import Path
from typing import Tuple, Union

import numpy as np
import pandas as pd
import torch
from torch import nn


class ConvNet(nn.Module):
    def __init__(self):
        super(ConvNet, self).__init__()

        self.block = nn.Sequential(
            nn.Conv1d(2, 16, kernel_size=3, padding=1),
            nn.MaxPool1d(2),
            nn.ReLU(),
            nn.Conv1d(16, 32, kernel_size=3, padding=1),
            nn.MaxPool1d(2),
            nn.ReLU(),
            nn.Conv1d(32, 32, kernel_size=5, padding=1),
            nn.ReLU(),
        )

        self.fc = nn.Sequential(
            nn.Linear(32, 16),
            nn.ReLU(),
            nn.Linear(16, 2),
        )

    def forward(self, x):
        x = self.block(x)
        x = self.fc(x.permute(0, 2, 1))
        x = x.mean((1,))
        return x


class PeriodPredictor:
    def __init__(
        self,
        model: nn.Module,
        dates_range: Tuple[str, str],
        dates_freq: str = "ME",
    ) -> None:
        """
        Initialize the PeriodPredictor.

        Args:
            model (nn.Module): The trained model for predicting periods.
            dates_range (Tuple[str, str]): The start and end dates for the prediction.
            dates_freq (str, optional): The frequency of the dates. Defaults to "ME".
        """
        self._model = model
        self._dates_freq = dates_freq

        start_date, end_date = dates_range
        dates = pd.date_range(start=start_date, end=end_date, freq=dates_freq)
        self._all_dates_df = pd.DataFrame({"date": dates})

        dates = dates.astype(int)
        self._normalized_dates = (dates - dates.min()) / (dates.max() - dates.min())

    @classmethod
    def load_from_checkpoint(cls, path: Union[str, Path]) -> "PeriodPredictor":
        """
        Load the PeriodPredictor from a checkpoint.

        Args:
            path (Union[str, Path]): The path to the checkpoint directory.

        Returns:
            PeriodPredictor: The loaded PeriodPredictor.
        """
        path = Path(path)
        model_weights = torch.load(path / "model.pt", map_location="cpu")

        model = ConvNet()
        model.load_state_dict(model_weights)

        with open(path / "metadata.json") as f:
            metadata = json.load(f)

        return cls(
            model=model,
            dates_range=metadata["dates_range"],
            dates_freq=metadata["dates_freq"],
        )

    @property
    def dates_freq(self):
        return self._dates_freq

    @torch.no_grad
    def predict(
        self,
        df: pd.DataFrame,
        dates_column: str,
        values_column: str,
        device: str = "cpu",
    ) -> Tuple[float, pd.Timestamp]:
        """
        Predict the period and last purchase date from the given DataFrame.

        Args:
            df (pd.DataFrame): The input DataFrame containing the dates and values.
            dates_column (str): The name of the column containing the dates.
            values_column (str): The name of the column containing the values.
            device (str, optional): The device to perform the prediction on. Defaults to "cpu".

        Returns:
            Tuple[float, pd.Timestamp]: A tuple containing the predicted period and the last purchase date.
        """
        self._model.to(device)
        df = df[[dates_column, values_column]].copy()
        df[dates_column] = pd.to_datetime(df[dates_column], dayfirst=True)

        df = (
            df.resample(self._dates_freq, on=dates_column)[[values_column]]
            .sum()
            .reset_index()
        )
        values = df[values_column]
        if len(np.unique(values)) == 1:
            values /= values
        else:
            df[values_column] = (values - values.min()) / (values.max() - values.min())

        expanded_df = self._all_dates_df.merge(
            df, how="left", left_on="date", right_on=dates_column
        )
        expanded_df[values_column] = expanded_df[values_column].fillna(0.0)

        values = expanded_df[values_column].to_numpy()
        assert len(values) == len(self._normalized_dates)

        X = np.stack([self._normalized_dates, values])
        X = torch.tensor(X, dtype=torch.float32, device=device)[None]

        period, last_purchase_index = self._model(X)[0]
        period = period.item()
        last_purchase_index = int(round(last_purchase_index.item()))
        last_purchase_index = (
            last_purchase_index
            if last_purchase_index < len(self._all_dates_df)
            else len(self._all_dates_df) - 1
        )
        last_purchase_date = self._all_dates_df.iloc[last_purchase_index, 0]

        return period, last_purchase_date
