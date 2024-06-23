import json
import logging
import math
import os
from concurrent import futures
from datetime import datetime, timedelta
from pprint import pprint
from typing import Dict, List, Optional, Set, Tuple, Union

import pandas as pd
import pymongo
from api import predictor_pb2_grpc
from api.predictor_pb2 import (PredictReq, PredictResp, PrepareDataReq,
                               UniqueCode, UniqueCodesReq, UniqueCodesResp)
from api.prompter_pb2 import QueryType
from data_process import (get_codes_data, get_merged_df, parse_sources,
                          prepare_stocks_df, filter_forecast)
from dateutil.relativedelta import relativedelta
from dotenv import load_dotenv
from google.protobuf.empty_pb2 import Empty
from grpc import ServicerContext, server
from matcher import ColbertMatcher, YaMatcher
from utils import (convert_datetime_to_str, convert_to_datetime, get_tracer,
                   mdb_instert_many, trace_function, convert_float_nan_to_none)

from forecast_model import Model
from period_model import PeriodPredictor

load_dotenv(".env")

mongo_url = f"mongodb://{os.getenv('MONGO_HOST')}:{os.getenv('MONGO_PORT')}"
mongo_client = pymongo.MongoClient(mongo_url)

tracer = get_tracer("jaeger:4317")


class Predictor(predictor_pb2_grpc.PredictorServicer):
    def __init__(
        self,
        period_model: PeriodPredictor,
        code_matcher: Union[ColbertMatcher, YaMatcher],
        name_matcher: Union[ColbertMatcher, YaMatcher],
    ) -> None:
        """
        Initialize the Predictor class.

        Args:
            period_model (PeriodPredictor): The period prediction model.
            code_matcher (Union[ColbertMatcher, YaMatcher]): The code matcher.
            name_matcher (Union[ColbertMatcher, YaMatcher]): The name matcher.

        Returns:
            None
        """
        self._period_model = period_model
        self._code_matcher = code_matcher
        self._name_matcher = name_matcher

    @trace_function(tracer, "get_forecast_from_model")
    def get_main_model_outputs(
        self,
        merged_df: pd.DataFrame,
        date: str = "2023-01-01",
        months: int = 120,
    ) -> Tuple[Dict[str, List[Dict[str, Union[str, float]]]], Set[str]]:
        """
        Get the main model outputs for the given dataframe.

        Args:
            merged_df (pd.DataFrame): The merged dataframe.
            date (str, optional): The start date. Defaults to "2023-01-01".
            months (int, optional): The number of months. Defaults to 120.

        Returns:
            Tuple[Dict[str, List[Dict[str, Union[str, float]]]], Set[str]]: A tuple containing
            a dictionary of forecast data for each segment and a set of regular codes.
        """
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
    def get_stocks(
        self,
        query: str,
        organization: str,
    ) -> Dict[str, List[Dict[str, Union[int, str, List[Dict[str, Union[str, int]]]]]]]:
        """
        Get stocks information from MongoDB.

        Args:
            query (str): The query string for stock names.
            organization (str): The organization name.

        Returns:
            Dict[str, List[Dict[str, Union[int, str, List[Dict[str, Union[str, int]]]]]]]:
            A dictionary containing stocks information.
        """
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
    def get_forecast(
        self,
        product: str,
        period: Union[str, int],
        organization: str,
        code: Optional[str] = None,
        start_dt: Optional[str] = None,
    ) -> Dict[str, Union[int, str, List, None]]:
        """
        Get forecast information from MongoDB.

        Args:
            product (str): The product name.
            period (Union[str, int]): The period for forecast.
            organization (str): The organization name.
            code (str): The code of the product. Defaults to None.
            start_dt (str): Tghe start date of the forecast. Defaults to None.
        Returns:
            Dict[str, Union[int, str, List, None]]:
            A dictionary containing forecast and other information.
        """
        if code is None:
            code = self._code_matcher.match_code(product)

        collection_name = "codes"
        collection = mongo_client[organization][collection_name]
        code_info = collection.find_one({"code": code}, {"_id": False})

        if code_info is None:
            pass  # TODO
        
        if start_dt is None:
            start_dt = self.cur_date
        else:
            start_dt = datetime.strptime(start_dt, "%Y-%m-%d")
            
        end_dt = start_dt + relativedelta(
            years=int(period) // 12, months=int(period) % 12
        )

        return filter_forecast(code_info, organization, code, start_dt, end_dt)

    @property
    def cur_date(self) -> datetime:
        """
        Returns the current date and time as a datetime object.

        :return: A datetime object representing the current date and time.
        :rtype: datetime.datetime
        """
        return convert_to_datetime("2023-01-01T10:00:20.021Z")

    @trace_function(tracer, "prepare_data")
    def PrepareData(
        self,
        request: PrepareDataReq,
        context: ServicerContext,
    ) -> Empty:
        """
        Prepare data from Excel files and store them in MongoDB.

        Args:
            request (PrepareDataReq): The request object containing the paths to the Excel files.
            context (ServicerContext): The gRPC context object.

        Returns:
            Empty: An empty response object.
        """
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
    def Predict(
        self,
        request: PredictReq,  
        context: ServicerContext,  
    ) -> PredictResp:
        """
        Predicts stocks or forecast based on the given request.

        Args:
            request (PredictReq): The request object containing the type and other parameters.
            context (ServicerContext): The gRPC context object.

        Returns:
            PredictResp: The response object containing the predicted data.
        """
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
    def UniqueCodes(
        self,
        request: UniqueCodesReq,  
        context: ServicerContext,  
    ) -> "UniqueCodesResp":
        """
        Get unique codes from database.

        Args:
            request (UniqueCodesReq): The request object containing the organization.
            context (ServicerContext): The gRPC context object.

        Returns:
            UniqueCodesResp: The response object containing the unique codes.
        """
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
