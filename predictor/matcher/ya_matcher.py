import os

import numpy as np
import requests
from scipy.spatial.distance import cdist


class YaMatcher:
    def __init__(self, path2embeds_dir: str, folder_id: str, api_key: str):
        self._stocks_embeds = np.load(f"{path2embeds_dir}/stocks_embeds.npy")
        self._stocks_names = np.load(f"{path2embeds_dir}/stocks_names.npy")

        self._code_embeds = np.load(f"{path2embeds_dir}/code_embeds.npy")
        self._code_names = np.load(f"{path2embeds_dir}/code_kpgz.npy")

        self._doc_uri = f"emb://{folder_id}/text-search-doc/latest"
        self._query_uri = f"emb://{folder_id}/text-search-query/latest"

        self._embed_url = (
            "https://llm.api.cloud.yandex.net/foundationModels/v1/textEmbedding"
        )
        self._headers = {
            "Content-Type": "application/json",
            "Authorization": f"Api-Key {api_key}",
            "x-folder-id": f"{folder_id}",
        }

    def get_embedding(self, text: str, text_type: str = "query") -> np.array:
        query_data = {
            "modelUri": self._doc_uri if text_type == "doc" else self._query_uri,
            "text": text,
        }
        req = requests.post(self._embed_url, json=query_data, headers=self._headers)

        return np.array(req.json()["embedding"])

    def match_stocks(self, query: str, top_k: int = 5):
        query_embedding = self.get_embedding(query, text_type="query")
        dist = cdist(query_embedding[None, :], self._stocks_embeds, metric="cosine")
        return self._stocks_names[np.argsort(dist)[0, :top_k]]

    def match_code(self, query: str):
        query_embedding = self.get_embedding(query, text_type="query")
        dist = cdist(query_embedding[None, :], self._code_embeds, metric="cosine")
        return self._code_names[np.argmin(dist[0])]
