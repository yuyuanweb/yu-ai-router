"""HTTP client wrapper for SDK."""

from __future__ import annotations

import json
import time

import httpx

from yu_ai_router_sdk.callback import StreamCallback
from yu_ai_router_sdk.config import ClientConfig
from yu_ai_router_sdk.exceptions import AuthException, RateLimitException, YuAIException
from yu_ai_router_sdk.models import ChatChunk, ChatRequest, ChatResponse, StreamResponse


class HttpClient:
    def __init__(self, config: ClientConfig) -> None:
        self.config = config
        timeout = httpx.Timeout(
            connect=config.connect_timeout / 1000,
            read=config.read_timeout / 1000,
            write=config.write_timeout / 1000,
            pool=config.connect_timeout / 1000,
        )
        self.client = httpx.Client(timeout=timeout)

    def chat(self, request: ChatRequest) -> ChatResponse:
        request.stream = False
        payload = request.to_payload()
        url = f"{self.config.base_url}/v1/chat/completions"
        last_error: Exception | None = None

        for attempt in range(self.config.max_retries + 1):
            try:
                response = self.client.post(url, json=payload, headers=self._headers())
                self._raise_for_status(response)
                return ChatResponse.from_dict(response.json())
            except (AuthException, RateLimitException):
                raise
            except Exception as exc:  # noqa: PERF203
                last_error = exc
                if attempt < self.config.max_retries:
                    time.sleep((self.config.retry_delay * (attempt + 1)) / 1000)

        raise YuAIException(
            f"Request failed after {self.config.max_retries} retries: {last_error}",
        )

    def chat_stream(self, request: ChatRequest, callback: StreamCallback) -> None:
        request.stream = True
        payload = request.to_payload()
        url = f"{self.config.base_url}/v1/chat/completions"

        try:
            with self.client.stream("POST", url, json=payload, headers=self._headers()) as response:
                self._raise_for_status(response)
                completed = False
                for line in response.iter_lines():
                    if not line:
                        continue
                    raw = line.strip()
                    if not raw.startswith("data:"):
                        continue
                    data = raw[5:].strip()
                    if not data:
                        continue
                    try:
                        parsed = json.loads(data)
                        stream_response = StreamResponse.from_dict(parsed)
                    except Exception:
                        continue
                    if not stream_response.choices:
                        continue
                    choice = stream_response.choices[0]
                    if choice.finish_reason == "stop":
                        callback.on_complete()
                        completed = True
                        break
                    if choice.delta.content is None and choice.delta.reasoning_content is None:
                        continue
                    callback.on_message(
                        ChatChunk(
                            content=choice.delta.content,
                            reasoning_content=choice.delta.reasoning_content,
                            done=False,
                            model=stream_response.model,
                        )
                    )
                if not completed:
                    callback.on_complete()
        except Exception as exc:
            callback.on_error(exc)

    def close(self) -> None:
        self.client.close()

    def _headers(self) -> dict[str, str]:
        return {
            "Authorization": f"Bearer {self.config.api_key}",
            "Content-Type": "application/json",
        }

    @staticmethod
    def _raise_for_status(response: httpx.Response) -> None:
        if response.status_code < 400:
            return
        message = response.text or f"Request failed with code: {response.status_code}"
        if response.status_code == 401:
            raise AuthException(message, status_code=401)
        if response.status_code == 429:
            raise RateLimitException(message, status_code=429)
        raise YuAIException(message, status_code=response.status_code)
