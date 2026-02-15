"""Image generation schemas."""

from __future__ import annotations

from datetime import datetime

from pydantic import AliasChoices, Field, field_serializer

from app.schemas.common import CamelBaseModel, LongIdModel


class ImageGenerationRequest(CamelBaseModel):
    prompt: str
    model: str | None = None
    n: int | None = 1
    size: str | None = None
    quality: str | None = None
    response_format: str | None = Field(
        default=None,
        alias="response_format",
        validation_alias=AliasChoices("response_format", "responseFormat"),
    )
    user: str | None = None


class ImageData(CamelBaseModel):
    url: str | None = None
    b64_json: str | None = Field(default=None, alias="b64Json")
    revised_prompt: str | None = Field(default=None, alias="revisedPrompt")


class ImageGenerationResponse(CamelBaseModel):
    created: int
    data: list[ImageData]


class ImageGenerationRecordVO(LongIdModel):
    user_id: int = Field(alias="userId")
    api_key_id: int | None = Field(default=None, alias="apiKeyId")
    model_id: int = Field(alias="modelId")
    model_key: str = Field(alias="modelKey")
    prompt: str
    revised_prompt: str | None = Field(default=None, alias="revisedPrompt")
    image_url: str | None = Field(default=None, alias="imageUrl")
    image_data: str | None = Field(default=None, alias="imageData")
    size: str | None = None
    quality: str | None = None
    status: str
    cost: float
    duration: int | None = None
    error_message: str | None = Field(default=None, alias="errorMessage")
    client_ip: str | None = Field(default=None, alias="clientIp")
    create_time: datetime | None = Field(default=None, alias="createTime")

    @field_serializer("user_id", when_used="json")
    def serialize_user_id(self, value: int) -> str:
        return str(value)

    @field_serializer("api_key_id", when_used="json")
    def serialize_api_key_id(self, value: int | None) -> str | None:
        return str(value) if value is not None else None

    @field_serializer("model_id", when_used="json")
    def serialize_model_id(self, value: int) -> str:
        return str(value)
