"""Plugin interface."""

from __future__ import annotations

from abc import ABC, abstractmethod

from app.plugin.plugin_context import PluginContext
from app.plugin.plugin_result import PluginResult


class Plugin(ABC):
    @abstractmethod
    def get_plugin_key(self) -> str:
        pass

    @abstractmethod
    def get_plugin_name(self) -> str:
        pass

    @abstractmethod
    def get_description(self) -> str:
        pass

    @abstractmethod
    async def execute(self, context: PluginContext) -> PluginResult:
        pass

    def supports(self, context: PluginContext) -> bool:
        _ = context
        return True

    def init(self, config: str | None) -> None:
        _ = config

    def destroy(self) -> None:
        return None
