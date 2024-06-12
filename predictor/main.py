import os
import logging
import re
from datetime import datetime
from dateutil.relativedelta import relativedelta

import pandas as pd
import pymongo
from api import predictor_pb2_grpc
from api.predictor_pb2 import (
    PredictReq,
    PredictResp,
    PrepareDataReq,
    ClientIdentifier,
    UniqueCodesResp,
)
from grpc import ServicerContext, server
from concurrent import futures
from pathlib import Path
from dotenv import load_dotenv

from data import merge_contracts_with_kpgz, prepare_contracts_df, prepare_kpgz_df
from forecast_model import Model
from period_model import PeriodPredictor
from process_stock import process_and_merge_stocks, StockType, Quarters, Stock
from matcher import ColbertMatcher

load_dotenv(".env")

mongo_url = f"mongodb://{os.getenv('MONGO_HOST')}:{os.getenv('MONGO_PORT')}/{os.getenv('MONGO_DB')}"
mongo_client = pymongo.MongoClient(mongo_url)
mongo_db = mongo_client[os.getenv("MONGO_DB")]

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


class Predictor(predictor_pb2_grpc.PredictorServicer):
    def __init__(self, period_model):
        self._period_model = period_model
        self._code_matcher = ColbertMatcher(
            checkpoint_name="3rd_level_codes.8bits",
            collection_path="./matcher/collections/collection_3rd_level_codes.json",
            category2code_path="./matcher/collections/category2code.json"
        )
        self._name_matcher = ColbertMatcher(
            checkpoint_name="full_names_stocks.8bits",
            collection_path="./matcher/collections/full_names_collection.json",
        )

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

            forecast = forecast_dict.get(code, None)

            all_data = cur_code_df.to_dict(orient="records")
            codes_data.append(
                {
                    "code": code,
                    "code_name": code_name,
                    "is_regular": code in regular_codes,
                    "forecast": forecast,
                    "median_execution_days": (
                        median_execution_duration
                        if median_execution_duration is not pd.NA
                        else None
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
                    "contracts_in_code": all_data,
                }
            )

        return codes_data

    def get_stocks(self, paths):
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

    def PrepareData(self, request: PrepareDataReq, context: ServicerContext):
        print(self._code_matcher.match_to_3rd_level_code('вода'))
        # в PrepareDataReq лежит путь до .csv/.xlsx файла

        logging.info(request.sources)

        assert len(request.sources) == 15

        contracts_path, kpgz_path, all_kpgz_codes_path = request.sources[:3]
        stocks = self.get_stocks(request.sources[3:])

        all_kpgz_codes = pd.read_csv(all_kpgz_codes_path)
        all_kpgz_codes = all_kpgz_codes.set_index("code")

        merged_df = self.get_merged_df(contracts_path, kpgz_path)

        forecast_dict, regular_codes = self.get_main_model_outputs(merged_df)
        codes_data = self.get_codes_data(
            merged_df, all_kpgz_codes, forecast_dict, regular_codes
        )

        collection_name = "codes"
        collection = mongo_db[collection_name]
        collection.insert_many(codes_data)

        collection_name = "stocks"
        collection = mongo_db[collection_name]
        collection.insert_many(stocks.to_dict(orient="records"))

        return ClientIdentifier(value="")

    def Predict(self, request: PredictReq, context: ServicerContext):
        logging.info(
            f"predict for ts={request.ts.ToJsonString()} months_count={request.months_count} segment={request.segment}"
        )
        # матчинг 
        collection_name = "codes"
        collection = mongo_db[collection_name]
        code_info = collection.find_one({"code": request.segment},{'_id':False})
        
        if code_info is None:
            pass #TODO
        
        start_dt = convert_to_datetime(request.ts.ToJsonString())
        end_dt = start_dt + relativedelta(years=request.months_count//12, months=request.months_count%12)
        forecast = [x for x in code_info['forecast'] if start_dt.timestamp() <= x['date'].timestamp() <= end_dt.timestamp() ]
        code_info['forecast'] = forecast
        
        return PredictResp(predicts=[]) #TODO

    def UniqueCodes(self, request: ClientIdentifier, context: ServicerContext):
        logging.info(f"unique codes for {request.value}")
        return UniqueCodesResp(codes=[])


def serve():
    period_model = PeriodPredictor.load_from_checkpoint("checkpoints/periods_model")

    s = server(futures.ThreadPoolExecutor(max_workers=10))
    predictor_pb2_grpc.add_PredictorServicer_to_server(Predictor(period_model), s)
    s.add_insecure_port("[::]:9980")
    print("starting server")
    s.start()
    s.wait_for_termination()


if __name__ == "__main__":
    serve()
