import numpy as np
import requests
from scipy.spatial.distance import cdist


class YaMatcher:
    def __init__(
        self,
        path2embeds_dir: str,
        folder_id: str,
        api_key: str,
    ) -> None:
        """
        Initialize the YaMatcher class.

        Args:
            path2embeds_dir (str): Path to the directory containing the embeddings.
            folder_id (str): Yandex.Cloud folder ID.
            api_key (str): Yandex.Cloud API key.

        Returns:
            None
        """
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

    def get_embedding(self, text: str, text_type: str = "query") -> np.ndarray:
        """
        Get the embedding of a given text.

        Args:
            text (str): The text to get the embedding for.
            text_type (str, optional): The type of the text. Defaults to "query".

        Returns:
            np.ndarray: The embedding of the text.
        """
        query_data = {
            "modelUri": self._doc_uri if text_type == "doc" else self._query_uri,
            "text": text,
        }
        req = requests.post(self._embed_url, json=query_data, headers=self._headers)

        return np.array(req.json()["embedding"])

    def match_stocks(self, query: str, top_k: int = 5) -> np.ndarray:
        """
        Match stocks based on query.

        Args:
            query (str): The query string for stock names.
            top_k (int): The number of top stocks to return. Default is 5.

        Returns:
            List[str]: A list of top stock names that match the query.
        """
        query_embedding = self.get_embedding(query, text_type="query")
        dist = cdist(query_embedding[None, :], self._stocks_embeds, metric="cosine")
        return self._stocks_names[np.argsort(dist)[0, :top_k]]

    def match_code(self, query: str) -> str:
        """
        Match a code based on query.

        Args:
            query (str): The query string for code names.

        Returns:
            str: The code name that best matches the query.
        """
        query_embedding = self.get_embedding(query, text_type="query")
        dist = cdist(query_embedding[None, :], self._code_embeds, metric="cosine")
        return self._code_names[np.argmin(dist[0])]
