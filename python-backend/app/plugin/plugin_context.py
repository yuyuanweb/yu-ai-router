"""Plugin execution context."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any


@dataclass(slots=True)
class PluginContext:
    user_id: int | None = None
    input: str | None = None
    file_url: str | None = None
    file_bytes: bytes | None = None
    file_type: str | None = None
    params: dict[str, Any] = field(default_factory=dict)

    def get_param(self, key: str, default: Any = None) -> Any:
        return self.params.get(key, default)
