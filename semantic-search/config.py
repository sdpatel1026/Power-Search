import enum
import os

from dotenv import load_dotenv

from databases.pinecone_vectorstore import PineconeClient
from pydantic_settings import BaseSettings


class LogLevel(str, enum.Enum):  # noqa: WPS600
    """Possible log levels."""

    NOTSET = "NOTSET"
    DEBUG = "DEBUG"
    INFO = "INFO"
    WARNING = "WARNING"
    ERROR = "ERROR"
    FATAL = "FATAL"


class Config(BaseSettings):
    """
    Application settings.

    These parameters can be configured
    with environment variables.
    """

    host: str = os.environ.get("HOST", "127.0.0.1")
    port: int = os.environ.get("PORT", 8000)
    # quantity of workers for uvicorn
    workers_count: int = os.environ.get("WORKER_COUNT", 1)
    # Enable uvicorn reloading
    if os.environ.get("APP_ENVIRONMENT", "prod") == "prod":
        reload: bool = False
        log_level: LogLevel = LogLevel.WARNING
    else:
        reload: bool = True
        log_level: LogLevel = LogLevel.INFO

    # Current environment
    # environment: str = "dev"


load_dotenv()
config = Config()
pinecone_client = PineconeClient(os.getenv("PINECONE_API_KEY"), os.getenv("PINECONE_ENVIRONMENT"),
                                 os.getenv("PINECONE_INDEX_NAME"))
