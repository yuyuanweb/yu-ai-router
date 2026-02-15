"""Plugin schemas."""

from __future__ import annotations

from typing import Any

from pydantic import AliasChoices, Field

from app.schemas.common import CamelBaseModel, LongIdModel, TimeModel


class PluginExecuteRequest(CamelBaseModel):
    plugin_key: str = Field(
        alias="pluginKey",
        validation_alias=AliasChoices("pluginKey", "plugin_key"),
    )
    input: str | None = None
    file_url: str | None = Field(
        default=None,
        alias="fileUrl",
        validation_alias=AliasChoices("fileUrl", "file_url"),
    )
    file_bytes: bytes | None = Field(default=None, exclude=True)
    file_type: str | None = Field(
        default=None,
        alias="fileType",
        validation_alias=AliasChoices("fileType", "file_type"),
    )
    params: dict[str, Any] | None = None


class PluginUpdateRequest(LongIdModel):
    plugin_name: str | None = Field(default=None, alias="pluginName")
    description: str | None = None
    config: str | None = None
    status: str | None = None
    priority: int | None = None


class PluginExecuteVO(CamelBaseModel):
    success: bool
    plugin_key: str = Field(alias="pluginKey")
    content: str | None = None
    error_message: str | None = Field(default=None, alias="errorMessage")
    duration: int
    data: dict[str, Any] | None = None


class PluginConfigVO(LongIdModel, TimeModel):
    plugin_key: str = Field(alias="pluginKey")
    plugin_name: str = Field(alias="pluginName")
    plugin_type: str = Field(alias="pluginType")
    description: str | None = None
    config: str | None = None
    status: str
    priority: int
