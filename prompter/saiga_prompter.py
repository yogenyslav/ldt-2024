import enum
import json
from dataclasses import dataclass
import torch
from peft import PeftConfig, PeftModel
from transformers import AutoModelForCausalLM, AutoTokenizer, GenerationConfig
from api.prompter_pb2 import QueryType

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
        self.messages = [{
            "role": "system",
            "content": system_prompt
        }]

    def add_user_message(self, message):
        self.messages.append({
            "role": "user",
            "content": message
        })

    def add_bot_message(self, message):
        self.messages.append({
            "role": "bot",
            "content": message
        })

    def get_prompt(self, tokenizer):
        final_text = ""
        for message in self.messages:
            message_text = self.message_template.format(**message)
            final_text += message_text
        final_text += "<s>bot\n"
        return final_text.strip()

class SaigaPrompter:
    def __init__(
            self,
            model_name: str = "IlyaGusev/saiga_7b_lora",
            prompts_path: str = "./prompts.json",
    ) -> SaigaOutput:
        self.config = PeftConfig.from_pretrained(model_name)
        self.model = AutoModelForCausalLM.from_pretrained(
            self.config.base_model_name_or_path,
            load_in_8bit=True,
            torch_dtype=torch.float16,
            device_map="auto",
        )
        self.model = PeftModel.from_pretrained(
            self.model,
            model_name,
            torch_dtype=torch.float16,
        )

        self.model.eval()

        self.tokenizer = AutoTokenizer.from_pretrained(model_name, use_fast=True)
        self.generation_config = GenerationConfig.from_pretrained(model_name)
        self.prompts = json.load(open(prompts_path))

    def generate(self, prompt: str):
        data = self.tokenizer(prompt, return_tensors="pt", add_special_tokens=False)
        data = {k: v.to(self.model.device) for k, v in data.items()}
        output_ids = self.model.generate(
            **data,
            generation_config=self.generation_config
        )[0]
        output_ids = output_ids[len(data["input_ids"][0]):]
        output = self.tokenizer.decode(output_ids, skip_special_tokens=True)
        return output.strip()

    def generate_response(self, prompt: str):
        conversation = Conversation()
        conversation.add_user_message(prompt)
        prompt = conversation.get_prompt(self.tokenizer)
        response = self.generate(prompt)
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
            if not("Название продукта" in outp and "Период прогнозирования" in outp):
                saiga_output.type = QueryType.UNDEFINED
            else:
                saiga_output.product = outp.split("Название продукта:")[1].split("\n")[0]
                saiga_output.period = outp.split("Период прогнозирования:")[1].split("\n")[0]
                
        if saiga_output.type == QueryType.PREDICTION:

            inp = self.prompts["time_normalizer"].format(request=saiga_output.period)
            outp = self.generate_response(inp)
            if "Период (в месяцах):" not in outp:
                saiga_output.type = QueryType.UNDEFINED
            else:
                saiga_output.period = outp.split("Период (в месяцах): ")[1].split("\n")[0]
                if not saiga_output.period.isdigit():
                    saiga_output.type = QueryType.UNDEFINED
        
        if saiga_output.type == QueryType.STOCK:
            saiga_output.period = None
        return saiga_output