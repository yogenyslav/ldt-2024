import ollama
import json
import os
from dataclasses import dataclass
from api.prompter_pb2 import QueryType
from dotenv import load_dotenv
from enum import Enum

load_dotenv(".env")


@dataclass
class SaigaOutput:
    type: QueryType
    product: str | None = None
    period: str | None = None

class PromptType(Enum):
    CLASSFIER = "classifier"
    PRODUCT_EXTRACTOR = "product_extractor"
    TIME_NORMALIZER = "time_normalizer"
    FINAL_PREDICTION_PART1 = "final_prediction_part1"
    FINAL_PREDICTION_PART2 = "final_prediction_part2"


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

    def get_prompt(self):
        final_text = ""
        for message in self.messages:
            message_text = self.message_template.format(**message)
            final_text += message_text
        final_text += "<s>bot\n"
        return final_text.strip()


class SaigaPrompter:
    def __init__(
        self,
        prompts_path: str = "./prompts.json",
    ) -> SaigaOutput:
        host = os.getenv("OLLAMA_HOST")
        self.client = ollama.Client(host=host)
        self.prompts = json.load(open(prompts_path))

    def generate(self, prompt: str):
        output = self.client.generate(model="saiga", prompt=prompt, stream=False)
        return output["response"]

    def generate(self, prompt: str, temp: float = 0.8):
        output = self.client.chat(
            model="saiga", 
            messages=[{"role": "user", "content": prompt}],
            stream=False,
            options={"temperature": temp, "num_predict": -2, "top_k": 1})
        return output["message"]["content"]
    
    def generate_stream(self, prompt: str, prompt_type: PromptType):
        if prompt_type == PromptType.FINAL_PREDICTION_PART1:
            temp = 0.9
            top_k = 40
        elif prompt_type == PromptType.FINAL_PREDICTION_PART2:
            temp = 0.85
            top_k = 1
        output = self.client.chat(
                model="saiga",
                messages=[{"role": "user", "content": prompt}],
                stream=True,
                options={"temperature": temp, "num_predict": -1, "top_k": top_k},
            )
        return output
    
    def generate_response(self, prompt: str, request_type: PromptType):
        conversation = Conversation()
        conversation.add_user_message(prompt)
        prompt = conversation.get_prompt()
        temp = 0.92 if request_type == PromptType.TIME_NORMALIZER else 0.8
        response = self.generate(prompt, temp)
        return response
    
    def generate_response_stream(self, prompt: str, prompt_type: PromptType):
        conversation = Conversation()
        conversation.add_user_message(prompt)
        prompt = conversation.get_prompt()
        response = self.generate_stream(prompt, prompt_type)
        return response
    
    def prepare_prompt1(self, data: dict) -> str:

        prompt = ""
        prompt += f"Код категории: {data['code']}\n"
        prompt += f"Наименование категории: {data['code_name']}\n"
        prompt += f"Регулярная/нерегулярная: {'регулярная' if data['is_regular'] else 'нерегулярная'}\n"
        prompt += f"Прогноз закупок:\n"
        for i, forecast in enumerate(data['forecast'], start=1):
            prompt += f"Закупка № {i}\n"
            prompt += f"Рекомендуемая дата заключения: {forecast['date']}\n"
            prompt += f"Рекомендуемая сумма закупки: {forecast['value']}\n"

        prompt += f"Медианное время выполнения закупки по категории: {data['median_execution_days']}\n"
        prompt += f"Среднее время до начала выполнения закупки по категории: {data['mean_start_to_execute_days']}\n"
        prompt += f"Средняя референсная цена: {data['mean_ref_price']}\n"

        prompt += "Топ 5 поставщиков этой категории по объему закупок:\n"
        for i, seller in enumerate(data['top5_providers'], start=1):
            prompt += f"Поставщик {i}, код исполнителя: {seller}\n"
        return prompt
    
    def prepare_prompt2(self, data: dict) -> str:
        prompt = ""
        for i, deal in enumerate(data['contracts_in_code'], start=1):
            prompt += f"Контракт № {i}:\n"
            prompt += f"Код СПГЗ: {deal['id_spgz']}\n"
            prompt += f"Конечное наименование КПГЗ: {deal['name_spgz']}\n"
            prompt += f"Наименование ГК: {deal['item_name_gk']}\n"
            prompt += f"Дата регистрации контракта: {deal['conclusion_date']}\n"
            prompt += f"Дата начала выполнения контракта: {deal['execution_term_from']}\n"
            prompt += f"Дата окончания выполнения контракта: {deal['execution_term_until']}\n"
            prompt += f"Дата окончания срока действия: {deal['end_date_of_validity']}\n"
            prompt += f"Оплачено, руб.: {deal['paid_rub']}\n"
            prompt += f"Цена ГК при заключении, руб.: {deal['gk_price_rub']}\n"
            prompt += f"Конечный код КПГЗ: {deal['final_code_kpgz']}\n"
            prompt += f"Конечное наименование КПГЗ: {deal['final_name_kpgz']}\n"
            prompt += f"Реестровый номер в РК: {deal['registry_number_in_rk']}\n"
            prompt += f"Код поставщика: {deal['provider']}\n"
            prompt += f"Референсная цена: {deal['ref_price']}\n"
            if i == 4:
                break

        return prompt
    
    def process_final_request(self, data: str, prompt_type: PromptType):
        inp = json.loads(data)
        if prompt_type == PromptType.FINAL_PREDICTION_PART1:
            request = self.prepare_prompt1(inp)
            request += "ЗАПРОС: Оформи отчет в MARKDOWN. Убери None, где нет информации. В своем отчете укажи всю предоставленную информацию. Не добавляй информацию, которой нет в исходных данных."
        if prompt_type == PromptType.FINAL_PREDICTION_PART2:
            request = self.prepare_prompt1(inp)
            request += "ЗАПРОС: Оформи эту информацию в MARKDOWN таблице"
        return self.generate_response_stream(request, prompt_type)

    def process_request(self, request: str):

        saiga_output = SaigaOutput(type=QueryType.UNDEFINED)

        inp = self.prompts["classifier"].format(request=request)

        classifier_outp = self.generate_response(inp, PromptType.CLASSFIER)

        if "склад" in classifier_outp.lower():
            saiga_output.type = QueryType.STOCK

        elif "закупки" in classifier_outp.lower():
            saiga_output.type = QueryType.PREDICTION

        else:
            saiga_output.type = QueryType.UNDEFINED

        if saiga_output.type != QueryType.UNDEFINED:
            inp = self.prompts["product_extractor"].format(request=request)
            outp = self.generate_response(inp, PromptType.PRODUCT_EXTRACTOR)
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
            outp = self.generate_response(inp, PromptType.TIME_NORMALIZER)
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
