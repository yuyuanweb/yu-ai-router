"""Chat service."""

from __future__ import annotations

import logging
import time
import uuid
from collections.abc import AsyncGenerator

from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import (
    ErrorCode,
    MAX_FALLBACK_RETRIES,
    MODEL_TYPE_CHAT,
    REQUEST_STATUS_FAILED,
    REQUEST_STATUS_SUCCESS,
    ROUTING_STRATEGY_AUTO,
    ROUTING_STRATEGY_FIXED,
)
from app.exceptions.business_exception import BusinessException
from app.models.model import Model
from app.models.model_provider import ModelProvider
from app.schemas.chat import ChatRequest, ChatResponse
from app.services.model_invoke_service import ModelInvokeService
from app.services.model_provider_service import ModelProviderService
from app.services.request_log_service import RequestLogService
from app.services.routing_service import RoutingService

logger = logging.getLogger("app")


class ChatService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db
        self.routing_service = RoutingService(db)
        self.model_provider_service = ModelProviderService(db)
        self.request_log_service = RequestLogService(db)
        self.model_invoke_service = ModelInvokeService()

    async def chat(
        self,
        chat_request: ChatRequest,
        user_id: int,
        api_key_id: int | None,
        client_ip: str | None = None,
        user_agent: str | None = None,
    ) -> ChatResponse:
        start = time.perf_counter()
        trace_id = uuid.uuid4().hex
        strategy_type = self._determine_strategy_type(chat_request.routing_strategy, chat_request.model)
        requested_model = chat_request.model
        selected_model = await self.routing_service.select_model(strategy_type, MODEL_TYPE_CHAT, requested_model)
        if selected_model is None:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "没有可用的模型")
        fallback_models = await self.routing_service.get_fallback_models(
            strategy_type, MODEL_TYPE_CHAT, requested_model
        )
        try:
            return await self._invoke_with_fallback(
                selected_model=selected_model,
                fallback_models=fallback_models,
                chat_request=chat_request,
                user_id=user_id,
                api_key_id=api_key_id,
                trace_id=trace_id,
                start=start,
                strategy_type=strategy_type,
                client_ip=client_ip,
                user_agent=user_agent,
            )
        except Exception as exc:
            await self.request_log_service.log_request(
                trace_id=trace_id,
                user_id=user_id,
                api_key_id=api_key_id,
                model_name=requested_model or "",
                request_model=requested_model,
                prompt_tokens=0,
                completion_tokens=0,
                total_tokens=0,
                duration=int((time.perf_counter() - start) * 1000),
                status=REQUEST_STATUS_FAILED,
                error_message=str(exc),
                error_code="SYSTEM_ERROR",
                routing_strategy=strategy_type,
                is_fallback=False,
                source="api" if api_key_id else "web",
                client_ip=client_ip,
                user_agent=user_agent,
            )
            raise BusinessException(ErrorCode.SYSTEM_ERROR, f"调用模型失败: {exc}") from exc

    async def chat_stream(
        self,
        chat_request: ChatRequest,
        user_id: int,
        api_key_id: int | None,
        client_ip: str | None = None,
        user_agent: str | None = None,
    ) -> AsyncGenerator[str, None]:
        start = time.perf_counter()
        trace_id = uuid.uuid4().hex
        strategy_type = self._determine_strategy_type(chat_request.routing_strategy, chat_request.model)
        requested_model = chat_request.model
        prompt_tokens = 0
        completion_tokens = 0
        thinking_started = False
        thinking_ended = False
        try:
            model = await self.routing_service.select_model(strategy_type, MODEL_TYPE_CHAT, requested_model)
            if model is None:
                raise BusinessException(ErrorCode.PARAMS_ERROR, "没有可用的模型")
            provider = await self.model_provider_service.get_by_id(model.provider_id)
            if provider is None:
                raise BusinessException(ErrorCode.PARAMS_ERROR, "模型提供者不存在")
            async for chunk in self.model_invoke_service.invoke_stream_chunk(model, provider, chat_request):
                if chunk.prompt_tokens and chunk.prompt_tokens > 0:
                    prompt_tokens = chunk.prompt_tokens
                if chunk.completion_tokens and chunk.completion_tokens > 0:
                    completion_tokens = chunk.completion_tokens

                if bool(chat_request.enable_reasoning) and chunk.reasoning_content:
                    if not thinking_started:
                        thinking_started = True
                        yield f"data: [THINKING]{self._escape_newlines(chunk.reasoning_content)}\n\n"
                    elif not thinking_ended:
                        yield f"data: {self._escape_newlines(chunk.reasoning_content)}\n\n"
                    continue

                if bool(chat_request.enable_reasoning) and thinking_started and not thinking_ended:
                    thinking_ended = True
                    if chunk.text:
                        yield f"data: [/THINKING]\n\n"
                        yield f"data: {self._escape_newlines(chunk.text)}\n\n"
                    continue

                if chunk.text:
                    yield f"data: {self._escape_newlines(chunk.text)}\n\n"

            yield "data: [DONE]\n\n"
            await self.request_log_service.log_request(
                trace_id=trace_id,
                user_id=user_id,
                api_key_id=api_key_id,
                model_id=model.id,
                model_name=model.model_key,
                request_model=requested_model,
                prompt_tokens=prompt_tokens,
                completion_tokens=completion_tokens,
                total_tokens=prompt_tokens + completion_tokens,
                duration=int((time.perf_counter() - start) * 1000),
                status=REQUEST_STATUS_SUCCESS,
                error_message=None,
                routing_strategy=strategy_type,
                is_fallback=False,
                source="api" if api_key_id else "web",
                client_ip=client_ip,
                user_agent=user_agent,
            )
        except Exception as exc:
            await self.request_log_service.log_request(
                trace_id=trace_id,
                user_id=user_id,
                api_key_id=api_key_id,
                model_name=requested_model or "",
                request_model=requested_model,
                prompt_tokens=0,
                completion_tokens=0,
                total_tokens=0,
                duration=int((time.perf_counter() - start) * 1000),
                status=REQUEST_STATUS_FAILED,
                error_message=str(exc),
                error_code="STREAM_ERROR",
                routing_strategy=strategy_type,
                is_fallback=False,
                source="api" if api_key_id else "web",
                client_ip=client_ip,
                user_agent=user_agent,
            )
            logger.error("流式调用模型失败: %s", exc, exc_info=True)
            return

    async def _invoke_with_fallback(
        self,
        selected_model: Model,
        fallback_models: list[Model],
        chat_request: ChatRequest,
        user_id: int,
        api_key_id: int | None,
        trace_id: str,
        start: float,
        strategy_type: str,
        client_ip: str | None,
        user_agent: str | None,
    ) -> ChatResponse:
        try:
            return await self._call_single_model(
                model=selected_model,
                chat_request=chat_request,
                user_id=user_id,
                api_key_id=api_key_id,
                trace_id=trace_id,
                start=start,
                strategy_type=strategy_type,
                is_fallback=False,
                client_ip=client_ip,
                user_agent=user_agent,
            )
        except Exception:
            retries = min(len(fallback_models), MAX_FALLBACK_RETRIES)
            last_exc: Exception | None = None
            for index in range(retries):
                fallback_model = fallback_models[index]
                try:
                    return await self._call_single_model(
                        model=fallback_model,
                        chat_request=chat_request,
                        user_id=user_id,
                        api_key_id=api_key_id,
                        trace_id=trace_id,
                        start=start,
                        strategy_type=strategy_type,
                        is_fallback=True,
                        client_ip=client_ip,
                        user_agent=user_agent,
                    )
                except Exception as exc:  # noqa: PERF203
                    last_exc = exc
            if last_exc is not None:
                raise last_exc
            raise

    async def _call_single_model(
        self,
        model: Model,
        chat_request: ChatRequest,
        user_id: int,
        api_key_id: int | None,
        trace_id: str,
        start: float,
        strategy_type: str,
        is_fallback: bool,
        client_ip: str | None,
        user_agent: str | None,
    ) -> ChatResponse:
        provider = await self.model_provider_service.get_by_id(model.provider_id)
        if provider is None:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "模型提供者不存在")
        try:
            response = await self.model_invoke_service.invoke(model, provider, chat_request)
            usage = response.usage
            await self.request_log_service.log_request(
                trace_id=trace_id,
                user_id=user_id,
                api_key_id=api_key_id,
                model_id=model.id,
                model_name=model.model_key,
                request_model=chat_request.model,
                prompt_tokens=usage.prompt_tokens,
                completion_tokens=usage.completion_tokens,
                total_tokens=usage.total_tokens,
                duration=int((time.perf_counter() - start) * 1000),
                status=REQUEST_STATUS_SUCCESS,
                error_message=None,
                routing_strategy=strategy_type,
                is_fallback=is_fallback,
                source="api" if api_key_id else "web",
                client_ip=client_ip,
                user_agent=user_agent,
            )
            return response
        except Exception as exc:
            await self.request_log_service.log_request(
                trace_id=trace_id,
                user_id=user_id,
                api_key_id=api_key_id,
                model_id=model.id,
                model_name=model.model_key,
                request_model=chat_request.model,
                prompt_tokens=0,
                completion_tokens=0,
                total_tokens=0,
                duration=int((time.perf_counter() - start) * 1000),
                status=REQUEST_STATUS_FAILED,
                error_message=str(exc),
                error_code="MODEL_ERROR",
                routing_strategy=strategy_type,
                is_fallback=is_fallback,
                source="api" if api_key_id else "web",
                client_ip=client_ip,
                user_agent=user_agent,
            )
            raise

    @staticmethod
    def _determine_strategy_type(requested_strategy: str | None, requested_model: str | None) -> str:
        if requested_strategy:
            return requested_strategy
        if requested_model:
            return ROUTING_STRATEGY_FIXED
        return ROUTING_STRATEGY_AUTO

    @staticmethod
    def _escape_newlines(text: str) -> str:
        return text.replace("\n", "\\n") if text else ""
