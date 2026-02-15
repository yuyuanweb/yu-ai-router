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
from app.schemas.chat import ChatRequest, ChatResponse, StreamChoice, StreamDelta, StreamResponse
from app.services.model_invoke_service import ModelInvokeService
from app.services.model_provider_service import ModelProviderService
from app.services.quota_service import QuotaService
from app.services.request_log_service import RequestLogService
from app.services.routing_service import RoutingService
from app.services.user_service import UserService

logger = logging.getLogger("app")


class ChatService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db
        self.routing_service = RoutingService(db)
        self.model_provider_service = ModelProviderService(db)
        self.request_log_service = RequestLogService(db)
        self.model_invoke_service = ModelInvokeService()
        self.quota_service = QuotaService(db)
        self.user_service = UserService(db)

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
        if user_id and await self.user_service.is_user_disabled(user_id):
            raise BusinessException(ErrorCode.FORBIDDEN_ERROR, "账号已被禁用，无法使用服务")
        if user_id and not await self.quota_service.check_quota(user_id):
            raise BusinessException(ErrorCode.OPERATION_ERROR, "Token配额已用尽，请联系管理员增加配额")
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
        if user_id and await self.user_service.is_user_disabled(user_id):
            raise BusinessException(ErrorCode.FORBIDDEN_ERROR, "账号已被禁用，无法使用服务")
        if user_id and not await self.quota_service.check_quota(user_id):
            raise BusinessException(ErrorCode.OPERATION_ERROR, "Token配额已用尽，请联系管理员增加配额")
        strategy_type = self._determine_strategy_type(chat_request.routing_strategy, chat_request.model)
        requested_model = chat_request.model
        prompt_tokens = 0
        completion_tokens = 0
        created = int(time.time())
        first_chunk = True
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

                if not chunk.text and not chunk.reasoning_content:
                    continue

                delta = StreamDelta(
                    role="assistant" if first_chunk else None,
                    content=chunk.text,
                    reasoningContent=chunk.reasoning_content,
                )
                first_chunk = False
                response = StreamResponse(
                    id=trace_id,
                    object="chat.completion.chunk",
                    created=created,
                    model=model.model_key,
                    choices=[StreamChoice(index=0, delta=delta, finishReason=None)],
                )
                yield f"data: {response.model_dump_json(by_alias=True)}\n\n"

            finish_response = StreamResponse(
                id=trace_id,
                object="chat.completion.chunk",
                created=created,
                model=model.model_key,
                choices=[
                    StreamChoice(
                        index=0,
                        delta=StreamDelta(),
                        finishReason="stop",
                    )
                ],
            )
            yield f"data: {finish_response.model_dump_json(by_alias=True)}\n\n"
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
            if user_id and prompt_tokens + completion_tokens > 0:
                await self.quota_service.deduct_tokens(user_id, prompt_tokens + completion_tokens)
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
            if user_id and usage.total_tokens > 0:
                await self.quota_service.deduct_tokens(user_id, usage.total_tokens)
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

