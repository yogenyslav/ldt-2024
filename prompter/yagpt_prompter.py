import json
import os
import requests
import time
from dataclasses import dataclass
from dotenv import load_dotenv
from enum import Enum

from api.prompter_pb2 import QueryType

load_dotenv(".env")


@dataclass
class PrompterOutput:
    type: QueryType
    product: str | None = None
    period: str | None = None


class PromptType(Enum):
    CLASSFIER = "classifier"
    PRODUCT_EXTRACTOR = "product_extractor"
    TIME_NORMALIZER = "time_normalizer"
    FINAL_PREDICTION_PART1 = "final_prediction_part1"
    FINAL_PREDICTION_PART2 = "final_prediction_part2"


class YaGPTPrompter:
    def __init__(
        self,
        prompts_path: str = "./prompts.json",
    ) -> PrompterOutput:

        self._yandex_api_key = os.getenv("YANDEX_API_KEY")
        self._yandex_folder_id = os.getenv("YANDEX_FOLDER_ID")
        self._prompts = json.load(open(prompts_path))
        self._url = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"

        self._headers = {
            "Content-Type": "application/json",
            "Authorization": f"Api-Key {self._yandex_api_key}",
        }

        model_choice = os.getenv("MODEL_CHOICE")
        self._model_uri = f"gpt://{self._yandex_folder_id}/{model_choice}/latest"

    def _prepare_prompt(
        self, prompt: str, request_type: PromptType, stream: True
    ) -> dict:
        return {
            "modelUri": self._model_uri,
            "completionOptions": {
                "stream": stream,
                "temperature": 0.3,
                "maxTokens": "5000",
            },
            "messages": [
                {"role": "system", "text": "Ты — умный ассистент."},
                {"role": "user", "text": prompt},
            ],
        }

    def prepare_prompt1(self, data: dict) -> str:

        prompt = ""
        prompt += f"Код категории: {data['code']}\n"
        prompt += f"Наименование категории: {data['code_name']}\n"
        prompt += f"Регулярная/нерегулярная: {'регулярная' if data['is_regular'] else 'нерегулярная'}\n"
        prompt += f"Прогноз закупок:\n"
        for i, forecast in enumerate(data["forecast"], start=1):
            prompt += f"Закупка № {i}\n"
            prompt += f"Рекомендуемая дата заключения: {forecast['date']}\n"
            prompt += f"Рекомендуемая сумма закупки: {forecast['value']}\n"

        prompt += f"Медианное время выполнения закупки по категории: {data['median_execution_days']}\n"
        prompt += f"Среднее время до начала выполнения закупки по категории: {data['mean_start_to_execute_days']}\n"
        prompt += f"Средняя референсная цена: {data['mean_ref_price']}\n"

        prompt += "Топ 5 поставщиков этой категории по объему закупок:\n"
        for i, seller in enumerate(data["top5_providers"], start=1):
            prompt += f"Поставщик {i}, код исполнителя: {seller}\n"

        prompt += f"\n Выгрзука из файла для оформления закупок: \n"
        if data["closest_purchase"]["volume"]:
            prompt += f"Объем закупки (условные единицы): {data['closest_purchase']['volume']}\n"
        for i, delivery in enumerate(data["output_json"]["rows"], start=1):
            prompt += f"Позиция {i}:\n"
            prompt += (
                f"Дата начала поставки: {delivery['DeliverySchedule']['start_date']}\n"
            )
            prompt += (
                f"Дата окончания поставки: {delivery['DeliverySchedule']['end_date']}\n"
            )
            if delivery["DeliverySchedule"]["deliveryAmount"]:
                prompt += f"Объем поставки (условных единиц): {delivery['DeliverySchedule']['deliveryAmount']}\n"
            prompt += f"Номер версии: {delivery['id']}\n"
            prompt += f"Объем в рублях: {delivery['nmc']}\n"
            prompt += f"ID СПГЗ: {delivery['spgzCharacteristics']['spgzId']}"
            prompt += (
                f"Наименование СПГЗ: {delivery['spgzCharacteristics']['spgzName']}\n"
            )
            prompt += f"Код КПГЗ: {delivery['spgzCharacteristics']['kpgzCode']}\n"
            prompt += (
                f"Наименование КПГЗ: {delivery['spgzCharacteristics']['kpgzName']}\n"
            )
            if i == 4:
                break

        return prompt

    def prepare_prompt2(self, data: dict) -> str:
        prompt = ""
        for i, deal in enumerate(data["example_contracts_in_code"], start=1):
            prompt += f"Контракт № {i}:\n"
            prompt += f"Код СПГЗ: {deal['id_spgz']}\n"
            prompt += f"Конечное наименование КПГЗ: {deal['name_spgz']}\n"
            prompt += f"Наименование ГК: {deal['item_name_gk']}\n"
            prompt += f"Дата регистрации контракта: {deal['conclusion_date']}\n"
            prompt += (
                f"Дата начала выполнения контракта: {deal['execution_term_from']}\n"
            )
            prompt += (
                f"Дата окончания выполнения контракта: {deal['execution_term_until']}\n"
            )
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

    def prepare_prompt_with_nan(self, data: dict) -> str:
        prompt = ""
        prompt += f"Код категории: {data['code']}\n"
        prompt += f"Наименование категории: {data['code_name']}\n"
        return prompt

    def prepare_prompt_with_empty_forecast(self, data: dict) -> str:
        prompt = ""
        prompt += f"Код категории: {data['code']}\n"
        prompt += f"Наименование категории: {data['code_name']}\n"
        prompt += f"Регулярная/нерегулярная: {'регулярная' if data['is_regular'] else 'нерегулярная'}\n"
        prompt += f"Напиши о том, что период прогнозирования выбран слишком мал, поэтому ниже представлена информация о ближайшей закупке\n"
        prompt += f"Дата заключения контракта: {data['closest_purchase']['date']}\n"
        prompt += f"Сумма закупки: {round(data['closest_purchase']['value'], 2)}\n"
        prompt += f"Объем закупки (условные единицы): {data['closest_purchase']['volume'] if data['closest_purchase']['volume'] else 'Информация отсутствует'}\n"
        prompt += f"Медианное время выполнения закупки по категории: {data['median_execution_days']}\n"
        prompt += f"Среднее время до начала выполнения закупки по категории: {data['mean_start_to_execute_days']}\n"
        prompt += f"Средняя референсная цена: {data['mean_ref_price'] if data['mean_ref_price'] else 'Информация отсутствует'}\n"

        prompt += "Топ 5 поставщиков этой категории по объему закупок:\n"
        for i, seller in enumerate(data["top5_providers"], start=1):
            prompt += f"Поставщик {i}, код исполнителя: {seller}\n"

        prompt += f"Подробная информация о ближайшей закупке:\n"
        prompt += f"Дата заключения контракта: {data['closest_purchase']['date']}\n"
        prompt += f"Сумма закупки: {data['closest_purchase']['value']}\n"
        if data["closest_purchase"]["volume"]:
            prompt += f"Объем закупки (условные единицы): {data['closest_purchase']['volume']}\n"

        return prompt

    def _generate_responce(
        self, prompt: str, request_type: PromptType, stream: bool = False
    ):
        time.sleep(1)  # TODO: sorry for this abomination, I'll fix this later
        prepared_prompt = self._prepare_prompt(prompt, request_type, stream=stream)
        if not stream:
            retry_attempts = 10
            while retry_attempts >= 0:
                try:
                    response = requests.post(
                        self._url,
                        headers=self._headers,
                        data=json.dumps(prepared_prompt),
                        stream=stream,
                    )
                    response.raise_for_status()
                    break
                except requests.exceptions.SSLError as e:
                    print("SSL Error, retrying...")
                    retry_attempts -= 1
                    if retry_attempts > 0:
                        time.sleep(5)
                    else:
                        raise e
                except requests.exceptions.HTTPError as e:
                    print("HTTP Error, retrying...")
                    retry_attempts -= 1
                    if retry_attempts > 0:
                        time.sleep(5)
                    else:
                        raise e

            output = json.loads(response.content.decode("utf-8"))
            output = output["result"]["alternatives"][-1]["message"]["text"]
        else:
            retry_attempts = 10
            while retry_attempts >= 0:
                try:
                    response = requests.post(
                        self._url,
                        headers=self._headers,
                        data=json.dumps(prepared_prompt),
                        stream=stream,
                    )
                    response.raise_for_status()
                    break
                except requests.exceptions.SSLError as e:
                    print("SSL Error, retrying...")
                    retry_attempts -= 1
                    if retry_attempts > 0:
                        time.sleep(5)
                    else:
                        raise e

                except requests.exceptions.HTTPError as e:
                    print("HTTP Error, retrying...")
                    retry_attempts -= 1
                    if retry_attempts > 0:
                        time.sleep(5)
                    else:
                        raise e

            output = response
        return output

    def process_request(self, request: str):

        prompter_output = PrompterOutput(type=QueryType.UNDEFINED)

        inp = self._prompts["classifier"].format(request=request)

        classifier_outp = self._generate_responce(inp, PromptType.CLASSFIER)

        if "склад" in classifier_outp.lower():
            prompter_output.type = QueryType.STOCK

        elif "закупки" in classifier_outp.lower():
            prompter_output.type = QueryType.PREDICTION

        else:
            prompter_output.type = QueryType.UNDEFINED

        if prompter_output.type != QueryType.UNDEFINED:
            inp = self._prompts["product_extractor"].format(request=request)
            outp = self._generate_responce(inp, PromptType.PRODUCT_EXTRACTOR)
            if not ("Название продукта" in outp and "Период прогнозирования" in outp):
                prompter_output.type = QueryType.UNDEFINED
            else:
                prompter_output.product = outp.split("Название продукта:")[1].split(
                    "\n"
                )[0]
                prompter_output.period = outp.split("Период прогнозирования:")[1].split(
                    "\n"
                )[0]

        if prompter_output.type == QueryType.PREDICTION:
            inp = self._prompts["time_normalizer"].format(request=request)
            outp = self._generate_responce(inp, PromptType.TIME_NORMALIZER)
            if "Период (в месяцах)" not in outp:
                prompter_output.period = outp
            else:
                prompter_output.period = (
                    outp.split("Период (в месяцах):")[1].split("\n")[0].replace(" ", "")
                )
                if not prompter_output.period.isdigit():
                    prompter_output.type = QueryType.UNDEFINED

        if prompter_output.type == QueryType.STOCK:
            prompter_output.period = None
        return prompter_output

    def process_final_request(self, data: str, prompt_type: PromptType):
        inp = json.loads(data)
        if inp["forecast"] is None:
            request = self.prepare_prompt_with_nan(inp)
            request += """ЗАПРОС: Напиши ответ о том, что пользователь товар, который не закупается регулярно, 
            для неё невозможность построить прогноз (2-3 предложения на ответ)
            Не давай лишней информации, посоветуй пользователю ввести другой товар. Не советуй пользователю ничего"""
        elif inp["forecast"] == []:
            request = self.prepare_prompt_with_empty_forecast(inp)
            request += """ЗАПРОС 1: Оформи отчет в MARKDOWN, добавь таблицы, где нужно, если они слишком широкие (столбцов больше, чем строк), поменяй строки с столбцы местами. Убери None, где нет информации. 
            В своем отчете укажи всю предоставленную информацию. Не добавляй информацию, которой нет в исходных данных. 
            Добавь суммаризацию, в которой скажи то, что информацию о закупках можно поменять в соответствии с потребностями заказчика."""
        else:
            request = self.prepare_prompt1(inp)
            request += """ЗАПРОС 1: Оформи отчет в MARKDOWN, добавь таблицы, где нужно, если они слишком широкие (столбцов больше, чем строк), поменяй строки с столбцы местами. Убери None, где нет информации. 
            В своем отчете укажи всю предоставленную информацию. Не добавляй информацию, которой нет в исходных данных. 
            Добавь суммаризацию, в которой скажи то, что информацию о закупках можно поменять в соответствии с потребностями заказчика."""
        return self._generate_responce(request, prompt_type, stream=True)
