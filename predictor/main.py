import os
import logging

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

from data import (merge_contracts_with_kpgz, prepare_contracts_df,
                  prepare_kpgz_df)
from forecast_model import Model
from period_model import PeriodPredictor

load_dotenv(".env")

mongo_url = f"mongodb://{os.getenv('MONGO_HOST')}:{os.getenv('MONGO_PORT')}/{os.getenv('MONGO_DB')}"
mongo_client = pymongo.MongoClient(mongo_url)
mongo_db = mongo_client[os.getenv('MONGO_DB')]

class Predictor(predictor_pb2_grpc.PredictorServicer):
    def __init__(self, period_model):
        self._period_model = period_model

    def PrepareData(self, request: PrepareDataReq, context: ServicerContext):
        # в PrepareDataReq лежит путь до .csv/.xlsx файла

        logging.info(request.sources)
        
        assert len(request.sources) == 2
        
        contracts_path, kpgz_path = request.sources
        
        contracts_df = pd.read_excel(contracts_path, nrows=3699)
        kpgz = pd.read_excel(kpgz_path).iloc[:1889]

        contracts_df: pd.DataFrame = prepare_contracts_df(contracts_df)
        kpgz: pd.DataFrame = prepare_kpgz_df(kpgz)
        merged_df: pd.DataFrame = merge_contracts_with_kpgz(contracts_df, kpgz)
        merged_df = merged_df.dropna(subset=['paid_rub'])
        
        model = Model(merged_df, self._period_model)
        forecast_dicts = []
        for segment in model.filtered_segments:
            out = model.predict(
                pd.to_datetime("2023-01-01", dayfirst=True), 120, segment=segment
            )
            forecast_dicts.append({segment.replace('.', '_'): [{'date': list(x.keys())[0], 'value': list(x.values())[0]} for x in out]}) 

        uniq_segments = merged_df["depth3_code_kpgz"].unique()
        regular_mask = uniq_segments.isin(model.filtered_segments)
        regular_dicts = [{'segment': x[0].replace('.', '_'), 'is_regular': x[1].item()} for x in zip(uniq_segments, regular_mask)]
        
        data_dict = merged_df.to_dict(orient="records")
        
        collection_name = 'contracts_with_kpgz'
        collection = mongo_db[collection_name]
        collection.insert_many(data_dict)
        
        collection_name = 'forecast'
        collection = mongo_db[collection_name]
        collection.insert_many(forecast_dicts)
        
        collection_name = 'segments'
        collection = mongo_db[collection_name]
        collection.insert_many(regular_dicts)
        
        # 1. prepare data here -> mongo
        # 2. unique codes here -> mongo
        # 3. predict here -> mongo
        return ClientIdentifier(value="")

    def Predict(self, request: PredictReq, context: ServicerContext):
        logging.info(
            f"predict for ts={request.ts.ToJsonString()} months_count={request.months_count} segment={request.segment}"
        )
        
        collection_name = 'segments'
        collection = mongo_db[collection_name]
        collection.find_one
        
        return PredictResp(predicts=[])

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
