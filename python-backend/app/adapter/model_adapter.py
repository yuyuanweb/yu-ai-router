"""Model adapter interface."""

from __future__ import annotations

from abc import ABC, abstractmethod
from collections.abc import AsyncGenerator

from app.models.model import Model
from app.models.model_provider import ModelProvider
from app.schemas.chat import ChatRequest, ChatResponse, StreamChunk


class ModelAdapter(ABC):
    @abstractmethod
    async def invoke(self, model: Model, provider: ModelProvider, chat_request: ChatRequest) -> ChatResponse:
        raise NotImplementedError

    @abstractmethod
    async def invoke_stream_chunk(
        self, model: Model, provider: ModelProvider, chat_request: ChatRequest
    ) -> AsyncGenerator[StreamChunk, None]:
        raise NotImplementedError

    @abstractmethod
    def supports(self, provider_name: str) -> bool:
        raise NotImplementedError
