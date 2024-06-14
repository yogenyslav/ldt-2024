from api import prompter_pb2_grpc
from api.prompter_pb2 import ExtractReq, ExtractedPrompt, StreamReq, StreamResp
from grpc import ServicerContext, server
from concurrent import futures
from saiga_ollama import SaigaPrompter, SaigaOutput, PromptType
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

    def RespondStream(self, request: StreamReq, context):
        """
        generator_1 = ...
        for v in generator1:
            yield StreamResp(...)

        generator_2 = ...
        for v in generator2:
            yield StreamResp(...)
        """
        generator_1 = self.saiga.process_final_request(
            request.prompt, PromptType.FINAL_PREDICTION_PART1
        )
        for v in generator_1:
            chunk = v["message"]["content"]
            print(chunk)
            yield StreamResp(chunk=chunk)

        generator_2 = self.saiga.process_final_request(
            StreamReq.prompt, PromptType.FINAL_PREDICTION_PART2
        )
        for v in generator_2:
            chunk = v["message"]["content"]
            print(chunk)
            yield StreamResp(chunk=chunk)


def serve():
    s = server(futures.ThreadPoolExecutor(max_workers=10))
    prompter_pb2_grpc.add_PrompterServicer_to_server(Prompter(), s)
    s.add_insecure_port("[::]:9990")
    print("starting server")
    s.start()
    s.wait_for_termination()


if __name__ == "__main__":
    serve()
