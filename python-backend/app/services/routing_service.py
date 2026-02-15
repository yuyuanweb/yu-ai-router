"""Routing service."""

from __future__ import annotations

from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import ROUTING_STRATEGY_AUTO
from app.models.model import Model
from app.strategy.auto_routing_strategy import AutoRoutingStrategy
from app.strategy.cost_first_routing_strategy import CostFirstRoutingStrategy
from app.strategy.fixed_routing_strategy import FixedRoutingStrategy
from app.strategy.latency_first_routing_strategy import LatencyFirstRoutingStrategy
from app.strategy.routing_strategy import RoutingStrategy


class RoutingService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db
        self.routing_strategies: list[RoutingStrategy] = [
            AutoRoutingStrategy(db),
            FixedRoutingStrategy(db),
            CostFirstRoutingStrategy(db),
            LatencyFirstRoutingStrategy(db),
        ]

    async def select_model(
        self, strategy_type: str | None, model_type: str, requested_model: str | None
    ) -> Model | None:
        strategy = self._get_strategy(strategy_type) or self._get_strategy(ROUTING_STRATEGY_AUTO)
        if strategy is None:
            return None
        return await strategy.select_model(model_type, requested_model)

    async def get_fallback_models(
        self, strategy_type: str | None, model_type: str, requested_model: str | None
    ) -> list[Model]:
        strategy = self._get_strategy(strategy_type) or self._get_strategy(ROUTING_STRATEGY_AUTO)
        if strategy is None:
            return []
        return await strategy.get_fallback_models(model_type, requested_model)

    def _get_strategy(self, strategy_type: str | None) -> RoutingStrategy | None:
        if not strategy_type:
            return None
        for strategy in self.routing_strategies:
            if strategy.get_strategy_type() == strategy_type:
                return strategy
        return None
