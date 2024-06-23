import json
import os
from datetime import datetime, timedelta
from typing import Tuple

import pymongo
from data_process import filter_forecast
from utils import convert_datetime_to_str

mongo_url = f"mongodb://{os.getenv('MONGO_HOST')}:{os.getenv('MONGO_PORT')}"
mongo_client = pymongo.MongoClient(mongo_url)


def detect(organization: str, current_date: str = "2024-05-18") -> Tuple[str, str]:
    current_date = datetime.strptime(current_date, "%Y-%m-%d")

    mongo_db = mongo_client[organization]
    collection_name = "codes"
    collection = mongo_db[collection_name]

    codes_info = collection.find({"is_regular": True})
    codes_info = list(codes_info)

    all_mails_body = []

    for code_info in codes_info:
        forecast = code_info["forecast"]
        median_execution_days = code_info["median_execution_days"]
        code_name = code_info["code_name"]
        code = code_info["code"]

        current_purchase = None

        for purchase in forecast:
            purchase_date = purchase["date"]
            if (
                current_date > (purchase_date - timedelta(days=14))
            ) and current_date < purchase_date:
                current_purchase = purchase
                break

        if current_purchase is None:
            continue

        filtered_forecast = filter_forecast(
            code_info,
            organization,
            code,
            current_date,
            current_date + timedelta(days=15),
        )
        json_content = json.dumps(
            convert_datetime_to_str(filtered_forecast)["output_json"], indent=4
        )

        mail_body = f"""
                    Советуем Вам совершить покупку по категории <<{code_name} - {code}>>.
                    Предлагаемая дата заключения: {current_purchase['date']}
                    Предполагаемое время исполнения (в днях): {median_execution_days}
                    Сумма закупки: {current_purchase['value']} рублей
                    Объем закупки: {current_purchase['volume']} штук
                    
                    Файл с сформированной закупкой (в разрезе товаров) приложен к письму.
                    """

        all_mails_body.append((mail_body, json_content))

    return all_mails_body


if __name__ == "__main__":
    print(detect("test123"))
