"""Rate limit service based on Redis fixed window."""

from __future__ import annotations

from redis.asyncio import Redis

from app.core.constants import RATE_LIMIT_API_KEY_PREFIX, RATE_LIMIT_IP_PREFIX


class RateLimitService:
    def __init__(self, redis: Redis) -> None:
        self.redis = redis

    async def try_acquire(self, key: str, limit: int, window_seconds: int) -> bool:
        current = await self.redis.incr(key)
        if current == 1:
            await self.redis.expire(key, window_seconds)
        return int(current) <= limit

    async def get_available_permits(self, key: str, limit: int) -> int:
        current = await self.redis.get(key)
        used = int(current) if current else 0
        remain = limit - used
        return remain if remain > 0 else 0

    async def check_api_key_rate_limit(self, api_key: str, limit: int) -> bool:
        key = f"{RATE_LIMIT_API_KEY_PREFIX}{api_key}"
        return await self.try_acquire(key, limit, 1)

    async def check_ip_rate_limit(self, ip: str, limit: int) -> bool:
        key = f"{RATE_LIMIT_IP_PREFIX}{ip}"
        return await self.try_acquire(key, limit, 1)
