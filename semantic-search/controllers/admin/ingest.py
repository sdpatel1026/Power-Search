import os

from langchain.document_loaders.pdf import OnlinePDFLoader
from langchain.embeddings import OpenAIEmbeddings
from langchain.embeddings.base import Embeddings
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.vectorstores.pinecone import Pinecone
from pydantic import BaseModel
from fastapi import APIRouter

router = APIRouter()


class Document(BaseModel):
    url: str
    doc_id: int
    published_year: int


def load_documents(url: str):
    loader = OnlinePDFLoader(url)
    return loader.load()


def chunk_datas(documents, chunk_size=1000, overlap_size=40):
    text_splitter = RecursiveCharacterTextSplitter(chunk_size=chunk_size, chunk_overlap=overlap_size)
    chunks = text_splitter.split_documents(documents)
    return chunks


def store_embeddings(chunks: list[Document], embeddings: Embeddings, index_name: str, doc_url: str, doc_id: int,
                     publish_year, chunk_size=1000):
    texts = [c.page_content for c in chunks]
    metadatas = [{"url": doc_url, "doc_id": doc_id} for _ in chunks]
    Pinecone.from_texts(texts, embeddings, metadatas=metadatas, index_name=index_name, embeddings_chunk_size=chunk_size,
                        namespace=publish_year)


@router.post("/ingest")
def ingest(doc: Document):
    documents = load_documents(doc.url)
    chunks = chunk_datas(documents)
    store_embeddings(chunks, OpenAIEmbeddings(openai_api_key=os.getenv("OPENAI_API_KEY_EMBEDDING")),
                     os.getenv("PINECONE_INDEX_NAME"), doc.url, doc.doc_id, doc.published_year)
