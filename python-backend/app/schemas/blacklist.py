"""Blacklist schemas."""

from __future__ import annotations

from app.schemas.common import CamelBaseModel


class BlacklistRequest(CamelBaseModel):
    ip: str
    reason: str | None = None
