import os
import logging
import re
import json
import math
from collections import defaultdict
from datetime import datetime
from dateutil.relativedelta import relativedelta
from pprint import pprint

import pandas as pd
import pymongo
import numpy as np
from api import predictor_pb2_grpc
from api.predictor_pb2 import (
    PredictReq,
    PredictResp,
    PrepareDataReq,
    UniqueCode,
    UniqueCodesReq,
    UniqueCodesResp,
)
from api.prompter_pb2 import QueryType
from google.protobuf.empty_pb2 import Empty
from grpc import ServicerContext, server
from concurrent import futures
from pathlib import Path
from dotenv import load_dotenv

from data import merge_contracts_with_kpgz, prepare_contracts_df, prepare_kpgz_df
from forecast_model import Model
from period_model import PeriodPredictor
from process_stock import process_and_merge_stocks, StockType, Quarters, Stock
from matcher import ColbertMatcher, YaMatcher

load_dotenv(".env")

mongo_url = f"mongodb://{os.getenv('MONGO_HOST')}:{os.getenv('MONGO_PORT')}"
mongo_client = pymongo.MongoClient(mongo_url)


def convert_to_datetime(iso_str):
    iso_str = iso_str.replace("Z", "+00:00")
    dt = datetime.strptime(iso_str, "%Y-%m-%dT%H:%M:%S.%f%z")
    return dt


def parse_filename(filename):
    pattern = r".*на\s(\d{2}\.\d{2}\.\d{4})(?:г\.?)?\s*\(сч\.\s*(\d+)\).*\.xlsx"

    match = re.match(pattern, filename)
    if match:
        date_str, account = match.groups()
        date = datetime.strptime(date_str, "%d.%m.%Y")

        quarter = (date.month - 1) // 3 + 1
        year = date.year

        return quarter, year, account
    return None


def convert_datetime_to_str(obj):
    if isinstance(obj, dict):
        return {k: convert_datetime_to_str(v) for k, v in obj.items()}
    elif isinstance(obj, list):
        return [convert_datetime_to_str(i) for i in obj]
    elif isinstance(obj, datetime):
        return obj.strftime("%Y-%m-%d")
    else:
        return obj


class Predictor(predictor_pb2_grpc.PredictorServicer):
    def __init__(self, period_model, code_matcher, name_matcher):
        self._period_model = period_model
        self._code_matcher = code_matcher
        self._name_matcher = name_matcher

    def get_merged_df(self, contracts_path, kpgz_path):
        contracts_df = pd.read_excel(contracts_path, nrows=3699)
        kpgz = pd.read_excel(kpgz_path).iloc[:1889]

        contracts_df: pd.DataFrame = prepare_contracts_df(contracts_df)
        kpgz: pd.DataFrame = prepare_kpgz_df(kpgz)
        merged_df: pd.DataFrame = merge_contracts_with_kpgz(contracts_df, kpgz)
        merged_df = merged_df.dropna(subset=["paid_rub"])

        return merged_df

    def get_main_model_outputs(self, merged_df):
        model = Model(merged_df, self._period_model)
        forecast_dict = {}
        for segment in model.filtered_segments:
            out = model.predict(
                pd.to_datetime("2023-01-01", dayfirst=True), 120, segment=segment
            )
            forecast_dict[segment] = [
                {"date": list(x.keys())[0], "value": list(x.values())[0]} for x in out
            ]

        regular_codes = set(model.filtered_segments)

        return forecast_dict, regular_codes

    def prepare_timeseries(self, ts, ref_price):
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

    def get_history(self, df: pd.DataFrame):
        hisory = df[["conclusion_date", "paid_rub"]].copy()
        hisory = (
            hisory.resample("ME", on="conclusion_date")["paid_rub"].sum().reset_index()
        )
        hisory = hisory[hisory["paid_rub"] > 0]
        return hisory.rename({"conclusion_date": "date", "paid_rub": "value"}, axis=1)

    def get_formated_purchase(
        self, cur_code_df, code_distrib, forecast, median_execution_duration
    ):
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

    def get_codes_data(self, merged_df, all_kpgz_codes, forecast_dict, regular_codes):
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

            median_execution_duration = (
                cur_code_df["execution_duration"].quantile().days
            )

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

            forecast = self.prepare_timeseries(
                forecast_dict.get(code, None), mean_ref_price
            )
            history = self.get_history(cur_code_df).to_dict(orient="records")
            history = self.prepare_timeseries(history, mean_ref_price)

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

            rows_by_date = self.get_formated_purchase(
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

    def prepare_stocks_df(self, paths):
        stocks = []
        for path in paths:
            result = parse_filename(path)
            if result:
                quarter, year, stock_type = result
                stock_type = StockType._value2member_map_.get(int(stock_type), None)
                quarter = Quarters._value2member_map_.get(int(quarter), None)
            else:
                raise ValueError(f"Wrong pattern for sotcks file: {path}")

            if stock_type is None or quarter is None:
                raise ValueError(
                    f"Not valid stock info: stock_type={stock_type}, quarter={quarter}"
                )

            raw_stock = pd.read_excel(path)
            stock = Stock(raw_stock, stock_type, quarter, year)
            stocks.append(stock)

        return process_and_merge_stocks(stocks)

    def get_stocks(self, query, organization):
        full_names = self._name_matcher.match_stocks(query)

        top_stocks = []

        collection_name = "stocks"
        collection = mongo_client[organization][collection_name]

        for i, name in enumerate(full_names):
            stock_info = collection.find({"name": name}, {"_id": False})
            stock_info = list(stock_info)
            assert len(stock_info) == 4

            top_stocks.append(
                {
                    "id": i,
                    "name": name,
                    "history": stock_info,
                }
            )

        return {"data": top_stocks}

    def get_forecast(self, product, period, organization):
        code = self._code_matcher.match_code(product)

        collection_name = "codes"
        collection = mongo_client[organization][collection_name]
        code_info = collection.find_one({"code": code}, {"_id": False})

        if code_info is None:
            pass  # TODO

        start_dt = self.cur_date
        end_dt = start_dt + relativedelta(
            years=int(period) // 12, months=int(period) % 12
        )

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
            out_id = (hash(str(forecast[0]["id"])) + hash(code) + hash(organization))%int(1e9)
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
        return code_info

    @property
    def cur_date(self):
        return convert_to_datetime("2023-01-01T10:00:20.021Z")

    def PrepareData(self, request: PrepareDataReq, context: ServicerContext):
        # в PrepareDataReq лежит путь до .csv/.xlsx файла

        logging.info(request.sources)

        assert len(request.sources) == 15

        contracts_path, kpgz_path, all_kpgz_codes_path = request.sources[:3]
        stocks = self.prepare_stocks_df(request.sources[3:])

        all_kpgz_codes = pd.read_csv(all_kpgz_codes_path)
        all_kpgz_codes = all_kpgz_codes.set_index("code")

        merged_df = self.get_merged_df(contracts_path, kpgz_path)

        forecast_dict, regular_codes = self.get_main_model_outputs(merged_df)
        codes_data = self.get_codes_data(
            merged_df, all_kpgz_codes, forecast_dict, regular_codes
        )

        collection_name = "codes"
        collection = mongo_client[request.organization][collection_name]
        collection.insert_many(codes_data)

        collection_name = "stocks"
        collection = mongo_client[request.organization][collection_name]
        collection.insert_many(stocks.to_dict(orient="records"))

        return Empty()

    def Predict(self, request: PredictReq, context: ServicerContext):
        logging.info(f"predict for {str(request)}")

        if request.type == QueryType.STOCK:
            resp = self.get_stocks(request.product, request.organization)
        elif request.type == QueryType.PREDICTION:
            resp = self.get_forecast(request.product, request.period, request.organization)

        pprint(resp)

        return PredictResp(data=json.dumps(resp).encode("utf-8"))

    def UniqueCodes(self, request: UniqueCodesReq, context: ServicerContext):
        collection_name = "codes"
        collection = mongo_client[request.organization][collection_name]
        out = collection.find(
            {"code": {"$exists": True}},
            {"_id": False, "code": True, "code_name": True, "is_regular": True},
        )

        codes_info = []
        for code_info in out:
            print(type(code_info["code_name"]))

            if (
                isinstance(code_info["code_name"], float)
                or code_info["code_name"] is None
            ):
                code_name = ""
            else:
                code_name = code_info["code_name"]

            codes_info.append(
                UniqueCode(
                    segment=code_info["code"],
                    name=code_name,
                    regular=code_info["is_regular"],
                )
            )

        return UniqueCodesResp(
            codes=codes_info,
        )


def serve():
    period_model = PeriodPredictor.load_from_checkpoint("checkpoints/periods_model")

    matcher_type = os.getenv("MATCHER")

    if matcher_type == "colbert":
        code_matcher = ColbertMatcher(
            checkpoint_name="3rd_level_codes.8bits",
            collection_path="./matcher/collections/collection_3rd_level_codes.json",
            category2code_path="./matcher/collections/category2code.json",
        )
        name_matcher = ColbertMatcher(
            checkpoint_name="full_names_stocks.8bits",
            collection_path="./matcher/collections/full_names_collection.json",
        )
    elif matcher_type == "ya":
        folder_id = os.getenv("FOLDER_ID")
        api_key = os.getenv("API_KEY")
        name_matcher = YaMatcher(
            "./matcher/embeds", folder_id=folder_id, api_key=api_key
        )
        code_matcher = name_matcher

    s = server(futures.ThreadPoolExecutor(max_workers=10))
    predictor_pb2_grpc.add_PredictorServicer_to_server(
        Predictor(
            period_model=period_model,
            code_matcher=code_matcher,
            name_matcher=name_matcher,
        ),
        s,
    )
    s.add_insecure_port("[::]:9980")
    print("starting server")
    s.start()
    s.wait_for_termination()


if __name__ == "__main__":
    serve()
