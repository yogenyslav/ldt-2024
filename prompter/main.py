from api import prompter_pb2_grpc
from api.prompter_pb2 import ExtractReq, ExtractedPrompt, QueryType
from grpc import ServicerContext, server
from concurrent import futures


class Prompter(prompter_pb2_grpc.PrompterServicer):
    def Extract(self, request: ExtractReq, context: ServicerContext) -> ExtractedPrompt:
        # prompter action here
        return ExtractedPrompt(type=QueryType.UNDEFINED, product="test")


def serve():
    s = server(futures.ThreadPoolExecutor(max_workers=10))
    prompter_pb2_grpc.add_PrompterServicer_to_server(Prompter(), s)
    s.add_insecure_port("[::]:9990")
    print("starting server")
    s.start()
    s.wait_for_termination()


if __name__ == "__main__":
    serve()
