import os
from langchain.embeddings import OpenAIEmbeddings
from databases.pinecone_vectorstore import PineconeClient
from fastapi import APIRouter

router = APIRouter()


@router.get("/discovery")
def discover(query: str):
    model = OpenAIEmbeddings(openai_api_key=os.getenv("OPENAI_API_KEY"))
    vectors = model.embed_query(query)
    vector_store = pinecone_client = PineconeClient(os.getenv("PINECONE_API_KEY"), os.getenv("PINECONE_ENVIRONMENT"),
                                                    os.getenv("PINECONE_INDEX_NAME"))
    result = vector_store.index.query(vectors, include_metadata=True)
    return result
