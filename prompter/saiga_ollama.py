import ollama
import json
import os
from dataclasses import dataclass
from api.prompter_pb2 import QueryType
from dotenv import load_dotenv

load_dotenv(".env")


@dataclass
class SaigaOutput:
    type: QueryType
    product: str | None = None
    period: str | None = None


class Conversation:
    def __init__(
        self,
        message_template="<s>{role}\n{content}</s>",
        system_prompt="<s>bot\n",
        response_template="Ты — Сайга, русскоязычный автоматический ассистент. Ты разговариваешь с людьми и помогаешь им.",
    ):
        self.message_template = message_template
        self.response_template = response_template
        self.system_prompt = system_prompt
        self.messages = [{"role": "system", "content": self.system_prompt}]

    def add_user_message(self, message):
        self.messages.append({"role": "user", "content": message})

    def add_bot_message(self, message):
        self.messages.append({"role": "bot", "content": message})

    def get_prompt(self) -> list[dict[str, str]]:
        self.messages.append({"role": "system", "content": self.system_prompt})
        return self.messages


class SaigaPrompter:
    def __init__(
        self,
        prompts_path: str = "./prompts.json",
    ) -> SaigaOutput:
        self.client = ollama.Client(os.getenv("OLLAMA_HOST"))
        self.prompts = json.load(open(prompts_path))

    def generate(self, messages: list[dict[str, str]]):
        print(messages)
        output = self.client.chat(model="saiga", messages=messages, stream=False)
        return output["message"]["content"]

    def generate_response(self, prompt: str):
        conversation = Conversation()
        conversation.add_user_message(prompt)
        messages = conversation.get_prompt(self.tokenizer)
        response = self.generate(messages)
        return response

    def process_request(self, request: str):

        saiga_output = SaigaOutput(type=QueryType.UNDEFINED)

        inp = self.prompts["classifier"].format(request=request)

        classifier_outp = self.generate_response(inp)

        if "склад" in classifier_outp.lower():
            saiga_output.type = QueryType.STOCK

        elif "закупки" in classifier_outp.lower():
            saiga_output.type = QueryType.PREDICTION

        else:
            saiga_output.type = QueryType.UNDEFINED

        if saiga_output.type != QueryType.UNDEFINED:
            inp = self.prompts["product_extractor"].format(request=request)
            outp = self.generate_response(inp)
            if not ("Название продукта" in outp and "Период прогнозирования" in outp):
                saiga_output.type = QueryType.UNDEFINED
            else:
                saiga_output.product = outp.split("Название продукта:")[1].split("\n")[
                    0
                ]
                saiga_output.period = outp.split("Период прогнозирования:")[1].split(
                    "\n"
                )[0]

        if saiga_output.type == QueryType.PREDICTION:

            inp = self.prompts["time_normalizer"].format(request=saiga_output.period)
            outp = self.generate_response(inp)
            if "Период (в месяцах):" not in outp:
                saiga_output.type = QueryType.UNDEFINED
            else:
                saiga_output.period = outp.split("Период (в месяцах): ")[1].split("\n")[
                    0
                ]
                if not saiga_output.period.isdigit():
                    saiga_output.type = QueryType.UNDEFINED

        if saiga_output.type == QueryType.STOCK:
            saiga_output.period = None
        return saiga_output


if __name__ == "__main__":
    saiga = SaigaPrompter()
    print(saiga.process_request("Привет, построй план закупки чайников на 1 год."))
