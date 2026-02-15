"""Model service."""

from __future__ import annotations

from math import ceil

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import DEFAULT_PAGE_NUM, DEFAULT_PAGE_SIZE, MAX_PAGE_SIZE, MODEL_STATUS_ACTIVE
from app.models.model import Model
from app.models.model_provider import ModelProvider
from app.schemas.common import PageData
from app.schemas.model import ModelQueryRequest, ModelVO


class ModelService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db

    async def add_model(self, model: Model) -> int:
        self.db.add(model)
        await self.db.commit()
        await self.db.refresh(model)
        return model.id

    async def update_model(self, model: Model) -> bool:
        entity = await self.get_by_id(model.id)
        if entity is None:
            return False
        for attr in [
            "model_name",
            "description",
            "context_length",
            "input_price",
            "output_price",
            "status",
            "priority",
            "default_timeout",
            "capabilities",
        ]:
            value = getattr(model, attr, None)
            if value is not None:
                setattr(entity, attr, value)
        await self.db.commit()
        return True

    async def delete_model(self, model_id: int) -> bool:
        entity = await self.get_by_id(model_id)
        if entity is None:
            return False
        entity.is_delete = 1
        await self.db.commit()
        return True

    async def get_by_id(self, model_id: int) -> Model | None:
        stmt = select(Model).where(Model.id == model_id, Model.is_delete == 0)
        return await self.db.scalar(stmt)

    async def get_by_model_key(self, model_key: str) -> Model | None:
        stmt = select(Model).where(Model.model_key == model_key, Model.is_delete == 0)
        return await self.db.scalar(stmt)

    async def get_model_vo(self, model_id: int) -> ModelVO | None:
        model = await self.get_by_id(model_id)
        if model is None:
            return None
        return await self._to_model_vo(model)

    async def list_page_vo(self, query: ModelQueryRequest) -> PageData[ModelVO]:
        page_num = query.page_num if query.page_num > 0 else DEFAULT_PAGE_NUM
        page_size = min(query.page_size if query.page_size > 0 else DEFAULT_PAGE_SIZE, MAX_PAGE_SIZE)
        stmt = select(Model).where(Model.is_delete == 0)
        if query.provider_id:
            stmt = stmt.where(Model.provider_id == query.provider_id)
        if query.model_key:
            stmt = stmt.where(Model.model_key.like(f"%{query.model_key}%"))
        if query.model_name:
            stmt = stmt.where(Model.model_name.like(f"%{query.model_name}%"))
        if query.model_type:
            stmt = stmt.where(Model.model_type == query.model_type)
        if query.status:
            stmt = stmt.where(Model.status == query.status)
        stmt = stmt.order_by(Model.priority.desc(), Model.create_time.desc())
        total = await self.db.scalar(select(func.count()).select_from(stmt.subquery())) or 0
        rows = (await self.db.scalars(stmt.offset((page_num - 1) * page_size).limit(page_size))).all()
        records = [await self._to_model_vo(item) for item in rows]
        return PageData[ModelVO](
            records=records,
            pageNumber=page_num,
            pageSize=page_size,
            totalPage=ceil(total / page_size) if total else 0,
            totalRow=total,
            optimizeCountQuery=True,
        )

    async def list_vo(self) -> list[ModelVO]:
        stmt = select(Model).where(Model.is_delete == 0)
        rows = (await self.db.scalars(stmt)).all()
        return [await self._to_model_vo(item) for item in rows]

    async def list_active(self) -> list[ModelVO]:
        stmt = select(Model).where(Model.is_delete == 0, Model.status == MODEL_STATUS_ACTIVE)
        rows = (await self.db.scalars(stmt)).all()
        return [await self._to_model_vo(item) for item in rows]

    async def list_active_by_provider(self, provider_id: int) -> list[ModelVO]:
        stmt = select(Model).where(
            Model.is_delete == 0,
            Model.status == MODEL_STATUS_ACTIVE,
            Model.provider_id == provider_id,
        )
        rows = (await self.db.scalars(stmt)).all()
        return [await self._to_model_vo(item) for item in rows]

    async def list_active_by_type(self, model_type: str) -> list[ModelVO]:
        stmt = select(Model).where(
            Model.is_delete == 0,
            Model.status == MODEL_STATUS_ACTIVE,
            Model.model_type == model_type,
        )
        rows = (await self.db.scalars(stmt)).all()
        return [await self._to_model_vo(item) for item in rows]

    async def list_active_entities(self) -> list[Model]:
        stmt = select(Model).where(Model.is_delete == 0, Model.status == MODEL_STATUS_ACTIVE)
        return list((await self.db.scalars(stmt)).all())

    async def update_model_metrics(
        self,
        model_id: int,
        health_status: str,
        avg_latency: int,
        success_rate,
        score,
    ) -> None:
        entity = await self.get_by_id(model_id)
        if entity is None:
            return
        entity.health_status = health_status
        entity.avg_latency = avg_latency
        entity.success_rate = success_rate
        entity.score = score
        await self.db.commit()

    async def _to_model_vo(self, model: Model) -> ModelVO:
        provider_stmt = select(ModelProvider).where(
            ModelProvider.id == model.provider_id,
            ModelProvider.is_delete == 0,
        )
        provider = await self.db.scalar(provider_stmt)
        return ModelVO.model_validate(
            {
                "id": model.id,
                "providerId": model.provider_id,
                "providerName": provider.provider_name if provider else None,
                "providerDisplayName": provider.display_name if provider else None,
                "modelKey": model.model_key,
                "modelName": model.model_name,
                "modelType": model.model_type,
                "description": model.description,
                "contextLength": model.context_length,
                "inputPrice": model.input_price,
                "outputPrice": model.output_price,
                "status": model.status,
                "healthStatus": model.health_status,
                "avgLatency": model.avg_latency,
                "successRate": model.success_rate,
                "priority": model.priority,
                "defaultTimeout": model.default_timeout,
                "supportReasoning": model.support_reasoning,
                "capabilities": model.capabilities,
                "createTime": model.create_time,
                "updateTime": model.update_time,
            }
        )
