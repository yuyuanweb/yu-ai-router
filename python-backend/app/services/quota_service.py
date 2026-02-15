"""Quota service."""

from __future__ import annotations

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import UNLIMITED_QUOTA
from app.models.user import User


class QuotaService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db

    async def check_quota(self, user_id: int | None) -> bool:
        if user_id is None:
            return True
        user = await self.db.scalar(select(User).where(User.id == user_id, User.is_delete == 0))
        if user is None:
            return False
        if user.token_quota in (None, UNLIMITED_QUOTA):
            return True
        return (user.used_tokens or 0) < user.token_quota

    async def deduct_tokens(self, user_id: int | None, tokens: int) -> bool:
        if user_id is None or tokens <= 0:
            return True
        user = await self.db.scalar(select(User).where(User.id == user_id, User.is_delete == 0))
        if user is None:
            return False
        user.used_tokens = (user.used_tokens or 0) + tokens
        await self.db.commit()
        return True

    async def get_remaining_quota(self, user_id: int | None) -> int:
        if user_id is None:
            return UNLIMITED_QUOTA
        user = await self.db.scalar(select(User).where(User.id == user_id, User.is_delete == 0))
        if user is None:
            return 0
        if user.token_quota in (None, UNLIMITED_QUOTA):
            return UNLIMITED_QUOTA
        return max(0, user.token_quota - (user.used_tokens or 0))
