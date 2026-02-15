"""IP blacklist service based on Redis Set."""

from __future__ import annotations

import logging

from redis.asyncio import Redis

from app.core.constants import BLACKLIST_IP_KEY

logger = logging.getLogger("app")


class BlacklistService:
    def __init__(self, redis: Redis) -> None:
        self.redis = redis

    async def is_blocked(self, ip: str | None) -> bool:
        if not ip or not ip.strip():
            return False
        return bool(await self.redis.sismember(BLACKLIST_IP_KEY, ip))

    async def add_to_blacklist(self, ip: str | None, reason: str | None = None) -> None:
        if not ip or not ip.strip():
            return
        await self.redis.sadd(BLACKLIST_IP_KEY, ip)
        logger.info("IP added to blacklist: %s, reason: %s", ip, reason)

    async def remove_from_blacklist(self, ip: str | None) -> None:
        if not ip or not ip.strip():
            return
        await self.redis.srem(BLACKLIST_IP_KEY, ip)
        logger.info("IP removed from blacklist: %s", ip)

    async def get_all_blacklist(self) -> set[str]:
        members = await self.redis.smembers(BLACKLIST_IP_KEY)
        return set(members)

    async def clear_blacklist(self) -> None:
        await self.redis.delete(BLACKLIST_IP_KEY)
        logger.info("Blacklist cleared")

    async def get_blacklist_count(self) -> int:
        return int(await self.redis.scard(BLACKLIST_IP_KEY))
