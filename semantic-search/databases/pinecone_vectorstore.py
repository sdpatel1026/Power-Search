import threading
import pinecone
from pinecone import UnauthorizedException


def list_indexes():
    return pinecone.list_indexes()


class PineconeClient:
    _instances = {}
    _lock = threading.Lock()

    def __new__(cls, api_key, environment, index_name):
        if index_name not in cls._instances:
            with cls._lock:
                if index_name not in cls._instances:
                    instance = super(PineconeClient, cls).__new__(cls)
                    cls._instances[index_name] = instance
                    instance.api_key = api_key
                    instance.index_name = index_name
                    if api_key is None or environment is None:
                        raise ValueError("PineCone api_key or environment not found")
                    try:
                        pinecone.init(api_key=api_key, environment=environment)
                    except UnauthorizedException:
                        raise ValueError("Invalid pineCone api_key or environment")
                    instance.index = pinecone.Index(index_name)
        return cls._instances[index_name]

    def create_index(self, dimension, metric_type="cosine"):
        # self.index.create(dimension=dimension, metric_type=metric_type)
        pass

    def delete_index(self):
        pass
        # self.index.delete()

    def insert_vectors(self, ids, vectors):
        # self.index.upsert(ids=ids, embeddings=vectors)
        pass

    def query_vectors(self, query_vector, top_k=100):
        results = self.index.query(queries=[query_vector], top_k=top_k)
        return results[0]  # Assuming you are querying a single vector

    def delete_vectors(self, ids):
        self.index.delete(ids=ids)
