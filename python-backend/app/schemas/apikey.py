"""API key schemas."""

from __future__ import annotations

from datetime import datetime

from pydantic import Field

from app.schemas.common import CamelBaseModel, LongIdModel


class ApiKeyCreateRequest(CamelBaseModel):
    key_name: str | None = Field(default=None, alias="keyName")


class ApiKeyVO(LongIdModel):
    key_value: str | None = Field(default=None, alias="keyValue")
    key_name: str | None = Field(default=None, alias="keyName")
    status: str | None = None
    total_tokens: int | None = Field(default=None, alias="totalTokens")
    last_used_time: datetime | None = Field(default=None, alias="lastUsedTime")
    create_time: datetime | None = Field(default=None, alias="createTime")
