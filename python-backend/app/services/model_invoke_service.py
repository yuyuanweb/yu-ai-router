"""Model invoke service."""

from __future__ import annotations

from collections.abc import AsyncGenerator

from app.adapter.adapter_factory import ModelAdapterFactory
from app.models.model import Model
from app.models.model_provider import ModelProvider
from app.schemas.chat import ChatRequest, ChatResponse, StreamChunk


class ModelInvokeService:
    def __init__(self) -> None:
        self.adapter_factory = ModelAdapterFactory()

    async def invoke(self, model: Model, provider: ModelProvider, chat_request: ChatRequest) -> ChatResponse:
        adapter = self.adapter_factory.get_adapter(provider.provider_name)
        return await adapter.invoke(model, provider, chat_request)

    async def invoke_stream_chunk(
        self, model: Model, provider: ModelProvider, chat_request: ChatRequest
    ) -> AsyncGenerator[StreamChunk, None]:
        adapter = self.adapter_factory.get_adapter(provider.provider_name)
        async for chunk in adapter.invoke_stream_chunk(model, provider, chat_request):
            yield chunk
