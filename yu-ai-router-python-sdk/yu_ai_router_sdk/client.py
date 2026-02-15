"""Yu AI SDK main client."""

from __future__ import annotations

from yu_ai_router_sdk.callback import StreamCallback
from yu_ai_router_sdk.config import ClientConfig
from yu_ai_router_sdk.http_client import HttpClient
from yu_ai_router_sdk.models import ChatRequest, ChatResponse


class YuAIClient:
    def __init__(self, config: ClientConfig) -> None:
        config.validate()
        self.config = config
        self.http_client = HttpClient(config)

    @classmethod
    def builder(cls) -> "ClientConfigBuilder":
        return ClientConfigBuilder()

    def chat(self, request_or_message: ChatRequest | str, model: str | None = None) -> ChatResponse:
        request = self._to_request(request_or_message, model)
        return self.http_client.chat(request)

    def chat_stream(
        self,
        request_or_message: ChatRequest | str,
        callback: StreamCallback,
        model: str | None = None,
    ) -> None:
        request = self._to_request(request_or_message, model)
        self.http_client.chat_stream(request, callback)

    def close(self) -> None:
        self.http_client.close()

    @staticmethod
    def _to_request(request_or_message: ChatRequest | str, model: str | None) -> ChatRequest:
        if isinstance(request_or_message, ChatRequest):
            return request_or_message
        if model:
            return ChatRequest.with_model(model, request_or_message)
        return ChatRequest.simple(request_or_message)


class ClientConfigBuilder:
    def __init__(self) -> None:
        self._api_key: str | None = None
        self._base_url = "http://localhost:8123/api"
        self._connect_timeout = 10000
        self._read_timeout = 30000
        self._write_timeout = 30000
        self._max_retries = 3
        self._retry_delay = 1000

    def api_key(self, api_key: str) -> "ClientConfigBuilder":
        self._api_key = api_key
        return self

    def base_url(self, base_url: str) -> "ClientConfigBuilder":
        self._base_url = base_url
        return self

    def connect_timeout(self, connect_timeout: int) -> "ClientConfigBuilder":
        self._connect_timeout = connect_timeout
        return self

    def read_timeout(self, read_timeout: int) -> "ClientConfigBuilder":
        self._read_timeout = read_timeout
        return self

    def write_timeout(self, write_timeout: int) -> "ClientConfigBuilder":
        self._write_timeout = write_timeout
        return self

    def max_retries(self, max_retries: int) -> "ClientConfigBuilder":
        self._max_retries = max_retries
        return self

    def retry_delay(self, retry_delay: int) -> "ClientConfigBuilder":
        self._retry_delay = retry_delay
        return self

    def build(self) -> YuAIClient:
        config = ClientConfig(
            api_key=self._api_key or "",
            base_url=self._base_url,
            connect_timeout=self._connect_timeout,
            read_timeout=self._read_timeout,
            write_timeout=self._write_timeout,
            max_retries=self._max_retries,
            retry_delay=self._retry_delay,
        )
        return YuAIClient(config)
