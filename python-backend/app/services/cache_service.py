"""Cache service for chat response."""

from __future__ import annotations

import hashlib
import json

from app.core.config import get_settings
from app.core.constants import CACHE_KEY_PREFIX
from app.infra.redis_client import redis_client
from app.schemas.chat import ChatRequest, ChatResponse


class CacheService:
    def __init__(self) -> None:
        self.settings = get_settings()

    async def get_cached_response(self, request: ChatRequest) -> ChatResponse | None:
        if not self.is_cache_enabled():
            return None
        cache_key = self.generate_cache_key(request)
        cached = await redis_client.get(cache_key)
        if not cached:
            return None
        try:
            return ChatResponse.model_validate_json(cached)
        except Exception:
            return None

    async def cache_response(self, request: ChatRequest, response: ChatResponse) -> None:
        if not self.is_cache_enabled():
            return
        cache_key = self.generate_cache_key(request)
        await redis_client.set(cache_key, response.model_dump_json(by_alias=True), ex=self.settings.ai_cache_ttl)

    @staticmethod
    def generate_cache_key(request: ChatRequest) -> str:
        model = request.model or "auto"
        messages_json = json.dumps([item.model_dump(by_alias=True) for item in request.messages], ensure_ascii=False)
        messages_md5 = hashlib.md5(messages_json.encode("utf-8")).hexdigest()  # noqa: S324
        key = f"{CACHE_KEY_PREFIX}{model}:{messages_md5}"
        if request.temperature is not None:
            key = f"{key}:t{request.temperature}"
        return key

    def clear_user_cache(self, user_id: int) -> None:
        _ = user_id

    def is_cache_enabled(self) -> bool:
        return self.settings.ai_cache_enabled
