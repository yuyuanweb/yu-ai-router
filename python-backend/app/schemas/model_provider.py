"""Model provider schemas."""

from __future__ import annotations

from decimal import Decimal

from pydantic import Field

from app.schemas.common import LongIdModel, PageRequest, TimeModel


class ProviderAddRequest(LongIdModel):
    provider_name: str = Field(alias="providerName")
    display_name: str = Field(alias="displayName")
    base_url: str = Field(alias="baseUrl")
    api_key: str = Field(alias="apiKey")
    priority: int | None = None
    config: str | None = None


class ProviderUpdateRequest(LongIdModel):
    display_name: str | None = Field(default=None, alias="displayName")
    base_url: str | None = Field(default=None, alias="baseUrl")
    api_key: str | None = Field(default=None, alias="apiKey")
    status: str | None = None
    priority: int | None = None
    config: str | None = None


class ProviderQueryRequest(PageRequest):
    provider_name: str | None = Field(default=None, alias="providerName")
    display_name: str | None = Field(default=None, alias="displayName")
    status: str | None = None
    health_status: str | None = Field(default=None, alias="healthStatus")


class ProviderVO(LongIdModel, TimeModel):
    provider_name: str = Field(alias="providerName")
    display_name: str = Field(alias="displayName")
    base_url: str = Field(alias="baseUrl")
    status: str
    health_status: str = Field(alias="healthStatus")
    avg_latency: int = Field(alias="avgLatency")
    success_rate: Decimal = Field(alias="successRate")
    priority: int
    config: str | None = None
