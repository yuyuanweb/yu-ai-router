"""BYOK schemas."""

from __future__ import annotations

from datetime import datetime

from pydantic import Field

from app.schemas.common import CamelBaseModel, LongIdModel


class UserProviderKeyAddRequest(CamelBaseModel):
    provider_id: int = Field(alias="providerId")
    api_key: str = Field(alias="apiKey")


class UserProviderKeyUpdateRequest(CamelBaseModel):
    id: int
    api_key: str | None = Field(default=None, alias="apiKey")
    status: str | None = None


class UserProviderKeyVO(LongIdModel):
    provider_id: int | None = Field(default=None, alias="providerId")
    provider_name: str | None = Field(default=None, alias="providerName")
    api_key: str | None = Field(default=None, alias="apiKey")
    status: str | None = None
    create_time: datetime | None = Field(default=None, alias="createTime")
    update_time: datetime | None = Field(default=None, alias="updateTime")
