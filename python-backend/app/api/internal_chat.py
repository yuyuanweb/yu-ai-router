"""Internal chat API."""

from __future__ import annotations

import json

from fastapi import APIRouter, Depends, Request
from fastapi.responses import StreamingResponse
from redis.asyncio import Redis
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.core.constants import ErrorCode, INTERNAL_CHAT_IP_RATE_LIMIT_PER_SECOND
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.infra.redis_client import get_redis_client
from app.middleware.auth import require_login
from app.models.user import User
from app.schemas.chat import ChatMessage, ChatRequest, ChatResponse
from app.schemas.common import BaseResponse
from app.services.chat_service import ChatService
from app.services.rate_limit_service import RateLimitService
from app.utils.request import get_client_ip

router = APIRouter(prefix="/internal/chat", tags=["internal-chat"])


@router.post("/completions", response_model=BaseResponse[ChatResponse] | None)
async def internal_chat_completions(
    payload: ChatRequest,
    request: Request,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
    redis: Redis = Depends(get_redis_client),
):
    if not payload.messages:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "messages 不能为空")
    if not payload.model:
        payload.model = "qwen-plus"
    service = ChatService(db)
    client_ip = get_client_ip(request)
    allowed = await RateLimitService(redis).check_ip_rate_limit(
        client_ip, INTERNAL_CHAT_IP_RATE_LIMIT_PER_SECOND
    )
    if not allowed:
        raise BusinessException(ErrorCode.TOO_MANY_REQUEST, "请求过于频繁，请稍后再试")
    user_agent = request.headers.get("user-agent")
    if payload.stream:
        return StreamingResponse(
            service.chat_stream(payload, login_user.id, None, client_ip, user_agent),
            media_type="text/event-stream",
        )
    response = await service.chat(payload, login_user.id, None, client_ip, user_agent)
    return success(response)


@router.post("/completions/upload", response_model=BaseResponse[ChatResponse] | None)
async def internal_chat_completions_upload(
    request: Request,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
    redis: Redis = Depends(get_redis_client),
):
    try:
        form_data = await request.form()
    except RuntimeError as exc:
        raise BusinessException(
            ErrorCode.OPERATION_ERROR,
            '服务缺少 "python-multipart" 依赖，无法处理文件上传',
        ) from exc

    file = form_data.get("file")
    messages = form_data.get("messages")
    model = form_data.get("model")
    stream_raw = str(form_data.get("stream", "false")).lower()
    stream = stream_raw in {"1", "true", "yes", "on"}
    routing_strategy = form_data.get("routing_strategy")
    plugin_key = form_data.get("plugin_key")
    enable_reasoning_raw = form_data.get("enable_reasoning")
    enable_reasoning: bool | None
    if enable_reasoning_raw is None:
        enable_reasoning = None
    else:
        enable_reasoning = str(enable_reasoning_raw).lower() in {"1", "true", "yes", "on"}

    try:
        if not messages:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "messages 不能为空")
        parsed_messages = json.loads(messages)
        if not isinstance(parsed_messages, list) or not parsed_messages:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "messages 不能为空")
        message_list = [ChatMessage.model_validate(item) for item in parsed_messages]
    except BusinessException:
        raise
    except Exception as exc:
        raise BusinessException(ErrorCode.PARAMS_ERROR, f"messages 解析失败: {exc}") from exc

    payload = ChatRequest(
        messages=message_list,
        model=model,
        stream=stream,
        routing_strategy=routing_strategy,
        plugin_key=plugin_key,
        enable_reasoning=enable_reasoning,
    )

    if file is not None and hasattr(file, "read"):
        file_bytes = await file.read()
        payload.file_bytes = file_bytes
        content_type = getattr(file, "content_type", None)
        payload.file_type = content_type
        if not payload.plugin_key and content_type:
            if content_type.startswith("image/"):
                payload.plugin_key = "image_recognition"
            elif content_type == "application/pdf":
                payload.plugin_key = "pdf_parser"

    if not payload.messages:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "messages 不能为空")
    if not payload.model:
        payload.model = "qwen-plus"

    service = ChatService(db)
    client_ip = get_client_ip(request)
    allowed = await RateLimitService(redis).check_ip_rate_limit(
        client_ip, INTERNAL_CHAT_IP_RATE_LIMIT_PER_SECOND
    )
    if not allowed:
        raise BusinessException(ErrorCode.TOO_MANY_REQUEST, "请求过于频繁，请稍后再试")
    user_agent = request.headers.get("user-agent")
    if payload.stream:
        return StreamingResponse(
            service.chat_stream(payload, login_user.id, None, client_ip, user_agent),
            media_type="text/event-stream",
        )
    response = await service.chat(payload, login_user.id, None, client_ip, user_agent)
    return success(response)
