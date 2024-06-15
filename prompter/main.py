from api import prompter_pb2_grpc
from api.prompter_pb2 import ExtractReq, ExtractedPrompt, StreamReq, StreamResp
from grpc import ServicerContext, server
from concurrent import futures
from saiga_ollama import SaigaPrompter, SaigaOutput, PromptType
from yagpt_prompter import YaGPTPrompter, PrompterOutput
from pathlib import Path
from dotenv import load_dotenv
import os
import logging
import time
import json

load_dotenv(".env")


class Prompter(prompter_pb2_grpc.PrompterServicer):

    def __init__(self):
        self.prompts_path = Path(__file__).parent / "prompts.json"
        self.model_choice = os.getenv("MODEL_CHOICE")
        if self.model_choice == "saiga":
            self.model = SaigaPrompter(prompts_path=self.prompts_path)
        elif self.model_choice == "yandexgpt":
            self.model = YaGPTPrompter(prompts_path=self.prompts_path)

    def Extract(self, request: ExtractReq, context: ServicerContext) -> ExtractedPrompt:
        # prompter action here
        print(f"got prompt {request.prompt}")
        output = self.model.process_request(request.prompt)
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
        if self.model_choice == "saiga":
            generator_1 = self.model.process_final_request(
                request.prompt, PromptType.FINAL_PREDICTION_PART1
            )
            for v in generator_1:
                chunk = v["message"]["content"]
                print(chunk)
                yield StreamResp(chunk=chunk)

            generator_2 = self.model.process_final_request(
                request.prompt, PromptType.FINAL_PREDICTION_PART2
            )
            for v in generator_2:
                chunk = v["message"]["content"]
                print(chunk)
                yield StreamResp(chunk=chunk)
        elif self.model_choice == "yandexgpt":
            generator_1 = self.model.process_final_request(
                request.prompt, PromptType.FINAL_PREDICTION_PART1
            )
            last_msg = ""
            for line in generator_1.iter_lines():
                if line:
                    decoded_line = line.decode("utf-8")
                    message = json.loads(decoded_line)
                    print(message)
                    msg_text = message["result"]["alternatives"][0]["message"]["text"]
                    outp_text = msg_text[len(last_msg) :]
                    last_msg = msg_text
                    for c in outp_text:
                        time.sleep(0.01)
                        yield StreamResp(chunk=c)


logging.basicConfig(level=logging.DEBUG)


def serve():
    s = server(futures.ThreadPoolExecutor(max_workers=10))
    prompter_pb2_grpc.add_PrompterServicer_to_server(Prompter(), s)
    s.add_insecure_port("[::]:9990")
    print("starting server")
    s.start()
    s.wait_for_termination()


if __name__ == "__main__":
    serve()
