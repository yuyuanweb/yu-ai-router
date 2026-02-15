"""Auto routing strategy."""

from __future__ import annotations

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import (
    HEALTH_STATUS_DEGRADED,
    HEALTH_STATUS_HEALTHY,
    HEALTH_STATUS_UNKNOWN,
    MODEL_STATUS_ACTIVE,
    ROUTING_STRATEGY_AUTO,
)
from app.models.model import Model
from app.strategy.routing_strategy import RoutingStrategy


class AutoRoutingStrategy(RoutingStrategy):
    def __init__(self, db: AsyncSession) -> None:
        self.db = db

    async def select_model(self, model_type: str, requested_model: str | None) -> Model | None:
        stmt = self._base_stmt(model_type).order_by(Model.score.asc()).limit(1)
        return await self.db.scalar(stmt)

    async def get_fallback_models(self, model_type: str, requested_model: str | None) -> list[Model]:
        stmt = self._base_stmt(model_type).order_by(Model.score.asc())
        models = (await self.db.scalars(stmt)).all()
        return list(models[1:]) if len(models) > 1 else []

    def get_strategy_type(self) -> str:
        return ROUTING_STRATEGY_AUTO

    @staticmethod
    def _base_stmt(model_type: str):
        return select(Model).where(
            Model.is_delete == 0,
            Model.status == MODEL_STATUS_ACTIVE,
            Model.health_status.in_(
                [HEALTH_STATUS_HEALTHY, HEALTH_STATUS_DEGRADED, HEALTH_STATUS_UNKNOWN]
            ),
            Model.model_type == model_type,
        )
