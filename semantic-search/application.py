from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import UJSONResponse

from controllers import router

# origins = [
#     "http://localhost",
#     "http://localhost:3000",
#     # Add other allowed origins as needed
# ]
# app.add_middleware(
#     CORSMiddleware,
#     allow_origins=origins,
#     allow_credentials=True,
#     allow_methods=["*"],
#     allow_headers=["*"],
# )


def get_app() -> FastAPI:
    app = FastAPI(
        title="Tickrlytics PY Backend",
        docs_url="/api/docs",
        redoc_url="/api/redoc",
        openapi_url="/api/openapi.json",
        default_response_class=UJSONResponse,
    )
    app.include_router(router=router.v1Router, prefix="/api")
    return app

# @app.post("/r2/discovery")
# async def search(query: str):
#         return discover(query)
#
# @app.post("/r2/admin/ingest")
# async def ingest_doc(doc: Document):
#         ingest(doc)
#         return "document ingested successfully"
