import pandas as pd
import numpy as np

from dataclasses import dataclass
from enum import Enum
import re

class StockType(Enum):
    ACC_105 = 105
    ACC_101 = 101
    ACC_21 = 21

class Quarters(Enum):
    Q1 = 1
    Q2 = 2
    Q3 = 3
    Q4 = 4 

@dataclass
class Stock:
    df: pd.DataFrame
    stock_type: StockType
    quarter: Quarters
    year: str

def extract_numbers(s):
    if not pd.isna(s):
        result = re.findall(r"(\d+\.\d+|\d+)", s)
        if not result or s != result[0]:
            return None 
        else:
            return re.findall(r"(\d+\.\d+|\d+)", s)[0]
    else:
        return None
    
def filter_rows(
        df: pd.DataFrame,
        vals_to_filter: list[str],
        filter_by: str,
) -> pd.DataFrame:
    
    return df[~df[filter_by].isin(vals_to_filter)]

def remove_nums(
        df: pd.DataFrame,
        filter_by: str,
) -> pd.DataFrame:
    
    df_processed = df.copy()
    for i, row in df.iterrows():
        numbers = extract_numbers(row[filter_by])
        if numbers:
            df_processed.drop(i, inplace=True)
    
    return df_processed

def remove_duplicated(
        df: pd.DataFrame,
        filter_by: str,
) -> pd.DataFrame:
    
    return df.groupby(filter_by, as_index=False).sum()

def rename_columns(
        df: pd.DataFrame,
        column_names: dict[str, str],
) -> pd.DataFrame:
    
    return df.rename(columns=column_names)

def rename_and_drop_cols(
        df: pd.DataFrame,
) -> pd.DataFrame:
    
    df.rename(
        columns={
            "Unnamed: 2": "name",
            "Unnamed: 19": "sum",
            "Unnamed: 20": "amount",
        },
        inplace=True,
    )

    df.dropna(subset=["name", "amount"], inplace=True)

    df["price"] = df["sum"] / df["amount"].replace(0, np.nan)
    df["price"].fillna(0, inplace=True)

    return df[["name", "price", "amount", "sum"]].groupby("name", as_index=False).sum()

def process_stock_type105(
        df: pd.DataFrame,
) -> pd.DataFrame:

    filter = ["Место хранения", "Счет", "КФО", "Номенклатура", "Кузнецов Максим Сергеевич", "Итого"]
    df_processed = df.copy()
    df = filter_rows(df_processed, filter, "МОЛ")    
    df = remove_nums(df_processed, "МОЛ")
    df = remove_duplicated(df_processed, "МОЛ")

    rename_dict = {
    "МОЛ": "name",
    "Цена": "price",
    "Количество": "amount",
    "Сумма": "sum",
    }

    df_processed = rename_columns(df_processed, rename_dict)
    return df_processed

def parse_stock_type101_21(
        df: pd.DataFrame,
) -> pd.DataFrame:
    
    return rename_and_drop_cols(df)

def process_stock(
        raw_stock: Stock,
) -> Stock:
    
    if raw_stock.stock_type == StockType.ACC_105:
        raw_stock.df = process_stock_type105(raw_stock.df)
    elif raw_stock.stock_type == StockType.ACC_101 or raw_stock.stock_type == StockType.ACC_21:
        raw_stock.df = parse_stock_type101_21(raw_stock.df)

    return raw_stock

def process_and_merge_stocks(
        stocks: list[Stock],
) -> pd.DataFrame:
    
    new_stocks = [process_stock(stock) for stock in stocks]

    for stock in new_stocks:
        stock.df.loc[:, "year"] = stock.year
        stock.df.loc[:, "quarter"] = stock.quarter.value

    df = pd.concat([stock.df for stock in stocks], axis=0)
    

    unique_dates = df[["year", "quarter"]].drop_duplicates()
    for name in df["name"].unique():
        for year, quarter in unique_dates.values:
            if not df[(df["name"] == name) & (df["year"] == year) & (df["quarter"] == quarter)].empty:
                continue
            df = pd.concat([df, pd.DataFrame({"name": [name], "year": [year], "quarter": [quarter], "price": [0], "amount": [0], "sum": [0]})], axis=0)
            # df = df.append({"name": name, "year": year, "quarter": quarter, "price": 0, "amount": 0, "sum": 0}, ignore_index=True)


    df.reset_index(drop=True, inplace=True)
    
    return df
