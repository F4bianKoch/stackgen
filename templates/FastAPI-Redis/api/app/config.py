import os

REDIS_URL = os.environ.get("REDIS_URL", "redis://redis:6379/0")
