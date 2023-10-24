from fastapi import APIRouter
from admin import ingest
from user import discovery

v1Router = APIRouter(prefix="/v1")
v1Router.include_router(router=ingest.router, prefix="/admin", tags=["admin"])
v1Router.include_router(router=discovery.router,tags=["user"])
