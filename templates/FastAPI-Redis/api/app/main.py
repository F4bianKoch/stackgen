import os
from contextlib import asynccontextmanager
from fastapi import FastAPI
from redis.asyncio import Redis

from .config import REDIS_URL

@asynccontextmanager
async def lifespan(app: FastAPI):
    redis = Redis.from_url(REDIS_URL, decode_responses=False)
    app.state.redis = redis
    await redis.ping()
    try:
        yield
    finally:
        await redis.aclose()

app = FastAPI(title="<project Name>", version="0.1.0", lifespan=lifespan)

@app.get("/health")
async def health():
    return {"status": "ok"}
