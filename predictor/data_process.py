import math
import re
from collections import defaultdict
from datetime import datetime
from typing import Any, Dict, List, Optional, Set, Tuple, Union

import numpy as np
import pandas as pd
from dateutil.relativedelta import relativedelta
from process_stock import Quarters, Stock, StockType, process_and_merge_stocks
from scipy.optimize import linear_sum_assignment
from scipy.spatial.distance import cdist
from sentence_transformers import SentenceTransformer
from utils import convert_datetime_to_str, convert_float_nan_to_none, get_tracer, trace_function

tracer = get_tracer("jaeger:4317")

contracts_map_columns = {
    "Дата заключения": "conclusion_date",
    "Срок исполнения по": "execution_term_until",
    "Срок исполнения с": "execution_term_from",
    "Дата окончания срока действия": "end_date_of_validity",
    "Оплачено, руб.": "paid_rub",
    "Цена ГК, руб.": "gk_price_rub",
    "ID СПГЗ": "id_spgz",
    "Конечный код КПГЗ": "final_code_kpgz",
    "Наименование СПГЗ": "name_spgz",
    "Наименование (предмет) ГК": "item_name_gk",
    "Конечное наименование КПГЗ": "final_name_kpgz",
    "Реестровый номер в РК": "registry_number_in_rk",
    "Исполнено поставщиком": "provider",
    "№ версии": "version_number",
}


kpgz_map_columns = {
    "Название СТЕ": "name_ste",
    "наименование характеристик": "characteristics_name",
    "реф.цена": "ref_price",
    "Реестровый номер в РК": "registry_number_in_rk",
}


@trace_function(tracer, "get_depth34_code_kpgz_column")
def get_depth34_code_kpgz(df: pd.DataFrame,
                          min_month_count: int = 3,
                          filter_service: bool = True) -> pd.Series:
    """
    Get the depth34 code kpgz column.

    Args:
        df (pd.DataFrame): The input dataframe.
        min_month_count (int): The minimum count of months.
        filter_service (bool): Whether to filter by service.

    Returns:
        pd.Series: The depth34 code kpgz column.
    """
    init_df = df[["final_code_kpgz", "conclusion_date"]]
    df = df[["final_code_kpgz", "conclusion_date"]].copy()
    if filter_service:
        df = df[df["final_code_kpgz"].str.startswith("01.")]
    df["str_dt"] = df["conclusion_date"].dt.strftime("%Y-%m")
    df = df.drop_duplicates(subset=["final_code_kpgz", "str_dt"])

    depth4 = df["final_code_kpgz"].str.slice(0, 11).value_counts().reset_index()
    depth4["depth3"] = depth4["final_code_kpgz"].str.slice(0, 8)
    depth4["depth3_count"] = depth4.groupby(["depth3"])["count"].transform("sum")

    good_depth4 = depth4.groupby(["depth3"])["count"].apply(
        lambda x: (x > min_month_count).all()
    )
    good_depth4 = set(good_depth4[good_depth4].index.to_list())

    depth34 = np.where(
        init_df["final_code_kpgz"].str.slice(0, 8).isin(good_depth4),
        init_df["final_code_kpgz"].str.slice(0, 11),
        init_df["final_code_kpgz"].str.slice(0, 8),
    )
    return depth34


@trace_function(tracer, "prepare_contracts_df")
def prepare_contracts_df(
    df: pd.DataFrame  # type: ignore
) -> pd.DataFrame:
    """
    Prepare the contracts dataframe.

    Args:
        df (pd.DataFrame): The input dataframe.

    Returns:
        pd.DataFrame: The prepared dataframe.
    """
    df = df[contracts_map_columns.keys()]
    df = df.rename(contracts_map_columns, axis=1)
    df = df.dropna(axis=0, subset=["final_code_kpgz", "conclusion_date"])

    df["conclusion_date"] = pd.to_datetime(df["conclusion_date"], dayfirst=True)
    df["execution_term_until"] = pd.to_datetime(
        df["execution_term_until"], dayfirst=True
    )
    df["execution_term_from"] = pd.to_datetime(df["execution_term_from"], dayfirst=True)
    df["end_date_of_validity"] = pd.to_datetime(
        df["end_date_of_validity"], dayfirst=True
    )

    df["depth3_code_kpgz"] = get_depth34_code_kpgz(df)  # TODO change names
    df["depth3_code_kpgz"] = df["depth3_code_kpgz"].astype(pd.StringDtype())

    df["registry_number_in_rk"] = df["registry_number_in_rk"].astype(pd.StringDtype())
    df = df.set_index("registry_number_in_rk", drop=False)

    df["id"] = np.arange(len(df))
    df["id"] = df["id"].convert_dtypes()

    return df


@trace_function(tracer, "prepare_kpgz_df")
def prepare_kpgz_df(df: pd.DataFrame) -> pd.DataFrame:
    """
    Prepare the kpgz dataframe.

    Args:
        df (pd.DataFrame): The input dataframe.

    Returns:
        pd.DataFrame: The prepared dataframe.
    """
    df = df[kpgz_map_columns.keys()]
    df = df.rename(kpgz_map_columns, axis=1)
    df = df.dropna(axis=0, subset=["name_ste"])

    df["registry_number_in_rk"] = df["registry_number_in_rk"].astype(pd.StringDtype())
    df = df.set_index("registry_number_in_rk", drop=True)

    df["id"] = np.arange(len(df))
    df["id"] = df["id"].convert_dtypes()

    return df


@trace_function(tracer, "merge_contracts_with_kpgz")
def merge_contracts_with_kpgz(
    contracts_df: pd.DataFrame, kpgz_df: pd.DataFrame
) -> pd.DataFrame:
    """
    Merge the contracts dataframe with the kpgz dataframe using the sentence transformers library.

    Args:
        contracts_df (pd.DataFrame): The contracts dataframe.
        kpgz_df (pd.DataFrame): The kpgz dataframe.

    Returns:
        pd.DataFrame: The merged dataframe.
    """
    model_name = "sentence-transformers/distilbert-multilingual-nli-stsb-quora-ranking"
    encoder = SentenceTransformer(model_name)

    contracts_names = contracts_df["name_spgz"].astype(pd.StringDtype())
    kpgz_names = kpgz_df["name_ste"].astype(pd.StringDtype())

    contracts_embedding = encoder.encode(contracts_names.tolist())
    kpgz_embedding = encoder.encode(kpgz_names.tolist())

    uniq_keys1 = contracts_df.index.values.unique().to_numpy()
    uniq_keys2 = kpgz_df.index.values.unique().to_numpy()

    uniq_keys1 = uniq_keys1[~pd.isna(uniq_keys1)]
    uniq_keys2 = uniq_keys2[~pd.isna(uniq_keys2)]

    uniq_keys = np.intersect1d(uniq_keys1, uniq_keys2, assume_unique=True)
    rest_keys_contracts = np.array(list(set(uniq_keys1) - set(uniq_keys)), dtype=object)

    to_concat_rows = []

    for key in uniq_keys:
        row1 = contracts_df.loc[key]
        row2 = kpgz_df.loc[key]

        if not isinstance(row1, pd.DataFrame):
            row1 = pd.DataFrame(row1).T
        if not isinstance(row2, pd.DataFrame):
            row2 = pd.DataFrame(row2).T

        embed_ids1 = row1["id"].to_numpy(dtype=int)
        embed_ids2 = row2["id"].to_numpy(dtype=int)

        embeds1 = contracts_embedding[embed_ids1]
        embeds2 = kpgz_embedding[embed_ids2]

        if embeds1.shape[0] == 1 and embeds2.shape[0] == 1:
            row_ind = np.array(
                [
                    0,
                ]
            )
            col_ind = np.array(
                [
                    0,
                ]
            )
        else:
            distance_matrix = cdist(embeds1, embeds2, metric="cosine")
            row_ind, col_ind = linear_sum_assignment(distance_matrix)

        matched_row1 = row1.iloc[row_ind]
        matched_row2 = row2.iloc[col_ind]

        concated_row = pd.concat(
            [matched_row1, matched_row2.drop(columns=["id"])], axis=1
        )

        to_concat_rows.append(concated_row)

    merged_df = pd.concat(to_concat_rows, axis=0)
    merged_df = pd.concat(
        [merged_df, contracts_df.loc[rest_keys_contracts]], axis=0
    ).convert_dtypes()
    merged_df = merged_df.sort_index()

    return merged_df


@trace_function(tracer, "get_merged_df")
def get_merged_df(
    contracts_path: str, kpgz_path: str
) -> pd.DataFrame:
    """
    Get the merged dataframe using contracts and kpgz dataframes.

    Args:
        contracts_path: The path to the contracts dataframe.
        kpgz_path: The path to the kpgz dataframe.

    Returns:
        pd.DataFrame: The merged dataframe.
    """
    contracts_df = pd.read_excel(contracts_path, nrows=3699)
    kpgz = pd.read_excel(kpgz_path).iloc[:1889]

    contracts_df: pd.DataFrame = prepare_contracts_df(contracts_df)
    kpgz: pd.DataFrame = prepare_kpgz_df(kpgz)
    merged_df: pd.DataFrame = merge_contracts_with_kpgz(contracts_df, kpgz)
    merged_df = merged_df.dropna(subset=["paid_rub"])

    return merged_df


@trace_function(tracer, "get_code_history")
def get_code_history(
    df: pd.DataFrame
) -> pd.DataFrame:
    """
    Get the code history from the given dataframe.

    Args:
        df (pd.DataFrame): The input dataframe containing "conclusion_date" and "paid_rub" columns.

    Returns:
        pd.DataFrame: The code history dataframe with "date" and "value" columns.
    """
    hisory = df[["conclusion_date", "paid_rub"]].copy()
    hisory = hisory.resample("ME", on="conclusion_date")["paid_rub"].sum().reset_index()
    hisory = hisory[hisory["paid_rub"] > 0]
    return hisory.rename({"conclusion_date": "date", "paid_rub": "value"}, axis=1)


@trace_function(tracer, "prepare_timeseries")
def prepare_timeseries(
    ts: Optional[List[Dict[str, Union[str, float]]]],
    ref_price: float,
) -> Optional[List[Dict[str, Union[int, str, float, None]]]]:
    """
    Prepare the timeseries data.

    Args:
        ts (Optional[List[Dict[str, Union[str, float]]]]): The timeseries data.
        ref_price (Union[float, pd.NA, np.nan]): The reference price.

    Returns:
        Optional[List[Dict[str, Union[int, str, float, None]]]]: The prepared timeseries data.
    """
    if ts is None:
        return None

    new_ts = []
    for i, d in enumerate(ts):
        if (
            ref_price is not pd.NA
            and not math.isnan(d["value"])
            and ref_price is not np.nan
        ):
            new_ts.append(
                {
                    "id": i,
                    "date": d["date"],
                    "value": d["value"],
                    "volume": int(d["value"] // ref_price),
                }
            )

        else:
            new_ts.append(
                {
                    "id": i,
                    "date": d["date"],
                    "value": d["value"],
                    "volume": None,
                }
            )

    return new_ts


@trace_function(tracer, "get_formated_purchase")
def get_formated_purchase(
    cur_code_df: pd.DataFrame,
    code_distrib: pd.DataFrame,
    forecast: Optional[List[Dict[str, Union[str, float]]]],
    median_execution_duration: Optional[int] = 30,
) -> Optional[Dict[str, List[Dict[str, Union[str, int, float]]]]]:
    """
    Get the formatted purchase data.

    Args:
        cur_code_df (pd.DataFrame): The current code dataframe.
        code_distrib (pd.DataFrame): The code distribution dataframe.
        forecast (Optional[List[Dict[str, Union[str, float]]]]): The forecast data.
        median_execution_duration (Optional[int], optional): The median execution duration. Defaults to 30.

    Returns:
        Optional[Dict[str, List[Dict[str, Union[str, int, float]]]]]: The formatted purchase data.
    """
    if forecast is None:
        return None
    rows_by_date = defaultdict(list)

    median_execution_duration = (
        median_execution_duration if median_execution_duration is not None else 30
    )

    for code, fraction in zip(
        code_distrib["final_code_kpgz"], code_distrib["paid_rub"]
    ):
        code_df = cur_code_df[
            ["name_spgz", "id_spgz", "version_number", "final_name_kpgz"]
        ][cur_code_df["final_code_kpgz"] == code].dropna(axis=1)
        name_spgz = code_df["name_spgz"].iloc[0]
        id_spgz = int(code_df["id_spgz"].iloc[0])
        version_number = int(code_df["version_number"].iloc[0])
        final_name_kpgz = code_df["final_name_kpgz"].iloc[0]

        for forecast_point in forecast:
            start_date = forecast_point["date"]
            end_date = forecast_point["date"] + relativedelta(
                days=median_execution_duration
            )
            year = end_date.year
            if forecast_point["volume"] is None:
                amount = None
            else:
                amount = int(forecast_point["volume"] * fraction)

            nmc = forecast_point["value"] * fraction
            raw = {
                "DeliverySchedule": {
                    "start_date": start_date,
                    "end_date": end_date,
                    "year": year,
                    "deliveryAmount": amount,
                },
                "entityId": id_spgz,
                "id": version_number,
                "nmc": nmc,
                "purchaseAmount": amount,
                "spgzCharacteristics": {
                    "spgzId": id_spgz,
                    "spgzName": name_spgz,
                    "kpgzCode": code,
                    "kpgzName": final_name_kpgz,
                },
            }

            rows_by_date[convert_datetime_to_str(start_date)].append(raw)

    return rows_by_date


@trace_function(tracer, "get_codes_data")
def get_codes_data(
    merged_df: pd.DataFrame,
    all_kpgz_codes: pd.DataFrame,
    forecast_dict: Dict[str, List[Dict[str, Any]]],
    regular_codes: Set[str],
) -> List[Dict[str, Union[str, int, float, bool, List[Dict[str, Any]], Dict[str, Any], None]]]:
    """
    Process the merged dataframe and return the codes data.

    Args:
        merged_df (pd.DataFrame): The merged dataframe.
        all_kpgz_codes (pd.DataFrame): The dataframe containing all kpgz codes.
        forecast_dict (Dict[str, List[Dict[str, Any]]]): The forecast dictionary.
        regular_codes (Set[str]): The set of regular codes.

    Returns:
        List[Dict[str, Union[str, int, float, bool, List[Dict[str, Any]], Dict[str, Any], None]]]: The codes data.
    """
    codes_stat = merged_df
    codes_stat["execution_duration"] = (
        codes_stat["execution_term_until"] - codes_stat["execution_term_from"]
    )
    codes_stat["start_to_execute_duration"] = (
        codes_stat["execution_term_from"] - codes_stat["conclusion_date"]
    )
    codes_stat["num_nans"] = codes_stat.isna().sum(axis=1)

    codes_stat = codes_stat.sort_values(by="num_nans")
    codes_stat = codes_stat.set_index(
        [
            "depth3_code_kpgz",
        ],
        drop=True,
    )

    codes_data = []
    for code in codes_stat.index.unique():
        try:
            code_name = all_kpgz_codes.loc[code, "name"]
        except KeyError:
            print(f"code: {code} is not found")
            code_name = None

        cur_code_df = codes_stat.loc[code]
        if isinstance(cur_code_df, pd.Series):
            cur_code_df = cur_code_df.to_frame().T

        median_execution_duration = cur_code_df["execution_duration"].quantile().days

        if (cur_code_df["start_to_execute_duration"] < pd.Timedelta(days=0)).any():
            mean_start_to_execute_duration = None
        elif len(cur_code_df) == 1:
            mean_start_to_execute_duration = (
                cur_code_df["start_to_execute_duration"].iloc[0].days
            )
        else:
            mean_start_to_execute_duration = (
                cur_code_df["start_to_execute_duration"].mean().days
            )

        mean_ref_price = cur_code_df["ref_price"].mean()
        top5_providers = cur_code_df["provider"].value_counts()[:5].index.tolist()

        cur_code_df = cur_code_df.drop(
            ["num_nans", "start_to_execute_duration", "execution_duration"], axis=1
        )

        forecast = prepare_timeseries(forecast_dict.get(code, None), mean_ref_price)
        history = get_code_history(cur_code_df).to_dict(orient="records")
        history = prepare_timeseries(history, mean_ref_price)

        examples = (
            cur_code_df.drop_duplicates(subset=["final_code_kpgz"])
            .iloc[:5]
            .to_dict(orient="records")
        )

        code_distrib = (
            cur_code_df.groupby("final_code_kpgz")["paid_rub"]
            .sum()
            .sort_values(ascending=False)[:20]
        )
        code_distrib = (code_distrib / code_distrib.sum()).reset_index()

        rows_by_date = get_formated_purchase(
            cur_code_df, code_distrib, forecast, median_execution_duration
        )

        codes_data.append(
            {
                "code": code,
                "code_name": code_name,
                "is_regular": code in regular_codes,
                "forecast": forecast,
                "history": history,
                "rows_by_date": rows_by_date,
                "median_execution_days": (
                    median_execution_duration
                    if median_execution_duration is not pd.NA
                    else 30
                ),
                "mean_start_to_execute_days": (
                    mean_start_to_execute_duration
                    if mean_start_to_execute_duration is not pd.NA
                    else None
                ),
                "mean_ref_price": (
                    mean_ref_price if mean_ref_price is not pd.NA else None
                ),
                "top5_providers": (
                    top5_providers if top5_providers is not pd.NA else None
                ),
                "example_contracts_in_code": examples,
            }
        )

    return codes_data


@trace_function(tracer, "prepare_stocks_df")
def prepare_stocks_df(
    paths: List
) -> pd.DataFrame:
    """
    Prepare stocks data from list of PathInfo objects.

    Args:
        paths (List): List of Path objects.

    Returns:
        pd.DataFrame: Processed and merged stocks data.
    """
    stocks = []
    for obj in paths:
        result = parse_filename(obj.name)
        if result:
            quarter, year, stock_type = result
            stock_type = StockType._value2member_map_.get(int(stock_type), None)
            quarter = Quarters._value2member_map_.get(int(quarter), None)
        else:
            raise ValueError(f"Wrong pattern for sotcks file: {obj.name}")

        if stock_type is None or quarter is None:
            raise ValueError(
                f"Not valid stock info: stock_type={stock_type}, quarter={quarter}"
            )

        raw_stock = pd.read_excel(obj.path)
        stock = Stock(raw_stock, stock_type, quarter, year)
        stocks.append(stock)

    return process_and_merge_stocks(stocks)


@trace_function(tracer, "parse_sources")
def parse_sources(
    sources: List
) -> Tuple[str, str, List]:
    """
    Parse sources and return paths to contracts, kpgz and stocks.

    Args:
        sources (List[PathInfo]): List of PathInfo objects.

    Returns:
        Tuple[str, str, List[PathInfo]]: Tuple with paths to contracts, kpgz and stocks.
    """
    contracts_path = [
        x.path for x in sources if x.name.startswith("Выгрузка контрактов по Заказчику")
    ][0]
    kpgz_path = [x.path for x in sources if x.name.startswith("КПГЗ ,СПГЗ, СТЕ")][0]
    stocks_path = [x for x in sources if x.name.startswith("Ведомость остатков")]

    return contracts_path, kpgz_path, stocks_path


@trace_function(tracer, "parse_filename")
def parse_filename(filename: str) -> Optional[Tuple[int, int, str]]:
    """
    Parses the filename and returns the quarter, year, and account number.

    Args:
        filename (str): The name of the file to parse.

    Returns:
        Optional[Tuple[int, int, str]]: A tuple containing the quarter, year, and account number,
        or None if the filename does not match the expected pattern.
    """
    pattern = r".*на\s(\d{2}\.\d{2}\.\d{4})(?:г\.?)?\s*\(сч\.\s*(\d+)\).*\.xlsx"

    match = re.match(pattern, filename)
    if match:
        date_str, account = match.groups()
        date = datetime.strptime(date_str, "%d.%m.%Y")

        quarter = (date.month - 1) // 3 + 1
        year = date.year

        return quarter, year, account
    return None

@trace_function(tracer, "filter_forecast")
def filter_forecast(
    code_info: Dict,
    organization: str,
    code: str,
    start_dt: datetime,
    end_dt: datetime,
):
    rows_by_date = code_info.pop("rows_by_date")

    if code_info["forecast"] is not None:
        forecast_start_filtered = [
            x
            for x in code_info["forecast"]
            if start_dt.timestamp() <= x["date"].timestamp()
        ]
        closest_purchase = forecast_start_filtered[0]
        forecast = [
            x
            for x in forecast_start_filtered
            if x["date"].timestamp() <= end_dt.timestamp()
        ]

    else:
        forecast = None
        closest_purchase = None

    if forecast is None:
        rows = None
        out_id = None
    elif len(forecast) > 0:
        rows = rows_by_date.get(convert_datetime_to_str(forecast[0]["date"]), None)
        out_id = (
            hash(str(forecast[0]["id"])) + hash(code) + hash(organization)
        ) % int(1e9)
    else:
        rows = []
        out_id = None

    output_json = {
        "id": out_id,
        "CustomerId": organization,
        "rows": rows,
    }

    code_info["forecast"] = forecast
    code_info["output_json"] = output_json
    code_info["closest_purchase"] = closest_purchase

    code_info = convert_datetime_to_str(code_info)
    code_info = convert_float_nan_to_none(code_info)
    return code_info