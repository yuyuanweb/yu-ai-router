"""API key service."""

from __future__ import annotations

import secrets
from datetime import datetime
from math import ceil

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import (
    API_KEY_PREFIX,
    API_KEY_STATUS_ACTIVE,
    API_KEY_STATUS_REVOKED,
    DEFAULT_PAGE_NUM,
    DEFAULT_PAGE_SIZE,
    ErrorCode,
    MAX_PAGE_SIZE,
)
from app.exceptions.business_exception import BusinessException
from app.models.api_key import ApiKey
from app.schemas.apikey import ApiKeyVO
from app.schemas.common import PageData


class ApiKeyService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db

    async def create_api_key(self, key_name: str | None, user_id: int) -> ApiKey:
        key_value = f"{API_KEY_PREFIX}{secrets.token_hex(16)}"
        entity = ApiKey(
            user_id=user_id,
            key_value=key_value,
            key_name=key_name,
            status=API_KEY_STATUS_ACTIVE,
            total_tokens=0,
            create_time=datetime.utcnow(),
            update_time=datetime.utcnow(),
        )
        self.db.add(entity)
        await self.db.commit()
        await self.db.refresh(entity)
        return entity

    async def get_by_id(self, api_key_id: int) -> ApiKey | None:
        stmt = select(ApiKey).where(ApiKey.id == api_key_id, ApiKey.is_delete == 0)
        return await self.db.scalar(stmt)

    async def get_by_key_value(self, key_value: str) -> ApiKey | None:
        stmt = select(ApiKey).where(
            ApiKey.key_value == key_value,
            ApiKey.status == API_KEY_STATUS_ACTIVE,
            ApiKey.is_delete == 0,
        )
        return await self.db.scalar(stmt)

    async def revoke_api_key(self, api_key_id: int, user_id: int) -> bool:
        stmt = select(ApiKey).where(
            ApiKey.id == api_key_id,
            ApiKey.user_id == user_id,
            ApiKey.is_delete == 0,
        )
        entity = await self.db.scalar(stmt)
        if entity is None:
            raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "API Key 不存在")
        entity.status = API_KEY_STATUS_REVOKED
        entity.update_time = datetime.utcnow()
        await self.db.commit()
        return True

    async def list_user_api_key_page(self, user_id: int, page_num: int, page_size: int) -> PageData[ApiKeyVO]:
        page_num = page_num if page_num > 0 else DEFAULT_PAGE_NUM
        page_size = min(page_size if page_size > 0 else DEFAULT_PAGE_SIZE, MAX_PAGE_SIZE)
        stmt = (
            select(ApiKey)
            .where(ApiKey.user_id == user_id, ApiKey.is_delete == 0)
            .order_by(ApiKey.create_time.desc())
        )
        count = await self.db.scalar(select(func.count()).select_from(stmt.subquery())) or 0
        total_page = ceil(count / page_size) if count else 0
        rows = (await self.db.scalars(stmt.offset((page_num - 1) * page_size).limit(page_size))).all()
        records = [self.to_api_key_vo(item, mask=True) for item in rows]
        return PageData[ApiKeyVO](
            records=records,
            pageNumber=page_num,
            pageSize=page_size,
            totalPage=total_page,
            totalRow=count,
            optimizeCountQuery=True,
        )

    async def update_usage_stats(self, api_key_id: int, tokens: int) -> None:
        entity = await self.get_by_id(api_key_id)
        if entity is None:
            return
        entity.total_tokens = (entity.total_tokens or 0) + tokens
        entity.last_used_time = datetime.utcnow()
        entity.update_time = datetime.utcnow()
        await self.db.commit()

    @staticmethod
    def to_api_key_vo(api_key: ApiKey, mask: bool) -> ApiKeyVO:
        vo = ApiKeyVO.model_validate(api_key)
        if mask and vo.key_value and len(vo.key_value) > 12:
            vo.key_value = f"{vo.key_value[:8]}****{vo.key_value[-4:]}"
        return vo
