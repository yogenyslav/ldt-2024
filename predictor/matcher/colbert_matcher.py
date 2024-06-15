import json
import colbert
from colbert import Indexer, Searcher
from colbert.infra import Run, RunConfig, ColBERTConfig
from colbert.data import Queries, Collection

class ColbertMatcher:

    def __init__(
            self,
            checkpoint_name: str,
            collection_path: str,
            category2code_path: str | None = None,
    ) -> None:
        
        with open(collection_path, "r") as read_file:
            collection = json.load(read_file)

        self.config = ColBERTConfig(
            doc_maxlen=400,
            nbits=8,
            kmeans_niters=10,
            ncells=5,
            centroid_score_threshold=0.7,
            ndocs=512,
            root='/app/matcher/experiments'
        )

        with Run().context(RunConfig(experiment="searcher", root='/app/matcher/experiments')):
            self.searcher = Searcher(index=checkpoint_name, collection=collection)
        
        if category2code_path:
            with open(category2code_path, "r") as read_file:
                self.category2code = json.load(read_file)

    def match_to_3rd_level_code(self, query: str) -> str:
        results = self.searcher.search(query, k=1)
        return self.category2code[self.searcher.collection[results[0][0]]]
    
    def match_to_full_name(self, query: str) -> str:
        results = self.searcher.search(query, k=5)
        return [self.searcher.collection[x] for x in results[0]]
