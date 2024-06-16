import json
import logging
import os
from concurrent import futures
from pprint import pprint

import pandas as pd
import pymongo
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
from data_process import get_codes_data, get_merged_df, parse_sources, prepare_stocks_df
from dateutil.relativedelta import relativedelta
from dotenv import load_dotenv
from google.protobuf.empty_pb2 import Empty
from grpc import ServicerContext, server
from matcher import ColbertMatcher, YaMatcher
from utils import (
    convert_datetime_to_str,
    convert_to_datetime,
    get_tracer,
    trace_function,
    mdb_instert_many,
)

from forecast_model import Model
from period_model import PeriodPredictor

load_dotenv(".env")

mongo_url = f"mongodb://{os.getenv('MONGO_HOST')}:{os.getenv('MONGO_PORT')}"
mongo_client = pymongo.MongoClient(mongo_url)

tracer = get_tracer("jaeger:4317")


class Predictor(predictor_pb2_grpc.PredictorServicer):
    def __init__(self, period_model, code_matcher, name_matcher):
        self._period_model = period_model
        self._code_matcher = code_matcher
        self._name_matcher = name_matcher

    @trace_function(tracer, "get_forecast_from_model")
    def get_main_model_outputs(self, merged_df, date="2023-01-01", months=120):
        model = Model(merged_df, self._period_model)
        forecast_dict = {}
        for segment in model.filtered_segments:
            out = model.predict(
                pd.to_datetime(date, dayfirst=True), months, segment=segment
            )
            forecast_dict[segment] = [
                {"date": list(x.keys())[0], "value": list(x.values())[0]} for x in out
            ]

        regular_codes = set(model.filtered_segments)

        return forecast_dict, regular_codes

    @trace_function(tracer, "get_stocks_from_mdb")
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

    @trace_function(tracer, "get_forecast_from_mdb")
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
        return code_info

    @property
    def cur_date(self):
        return convert_to_datetime("2023-01-01T10:00:20.021Z")

    @trace_function(tracer, "prepare_data")
    def PrepareData(self, request: PrepareDataReq, context: ServicerContext):
        # в PrepareDataReq лежит путь до .csv/.xlsx файла

        logging.info(request.sources)
        assert len(request.sources) == 14

        contracts_path, kpgz_path, stocks_path = parse_sources(request.sources)
        all_kpgz_codes_path = "./all_kpgz_codes2name.csv"

        stocks = prepare_stocks_df(stocks_path)

        all_kpgz_codes = pd.read_csv(all_kpgz_codes_path)
        all_kpgz_codes = all_kpgz_codes.set_index("code")

        merged_df = get_merged_df(contracts_path, kpgz_path)

        forecast_dict, regular_codes = self.get_main_model_outputs(merged_df)
        codes_data = get_codes_data(
            merged_df, all_kpgz_codes, forecast_dict, regular_codes
        )

        mdb = mongo_client[request.organization]
        
        mdb_instert_many(codes_data, mdb, 'codes')
        mdb_instert_many(stocks.to_dict(orient="records"), mdb, 'stocks')

        return Empty()

    @trace_function(tracer, "predict")
    def Predict(self, request: PredictReq, context: ServicerContext):
        logging.info(f"predict for {str(request)}")

        if request.type == QueryType.STOCK:
            resp = self.get_stocks(request.product, request.organization)
        elif request.type == QueryType.PREDICTION:
            resp = self.get_forecast(
                request.product, request.period, request.organization
            )

        pprint(resp)

        return PredictResp(data=json.dumps(resp).encode("utf-8"))

    @trace_function(tracer, "unique_codes")
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
