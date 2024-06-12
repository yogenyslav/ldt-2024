import numpy as np
import pandas as pd
from scipy.optimize import linear_sum_assignment
from scipy.spatial.distance import cdist
from sentence_transformers import SentenceTransformer

contracts_map_columns = {
    "Дата заключения": "conclusion_date",
    "Срок исполнения по": "execution_term_until",
    "Срок исполнения с": "execution_term_from",
    "Дата окончания срока действия": "end_date_of_validity",
    "Оплачено, руб.": "paid_rub",
    "Цена ГК, руб.": "gk_price_rub",
    "ID СПГЗ": "id_spgz",
    "Конечный код КПГЗ": "final_code_kpgz",
    "Наименование СПГЗ": "name_spgz",
    "Наименование (предмет) ГК": "item_name_gk",
    "Конечное наименование КПГЗ": "final_name_kpgz",
    "Реестровый номер в РК": "registry_number_in_rk",
    "Исполнено поставщиком": "provider",
}


kpgz_map_columns = {
    "Название СТЕ": "name_ste",
    "наименование характеристик": "characteristics_name",
    "реф.цена": "ref_price",
    "Реестровый номер в РК": "registry_number_in_rk",
}


def prepare_contracts_df(df: pd.DataFrame) -> pd.DataFrame:
    df = df[contracts_map_columns.keys()]
    df = df.rename(contracts_map_columns, axis=1)
    df = df.dropna(axis=0, subset=["final_code_kpgz", "conclusion_date"])

    df["conclusion_date"] = pd.to_datetime(df["conclusion_date"], dayfirst=True)
    df["execution_term_until"] = pd.to_datetime(
        df["execution_term_until"], dayfirst=True
    )
    df["execution_term_from"] = pd.to_datetime(df["execution_term_from"], dayfirst=True)
    df["end_date_of_validity"] = pd.to_datetime(
        df["end_date_of_validity"], dayfirst=True
    )

    df["depth3_code_kpgz"] = df["final_code_kpgz"].str.slice(0, 8).astype(str)

    df["registry_number_in_rk"] = df["registry_number_in_rk"].astype(pd.StringDtype())
    df = df.set_index("registry_number_in_rk", drop=False)

    df["id"] = np.arange(len(df))
    df["id"] = df["id"].convert_dtypes()

    return df


def prepare_kpgz_df(df: pd.DataFrame) -> pd.DataFrame:
    df = df[kpgz_map_columns.keys()]
    df = df.rename(kpgz_map_columns, axis=1)
    df = df.dropna(axis=0, subset=["name_ste"])

    df["registry_number_in_rk"] = df["registry_number_in_rk"].astype(pd.StringDtype())
    df = df.set_index("registry_number_in_rk", drop=True)

    df["id"] = np.arange(len(df))
    df["id"] = df["id"].convert_dtypes()

    return df


def merge_contracts_with_kpgz(contracts_df, kpgz_df):
    model_name = "sentence-transformers/distilbert-multilingual-nli-stsb-quora-ranking"
    encoder = SentenceTransformer(model_name)

    contracts_names = contracts_df["name_spgz"].astype(pd.StringDtype())
    kpgz_names = kpgz_df["name_ste"].astype(pd.StringDtype())

    contracts_embedding = encoder.encode(contracts_names.tolist())
    kpgz_embedding = encoder.encode(kpgz_names.tolist())

    uniq_keys1 = contracts_df.index.values.unique().to_numpy()
    uniq_keys2 = kpgz_df.index.values.unique().to_numpy()

    uniq_keys1 = uniq_keys1[~pd.isna(uniq_keys1)]
    uniq_keys2 = uniq_keys2[~pd.isna(uniq_keys2)]

    uniq_keys = np.intersect1d(uniq_keys1, uniq_keys2, assume_unique=True)
    rest_keys_contracts = np.array(list(set(uniq_keys1) - set(uniq_keys)), dtype=object)

    to_concat_rows = []

    for key in uniq_keys:
        row1 = contracts_df.loc[key]
        row2 = kpgz_df.loc[key]

        if not isinstance(row1, pd.DataFrame):
            row1 = pd.DataFrame(row1).T
        if not isinstance(row2, pd.DataFrame):
            row2 = pd.DataFrame(row2).T

        embed_ids1 = row1["id"].to_numpy(dtype=int)
        embed_ids2 = row2["id"].to_numpy(dtype=int)

        embeds1 = contracts_embedding[embed_ids1]
        embeds2 = kpgz_embedding[embed_ids2]

        if embeds1.shape[0] == 1 and embeds2.shape[0] == 1:
            row_ind = np.array(
                [
                    0,
                ]
            )
            col_ind = np.array(
                [
                    0,
                ]
            )
        else:
            distance_matrix = cdist(embeds1, embeds2, metric="cosine")
            row_ind, col_ind = linear_sum_assignment(distance_matrix)

        matched_row1 = row1.iloc[row_ind]
        matched_row2 = row2.iloc[col_ind]

        concated_row = pd.concat(
            [matched_row1, matched_row2.drop(columns=["id"])], axis=1
        )

        to_concat_rows.append(concated_row)

    merged_df = pd.concat(to_concat_rows, axis=0)
    merged_df = pd.concat(
        [merged_df, contracts_df.loc[rest_keys_contracts]], axis=0
    ).convert_dtypes()
    merged_df = merged_df.sort_index()

    return merged_df
