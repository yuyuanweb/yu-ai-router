"""Plugin execution result."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any


@dataclass(slots=True)
class PluginResult:
    success: bool
    content: str | None = None
    error_message: str | None = None
    duration: int = 0
    data: dict[str, Any] = field(default_factory=dict)

    @staticmethod
    def success_result(content: str, data: dict[str, Any] | None = None) -> "PluginResult":
        return PluginResult(success=True, content=content, data=data or {})

    @staticmethod
    def fail(error_message: str) -> "PluginResult":
        return PluginResult(success=False, error_message=error_message)
