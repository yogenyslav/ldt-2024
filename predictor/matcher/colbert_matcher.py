import json
from typing import List, Optional

from colbert import  Searcher
from colbert.infra import ColBERTConfig, Run, RunConfig


class ColbertMatcher:

    def __init__(
        self,
        checkpoint_name: str,
        collection_path: str,
        category2code_path: Optional[str] = None,
    ) -> None:
        """
        Initializes a ColbertMatcher object.

        Args:
            checkpoint_name (str): The name of the ColBERT checkpoint.
            collection_path (str): The path to the collection JSON file.
            category2code_path (Optional[str], optional): The path to the category to code
                JSON file. Defaults to None.
        """

        with open(collection_path, "r") as read_file:
            collection = json.load(read_file)

        self.config = ColBERTConfig(
            doc_maxlen=400,
            nbits=8,
            kmeans_niters=10,
            ncells=5,
            centroid_score_threshold=0.7,
            ndocs=512,
            root="/app/matcher/experiments",
        )

        with Run().context(
            RunConfig(experiment="searcher", root="/app/matcher/experiments")
        ):
            self.searcher = Searcher(index=checkpoint_name, collection=collection)

        if category2code_path:
            with open(category2code_path, "r") as read_file:
                self.category2code = json.load(read_file)

    def match_code(self, query: str) -> str:
        """
        Match a query to a code and return the corresponding code.

        Args:
            query (str): The query to match.

        Returns:
            str: The corresponding code.
        """
        results = self.searcher.search(query, k=1)
        return self.category2code[self.searcher.collection[results[0][0]]]

    def match_stocks(self, query: str, top_k: int = 5) -> List[str]:
        """
        Match a query to stocks and return the top K stocks.

        Args:
            query (str): The query to match.
            top_k (int): The number of top stocks to return. Defaults to 5.

        Returns:
            List[str]: The top K stocks matching the query.
        """
        results = self.searcher.search(query, k=top_k)
        return [self.searcher.collection[x] for x in results[0]]
