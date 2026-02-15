"""Routing strategy interface."""

from __future__ import annotations

from abc import ABC, abstractmethod

from app.models.model import Model


class RoutingStrategy(ABC):
    @abstractmethod
    async def select_model(self, model_type: str, requested_model: str | None) -> Model | None:
        raise NotImplementedError

    @abstractmethod
    async def get_fallback_models(self, model_type: str, requested_model: str | None) -> list[Model]:
        raise NotImplementedError

    @abstractmethod
    def get_strategy_type(self) -> str:
        raise NotImplementedError
