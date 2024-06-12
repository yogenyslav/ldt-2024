from api import prompter_pb2_grpc
from api.prompter_pb2 import ExtractReq, ExtractedPrompt, QueryType
from grpc import ServicerContext, server
from concurrent import futures
from saiga_ollama import SaigaPrompter, SaigaOutput
from pathlib import Path


class Prompter(prompter_pb2_grpc.PrompterServicer):

    def __init__(self):

        self.prompts_path = Path(__file__).parent / "prompts.json"
        self.saiga = SaigaPrompter(prompts_path=self.prompts_path)

    def Extract(self, request: ExtractReq, context: ServicerContext) -> ExtractedPrompt:
        # prompter action here
        print(f"got prompt {request.prompt}")
        output = self.saiga.process_request(request.prompt)
        return ExtractedPrompt(
            type=output.type, product=output.product, period=output.period
        )


def serve():
    s = server(futures.ThreadPoolExecutor(max_workers=10))
    prompter_pb2_grpc.add_PrompterServicer_to_server(Prompter(), s)
    s.add_insecure_port("[::]:9990")
    print("starting server")
    s.start()
    s.wait_for_termination()


if __name__ == "__main__":
    serve()
