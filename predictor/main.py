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


class Predictor(predictor_pb2_grpc.PredictorServicer):
    def __init__(self):
        pass

    def PrepareData(self, request: PrepareDataReq, context: ServicerContext):
        # в PrepareDataReq лежит путь до .csv/.xlsx файла

        print(f"prepare data for {request.source}")
        # 1. prepare data here -> mongo
        # 2. unique codes here -> mongo
        # 3. predict here -> mongo
        return ClientIdentifier(value="")

    def Predict(self, request: PredictReq, context: ServicerContext):
        print(
            f"predict for ts={request.ts.ToJsonString()} months_count={request.months_count} segment={request.segment}"
        )
        return PredictResp(predicts=[])

    def UniqueCodes(self, request: ClientIdentifier, context: ServicerContext):
        print(f"unique codes for {request.value}")
        return UniqueCodesResp(codes=[])


def serve():
    s = server(futures.ThreadPoolExecutor(max_workers=10))
    predictor_pb2_grpc.add_PredictorServicer_to_server(Predictor(), s)
    s.add_insecure_port("[::]:9980")
    print("starting server")
    s.start()
    s.wait_for_termination()


if __name__ == "__main__":
    serve()
