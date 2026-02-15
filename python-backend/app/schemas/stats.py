"""Stats schemas."""

from __future__ import annotations

from datetime import datetime

from pydantic import Field

from app.schemas.common import CamelBaseModel, LongIdModel


class RequestLogVO(LongIdModel):
    user_id: int | None = Field(default=None, alias="userId")
    api_key_id: int | None = Field(default=None, alias="apiKeyId")
    model_name: str | None = Field(default=None, alias="modelName")
    prompt_tokens: int | None = Field(default=None, alias="promptTokens")
    completion_tokens: int | None = Field(default=None, alias="completionTokens")
    total_tokens: int | None = Field(default=None, alias="totalTokens")
    duration: int | None = None
    status: str | None = None
    error_message: str | None = Field(default=None, alias="errorMessage")
    create_time: datetime | None = Field(default=None, alias="createTime")
    update_time: datetime | None = Field(default=None, alias="updateTime")


class TokenStatsVO(CamelBaseModel):
    total_tokens: int = Field(alias="totalTokens")
