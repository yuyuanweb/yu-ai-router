"""OpenAI-compatible model adapter."""

from __future__ import annotations

import time
import uuid
from collections.abc import AsyncGenerator, Iterable

from langchain_core.messages import AIMessage, HumanMessage, SystemMessage
from langchain_openai import ChatOpenAI
from openai import AsyncOpenAI

from app.adapter.model_adapter import ModelAdapter
from app.core.constants import CHAT_OBJECT
from app.models.model import Model
from app.models.model_provider import ModelProvider
from app.schemas.chat import ChatMessage, ChatRequest, ChatResponse, StreamChunk


class OpenAIAdapter(ModelAdapter):
    supported_providers = {"openai", "gpt"}

    async def invoke(self, model: Model, provider: ModelProvider, chat_request: ChatRequest) -> ChatResponse:
        llm = self._build_llm(model, provider, chat_request, stream=False)
        message: AIMessage = await llm.ainvoke(self._to_messages(chat_request.messages))
        usage = self._extract_usage(message)
        return ChatResponse(
            id=getattr(message, "id", None) or uuid.uuid4().hex,
            object=CHAT_OBJECT,
            created=int(time.time()),
            model=model.model_key,
            choices=[
                {
                    "index": 0,
                    "message": ChatMessage(role="assistant", content=self._normalize_text(message.content)),
                    "finishReason": self._extract_finish_reason(message),
                }
            ],
            usage={
                "promptTokens": usage["prompt_tokens"],
                "completionTokens": usage["completion_tokens"],
                "totalTokens": usage["total_tokens"],
            },
        )

    async def invoke_stream_chunk(
        self, model: Model, provider: ModelProvider, chat_request: ChatRequest
    ) -> AsyncGenerator[StreamChunk, None]:
        client = AsyncOpenAI(
            api_key=provider.api_key,
            base_url=self._resolve_base_url(provider.base_url),
            timeout=model.default_timeout / 1000 if model.default_timeout else 60,
        )
        payload: dict = {
            "model": model.model_key,
            "messages": self._to_openai_messages(chat_request.messages),
            "stream": True,
            "stream_options": {"include_usage": True},
        }
        if chat_request.temperature is not None:
            payload["temperature"] = chat_request.temperature
        if chat_request.max_tokens is not None:
            payload["max_tokens"] = chat_request.max_tokens
        extra_body = self._build_reasoning_extra_body(model, provider, chat_request)
        if extra_body:
            payload["extra_body"] = extra_body

        stream = await client.chat.completions.create(**payload)
        async for chunk in stream:
            text = ""
            reasoning = ""
            prompt_tokens = None
            completion_tokens = None

            usage_obj = getattr(chunk, "usage", None)
            if usage_obj is not None:
                prompt_tokens = getattr(usage_obj, "prompt_tokens", None)
                completion_tokens = getattr(usage_obj, "completion_tokens", None)

            choices = getattr(chunk, "choices", None) or []
            if choices:
                delta = getattr(choices[0], "delta", None)
                if delta is not None:
                    text = getattr(delta, "content", "") or ""
                    reasoning = getattr(delta, "reasoning_content", "") or ""
                    # 兜底：兼容不同 SDK/供应商字段名
                    model_extra = getattr(delta, "model_extra", None) or {}
                    if not reasoning and isinstance(model_extra, dict):
                        reasoning = (
                            str(
                                model_extra.get("reasoning_content")
                                or model_extra.get("reasoningContent")
                                or model_extra.get("thinking")
                                or ""
                            )
                            or ""
                        )

            stream_chunk = StreamChunk(
                text=text or None,
                reasoningContent=reasoning or None,
                promptTokens=prompt_tokens or None,
                completionTokens=completion_tokens or None,
                empty=not bool(text or reasoning or prompt_tokens or completion_tokens),
            )
            if not stream_chunk.empty:
                yield stream_chunk

    def supports(self, provider_name: str) -> bool:
        return provider_name.lower() in self.supported_providers

    def _build_llm(
        self, model: Model, provider: ModelProvider, chat_request: ChatRequest, stream: bool
    ) -> ChatOpenAI:
        kwargs: dict = {
            "api_key": provider.api_key,
            "base_url": self._resolve_base_url(provider.base_url),
            "model": model.model_key,
            "streaming": stream,
            "timeout": model.default_timeout / 1000 if model.default_timeout else 60,
        }
        if chat_request.temperature is not None:
            kwargs["temperature"] = chat_request.temperature
        if chat_request.max_tokens is not None:
            kwargs["max_tokens"] = chat_request.max_tokens
        model_kwargs: dict = {}
        if stream:
            model_kwargs["stream_options"] = {"include_usage": True}
        extra_body = self._build_reasoning_extra_body(model, provider, chat_request)
        if extra_body:
            kwargs["extra_body"] = extra_body
        if model_kwargs:
            kwargs["model_kwargs"] = model_kwargs
        return ChatOpenAI(**kwargs)

    @staticmethod
    def _build_reasoning_extra_body(
        model: Model, provider: ModelProvider, chat_request: ChatRequest
    ) -> dict:
        if not (model.support_reasoning == 1 and chat_request.enable_reasoning):
            return {}
        provider_name = provider.provider_name.lower()
        if provider_name in {"qwen", "dashscope", "tongyi", "aliyun"}:
            return {"enable_thinking": True}
        if provider_name in {"zhipu", "zhipuai", "glm"}:
            return {"thinking": {"type": "enabled"}}
        return {}

    @staticmethod
    def _to_messages(items: Iterable[ChatMessage]):
        messages = []
        for item in items:
            if item.role == "system":
                messages.append(SystemMessage(content=item.content))
            elif item.role == "assistant":
                messages.append(AIMessage(content=item.content))
            else:
                messages.append(HumanMessage(content=item.content))
        return messages

    @staticmethod
    def _to_openai_messages(items: Iterable[ChatMessage]) -> list[dict[str, str]]:
        return [{"role": item.role, "content": item.content} for item in items]

    @staticmethod
    def _extract_usage(message: AIMessage) -> dict[str, int]:
        usage = {"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0}
        usage_metadata = getattr(message, "usage_metadata", None) or {}
        if usage_metadata:
            usage["prompt_tokens"] = int(usage_metadata.get("input_tokens", 0) or 0)
            usage["completion_tokens"] = int(usage_metadata.get("output_tokens", 0) or 0)
            usage["total_tokens"] = int(usage_metadata.get("total_tokens", 0) or 0)
        response_metadata = getattr(message, "response_metadata", None) or {}
        token_usage = response_metadata.get("token_usage", {}) if isinstance(response_metadata, dict) else {}
        usage["prompt_tokens"] = usage["prompt_tokens"] or int(token_usage.get("prompt_tokens", 0) or 0)
        usage["completion_tokens"] = usage["completion_tokens"] or int(token_usage.get("completion_tokens", 0) or 0)
        usage["total_tokens"] = usage["total_tokens"] or int(token_usage.get("total_tokens", 0) or 0)
        if usage["total_tokens"] == 0:
            usage["total_tokens"] = usage["prompt_tokens"] + usage["completion_tokens"]
        return usage

    @staticmethod
    def _extract_finish_reason(message: AIMessage) -> str | None:
        metadata = getattr(message, "response_metadata", {}) or {}
        value = metadata.get("finish_reason")
        return str(value) if value is not None else None

    @staticmethod
    def _normalize_text(content) -> str:
        if content is None:
            return ""
        if isinstance(content, str):
            return content
        if isinstance(content, list):
            parts: list[str] = []
            for item in content:
                if isinstance(item, dict) and item.get("type") == "text":
                    parts.append(str(item.get("text", "")))
                else:
                    parts.append(str(item))
            return "".join(parts)
        return str(content)

    @staticmethod
    def _split_text_and_reasoning(message: AIMessage) -> tuple[str, str]:
        content = getattr(message, "content", None)
        text_parts: list[str] = []
        reasoning_parts: list[str] = []

        if isinstance(content, str):
            text_parts.append(content)
        elif isinstance(content, list):
            for item in content:
                if not isinstance(item, dict):
                    continue
                item_type = str(item.get("type", "")).lower()
                # 常见 OpenAI-compatible 推理块：type=reasoning / thinking
                if item_type in {"reasoning", "thinking"}:
                    value = (
                        item.get("text")
                        or item.get("reasoning_content")
                        or item.get("reasoningContent")
                        or item.get("content")
                    )
                    if value:
                        reasoning_parts.append(str(value))
                    continue
                if item_type == "text":
                    value = item.get("text") or item.get("content")
                    if value:
                        text_parts.append(str(value))
                    continue
                # 兜底：部分兼容实现直接把 reasoning 字段放在块里
                reasoning_value = item.get("reasoning_content") or item.get("reasoningContent")
                if reasoning_value:
                    reasoning_parts.append(str(reasoning_value))
                text_value = item.get("text") or item.get("content")
                if text_value:
                    text_parts.append(str(text_value))

        # 兜底：从 additional_kwargs / response_metadata 深度递归提取
        additional = getattr(message, "additional_kwargs", None) or {}
        response_metadata = getattr(message, "response_metadata", None) or {}
        reasoning_parts.extend(OpenAIAdapter._collect_reasoning_strings(additional))
        reasoning_parts.extend(OpenAIAdapter._collect_reasoning_strings(response_metadata))

        # 去重并保持顺序
        seen: set[str] = set()
        dedup_reasoning: list[str] = []
        for item in reasoning_parts:
            if item and item not in seen:
                seen.add(item)
                dedup_reasoning.append(item)
        return "".join(text_parts), "".join(dedup_reasoning)

    @staticmethod
    def _extract_reasoning_text(message: AIMessage) -> str:
        additional = getattr(message, "additional_kwargs", None) or {}
        for key in ("reasoning_content", "reasoningContent"):
            value = additional.get(key)
            if isinstance(value, str) and value:
                return value
        if isinstance(additional.get("thinking"), dict):
            thinking_text = additional["thinking"].get("content")
            if isinstance(thinking_text, str) and thinking_text:
                return thinking_text
        response_metadata = getattr(message, "response_metadata", None) or {}
        if isinstance(response_metadata, dict):
            value = response_metadata.get("reasoning_content") or response_metadata.get("reasoningContent")
            if isinstance(value, str):
                return value
        return ""

    @staticmethod
    def _collect_reasoning_strings(data) -> list[str]:
        results: list[str] = []

        def walk(node) -> None:
            if isinstance(node, dict):
                for key, value in node.items():
                    lowered = str(key).lower()
                    if lowered in {"reasoning", "reasoning_content", "reasoningcontent", "thinking"}:
                        if isinstance(value, str) and value:
                            results.append(value)
                        else:
                            walk(value)
                    else:
                        walk(value)
                return
            if isinstance(node, list):
                for item in node:
                    walk(item)

        walk(data)
        return results

    @staticmethod
    def _resolve_base_url(base_url: str) -> str:
        normalized = base_url.rstrip("/")
        if "dashscope.aliyuncs.com/compatible-mode" in normalized and not normalized.endswith("/v1"):
            normalized = f"{normalized}/v1"
        if "api.deepseek.com" in normalized and not normalized.endswith("/v1"):
            normalized = f"{normalized}/v1"
        if "open.bigmodel.cn/api/paas" in normalized and not normalized.endswith("/v4"):
            normalized = f"{normalized}/v4"
        return normalized
