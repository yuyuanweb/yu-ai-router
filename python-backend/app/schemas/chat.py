"""Chat schemas."""

from __future__ import annotations

from typing import Any

from pydantic import Field

from app.schemas.common import CamelBaseModel


class ChatMessage(CamelBaseModel):
    role: str
    content: str


class ChatRequest(CamelBaseModel):
    model: str | None = None
    messages: list[ChatMessage]
    stream: bool | None = False
    temperature: float | None = None
    max_tokens: int | None = Field(default=None, alias="max_tokens")


class ChatChoice(CamelBaseModel):
    index: int
    message: ChatMessage
    finish_reason: str | None = Field(default=None, alias="finishReason")


class ChatUsage(CamelBaseModel):
    prompt_tokens: int = Field(alias="promptTokens")
    completion_tokens: int = Field(alias="completionTokens")
    total_tokens: int = Field(alias="totalTokens")


class ChatResponse(CamelBaseModel):
    id: str
    object: str
    created: int
    model: str
    choices: list[ChatChoice]
    usage: ChatUsage


class OpenAiMessage(CamelBaseModel):
    role: str | None = None
    content: str | None = None


class OpenAiChoice(CamelBaseModel):
    index: int
    message: OpenAiMessage | None = None
    delta: OpenAiMessage | None = None
    finish_reason: str | None = Field(default=None, alias="finish_reason")


class OpenAiUsage(CamelBaseModel):
    prompt_tokens: int | None = 0
    completion_tokens: int | None = 0
    total_tokens: int | None = 0


class OpenAiChatResponse(CamelBaseModel):
    id: str | None = None
    object: str | None = None
    created: int | None = None
    model: str | None = None
    choices: list[OpenAiChoice] = Field(default_factory=list)
    usage: OpenAiUsage | None = None


class OpenAiStreamChunk(CamelBaseModel):
    choices: list[OpenAiChoice] = Field(default_factory=list)
    usage: OpenAiUsage | None = None
    raw: dict[str, Any] | None = None
