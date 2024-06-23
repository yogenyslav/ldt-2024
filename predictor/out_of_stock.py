import json
import os
import smtplib
import asyncpg
import asyncio
from email.message import EmailMessage
from datetime import datetime, timedelta
from typing import Tuple

import pymongo
from data_process import filter_forecast
from utils import convert_datetime_to_str
from dotenv import load_dotenv

load_dotenv(".env")

mongo_url = f"mongodb://{os.getenv('MONGO_HOST')}:{os.getenv('MONGO_PORT')}"
mongo_client = pymongo.MongoClient(mongo_url)
postgres_url = os.getenv(
    f"postgres+asyncpg://{os.getenv('POSTGRES_USER')}:{os.getenv('POSTGRES_PASSWORD')}@{os.getenv('POSTGRES_HOST')}:{os.getenv('POSTGRES_PORT')}/{os.getenv('POSTGRES_DB')}"
)


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


def send_email(
    server: smtplib.SMTP_SSL, data: tuple[str, str, str], body: str, content: str
):
    msg = EmailMessage()
    msg["Subject"] = "Предложение по закупке"
    msg["From"] = os.getenv("MAIL_USERNAME")
    msg["To"] = data[0]
    body = f"Здравствуйте, {data[1]} {data[2]}!\n\n{body}"
    msg.set_content(body)
    msg.add_attachment(
        content, filename="purchase.json", maintype="application", subtype="json"
    )
    server.send_message(msg)


async def get_user_data(
    conn: asyncpg.pool.PoolAcquireContext, organization: str
) -> list[tuple[str, str, str]]:
    org_id = organization.split("-")[1]
    res = await conn.fetch(
        "select email, first_name, last_name from chat.notification where organization_id = $1",
        org_id,
    )
    print(res)
    return [(r["email"], r["first_name"], r["last_name"]) for r in res]


async def main():
    with smtplib.SMTP_SSL("smtp.yandex.com", 465, timeout=60) as server:
        organizations = mongo_client.list_database_names()
        server.login(
            os.getenv("MAIL_USERNAME"),
            os.getenv("MAIL_PASSWORD"),
        )
        async with asyncpg.create_pool(postgres_url) as pool:
            async with pool.acquire() as conn:
                for organization in organizations:
                    user_data = await get_user_data(conn, organization)
                    for data in user_data:
                        mails_body = detect(organization)
                        for mail_body, content in mails_body:
                            send_email(server, data, mail_body, content)


if __name__ == "__main__":
    asyncio.run(main())
