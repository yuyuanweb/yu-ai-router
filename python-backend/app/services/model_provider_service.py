"""Model provider service."""

from __future__ import annotations

from math import ceil

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import (
    DEFAULT_PAGE_NUM,
    DEFAULT_PAGE_SIZE,
    HEALTH_STATUS_DEGRADED,
    HEALTH_STATUS_HEALTHY,
    MAX_PAGE_SIZE,
    PROVIDER_STATUS_ACTIVE,
)
from app.models.model_provider import ModelProvider
from app.schemas.common import PageData
from app.schemas.model_provider import ProviderQueryRequest, ProviderVO


class ModelProviderService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db

    async def add_provider(self, provider: ModelProvider) -> int:
        self.db.add(provider)
        await self.db.commit()
        await self.db.refresh(provider)
        return provider.id

    async def update_provider(self, provider: ModelProvider) -> bool:
        entity = await self.get_by_id(provider.id)
        if entity is None:
            return False
        for attr in ["display_name", "base_url", "api_key", "status", "priority", "config"]:
            value = getattr(provider, attr, None)
            if value is not None:
                setattr(entity, attr, value)
        await self.db.commit()
        return True

    async def delete_provider(self, provider_id: int) -> bool:
        entity = await self.get_by_id(provider_id)
        if entity is None:
            return False
        entity.is_delete = 1
        await self.db.commit()
        return True

    async def get_by_id(self, provider_id: int) -> ModelProvider | None:
        stmt = select(ModelProvider).where(ModelProvider.id == provider_id, ModelProvider.is_delete == 0)
        return await self.db.scalar(stmt)

    async def get_by_name(self, provider_name: str) -> ModelProvider | None:
        stmt = select(ModelProvider).where(
            ModelProvider.provider_name == provider_name,
            ModelProvider.is_delete == 0,
        )
        return await self.db.scalar(stmt)

    async def list_page_vo(self, query: ProviderQueryRequest) -> PageData[ProviderVO]:
        page_num = query.page_num if query.page_num > 0 else DEFAULT_PAGE_NUM
        page_size = min(query.page_size if query.page_size > 0 else DEFAULT_PAGE_SIZE, MAX_PAGE_SIZE)
        stmt = select(ModelProvider).where(ModelProvider.is_delete == 0)
        if query.provider_name:
            stmt = stmt.where(ModelProvider.provider_name.like(f"%{query.provider_name}%"))
        if query.display_name:
            stmt = stmt.where(ModelProvider.display_name.like(f"%{query.display_name}%"))
        if query.status:
            stmt = stmt.where(ModelProvider.status == query.status)
        if query.health_status:
            stmt = stmt.where(ModelProvider.health_status == query.health_status)
        stmt = stmt.order_by(ModelProvider.priority.desc(), ModelProvider.create_time.desc())
        total = await self.db.scalar(select(func.count()).select_from(stmt.subquery())) or 0
        rows = (await self.db.scalars(stmt.offset((page_num - 1) * page_size).limit(page_size))).all()
        records = [ProviderVO.model_validate(item) for item in rows]
        return PageData[ProviderVO](
            records=records,
            pageNumber=page_num,
            pageSize=page_size,
            totalPage=ceil(total / page_size) if total else 0,
            totalRow=total,
            optimizeCountQuery=True,
        )

    async def list_vo(self) -> list[ProviderVO]:
        stmt = select(ModelProvider).where(ModelProvider.is_delete == 0)
        rows = (await self.db.scalars(stmt)).all()
        return [ProviderVO.model_validate(item) for item in rows]

    async def list_healthy(self) -> list[ProviderVO]:
        stmt = (
            select(ModelProvider)
            .where(
                ModelProvider.is_delete == 0,
                ModelProvider.status == PROVIDER_STATUS_ACTIVE,
                ModelProvider.health_status.in_([HEALTH_STATUS_HEALTHY, HEALTH_STATUS_DEGRADED]),
            )
            .order_by(ModelProvider.priority.desc())
        )
        rows = (await self.db.scalars(stmt)).all()
        return [ProviderVO.model_validate(item) for item in rows]

    async def list_entities(self) -> list[ModelProvider]:
        stmt = select(ModelProvider).where(ModelProvider.is_delete == 0)
        return list((await self.db.scalars(stmt)).all())

    async def update_health_status(
        self,
        provider_id: int,
        health_status: str,
        avg_latency: int,
        success_rate,
    ) -> None:
        provider = await self.get_by_id(provider_id)
        if provider is None:
            return
        provider.health_status = health_status
        provider.avg_latency = avg_latency
        provider.success_rate = success_rate
        await self.db.commit()
