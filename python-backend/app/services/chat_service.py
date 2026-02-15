"""Chat service."""

from __future__ import annotations

import logging
import time
import uuid
from collections.abc import AsyncGenerator

from langchain_core.messages import AIMessage, HumanMessage, SystemMessage
from langchain_openai import ChatOpenAI
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.config import get_settings
from app.core.constants import CHAT_OBJECT, ErrorCode, REQUEST_STATUS_FAILED, REQUEST_STATUS_SUCCESS
from app.exceptions.business_exception import BusinessException
from app.schemas.chat import ChatMessage, ChatRequest, ChatResponse
from app.services.request_log_service import RequestLogService

settings = get_settings()
logger = logging.getLogger("app")


class ChatService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db
        self.request_log_service = RequestLogService(db)

    async def chat(self, chat_request: ChatRequest, user_id: int, api_key_id: int) -> ChatResponse:
        start = time.perf_counter()
        model_name = chat_request.model or settings.ai_model
        try:
            llm = self._build_llm(chat_request, stream=False)
            messages = self._to_langchain_messages(chat_request)
            ai_response: AIMessage = await llm.ainvoke(messages)

            content = self._normalize_content(ai_response.content)
            token_usage = self._extract_token_usage(ai_response)
            prompt_tokens = token_usage["prompt_tokens"]
            completion_tokens = token_usage["completion_tokens"]
            total_tokens = token_usage["total_tokens"]
            finish_reason = self._extract_finish_reason(ai_response)

            response = ChatResponse(
                id=self._extract_response_id(ai_response) or uuid.uuid4().hex,
                object=CHAT_OBJECT,
                created=int(time.time()),
                model=model_name,
                choices=[
                    {
                        "index": 0,
                        "message": ChatMessage(role="assistant", content=content),
                        "finishReason": finish_reason,
                    }
                ],
                usage={
                    "promptTokens": prompt_tokens,
                    "completionTokens": completion_tokens,
                    "totalTokens": total_tokens,
                },
            )
            duration = int((time.perf_counter() - start) * 1000)
            await self.request_log_service.log_request(
                user_id=user_id,
                api_key_id=api_key_id,
                model_name=model_name,
                prompt_tokens=prompt_tokens,
                completion_tokens=completion_tokens,
                total_tokens=total_tokens,
                duration=duration,
                status=REQUEST_STATUS_SUCCESS,
                error_message=None,
            )
            return response
        except Exception as exc:
            duration = int((time.perf_counter() - start) * 1000)
            await self.request_log_service.log_request(
                user_id=user_id,
                api_key_id=api_key_id,
                model_name=model_name,
                prompt_tokens=0,
                completion_tokens=0,
                total_tokens=0,
                duration=duration,
                status=REQUEST_STATUS_FAILED,
                error_message=str(exc),
            )
            raise BusinessException(ErrorCode.SYSTEM_ERROR, f"调用模型失败: {exc}") from exc

    async def chat_stream(
        self,
        chat_request: ChatRequest,
        user_id: int,
        api_key_id: int,
    ) -> AsyncGenerator[str, None]:
        start = time.perf_counter()
        model_name = chat_request.model or settings.ai_model
        prompt_tokens = 0
        completion_tokens = 0
        try:
            llm = self._build_llm(chat_request, stream=True)
            messages = self._to_langchain_messages(chat_request)
            async for chunk in llm.astream(messages):
                usage = self._extract_token_usage(chunk)
                if usage["prompt_tokens"] > 0:
                    prompt_tokens = usage["prompt_tokens"]
                if usage["completion_tokens"] > 0:
                    completion_tokens = usage["completion_tokens"]
                text = self._normalize_content(chunk.content)
                if not text:
                    continue
                escaped_text = text.replace(chr(10), "\\n")
                yield f"data: {escaped_text}\n\n"

            yield "data: [DONE]\n\n"

            duration = int((time.perf_counter() - start) * 1000)
            await self.request_log_service.log_request(
                user_id=user_id,
                api_key_id=api_key_id,
                model_name=model_name,
                prompt_tokens=prompt_tokens,
                completion_tokens=completion_tokens,
                total_tokens=prompt_tokens + completion_tokens,
                duration=duration,
                status=REQUEST_STATUS_SUCCESS,
                error_message=None,
            )
        except Exception as exc:
            duration = int((time.perf_counter() - start) * 1000)
            await self.request_log_service.log_request(
                user_id=user_id,
                api_key_id=api_key_id,
                model_name=model_name,
                prompt_tokens=0,
                completion_tokens=0,
                total_tokens=0,
                duration=duration,
                status=REQUEST_STATUS_FAILED,
                error_message=str(exc),
            )
            # SSE 响应一旦开始，不能再抛业务异常，否则会触发 "response already started"。
            # 这里与 Java 的 doOnError 行为对齐：记录日志并结束流。
            logger.error("流式调用模型失败: %s", exc, exc_info=True)
            return

    def _build_llm(self, chat_request: ChatRequest, stream: bool) -> ChatOpenAI:
        if not settings.ai_api_key:
            raise BusinessException(ErrorCode.SYSTEM_ERROR, "AI_API_KEY 未配置")
        openai_base_url = self._resolve_openai_base_url()
        kwargs: dict = {
            "api_key": settings.ai_api_key,
            "base_url": openai_base_url,
            "model": chat_request.model or settings.ai_model,
            "timeout": settings.ai_timeout_seconds,
            "streaming": stream,
        }
        if chat_request.temperature is not None:
            kwargs["temperature"] = chat_request.temperature
        if chat_request.max_tokens is not None:
            kwargs["max_tokens"] = chat_request.max_tokens
        if stream:
            kwargs["model_kwargs"] = {"stream_options": {"include_usage": True}}
        return ChatOpenAI(**kwargs)

    @staticmethod
    def _resolve_openai_base_url() -> str:
        base_url = settings.ai_base_url.rstrip("/")
        path = (settings.ai_chat_completions_path or "").strip()
        suffix = "/chat/completions"
        if path.endswith(suffix):
            version_prefix = path[: -len(suffix)]
            if version_prefix:
                if not version_prefix.startswith("/"):
                    version_prefix = f"/{version_prefix}"
                if not base_url.endswith(version_prefix):
                    base_url = f"{base_url}{version_prefix}"
        # 兼容通义千问 OpenAI 模式：如果未带 /v1，自动补齐。
        if "dashscope.aliyuncs.com/compatible-mode" in base_url and not base_url.endswith("/v1"):
            base_url = f"{base_url}/v1"
        return base_url

    @staticmethod
    def _to_langchain_messages(chat_request: ChatRequest):
        messages = []
        for item in chat_request.messages:
            if item.role == "system":
                messages.append(SystemMessage(content=item.content))
            elif item.role == "assistant":
                messages.append(AIMessage(content=item.content))
            else:
                messages.append(HumanMessage(content=item.content))
        return messages

    @staticmethod
    def _normalize_content(content) -> str:
        if content is None:
            return ""
        if isinstance(content, str):
            return content
        if isinstance(content, list):
            text_parts: list[str] = []
            for block in content:
                if isinstance(block, dict):
                    if block.get("type") == "text" and block.get("text"):
                        text_parts.append(str(block["text"]))
                else:
                    text_parts.append(str(block))
            return "".join(text_parts)
        return str(content)

    @staticmethod
    def _extract_response_id(message: AIMessage) -> str | None:
        return getattr(message, "id", None)

    @staticmethod
    def _extract_finish_reason(message: AIMessage) -> str | None:
        metadata = getattr(message, "response_metadata", {}) or {}
        finish_reason = metadata.get("finish_reason")
        return str(finish_reason) if finish_reason is not None else None

    @staticmethod
    def _extract_token_usage(message: AIMessage) -> dict[str, int]:
        usage = {"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0}
        usage_metadata = getattr(message, "usage_metadata", None) or {}
        if usage_metadata:
            usage["prompt_tokens"] = int(usage_metadata.get("input_tokens", 0) or 0)
            usage["completion_tokens"] = int(usage_metadata.get("output_tokens", 0) or 0)
            usage["total_tokens"] = int(usage_metadata.get("total_tokens", 0) or 0)
            if usage["total_tokens"] == 0:
                usage["total_tokens"] = usage["prompt_tokens"] + usage["completion_tokens"]
            return usage

        response_metadata = getattr(message, "response_metadata", None) or {}
        token_usage = response_metadata.get("token_usage", {}) if isinstance(response_metadata, dict) else {}
        usage["prompt_tokens"] = int(token_usage.get("prompt_tokens", 0) or 0)
        usage["completion_tokens"] = int(token_usage.get("completion_tokens", 0) or 0)
        usage["total_tokens"] = int(token_usage.get("total_tokens", 0) or 0)
        if usage["total_tokens"] == 0:
            usage["total_tokens"] = usage["prompt_tokens"] + usage["completion_tokens"]
        return usage
