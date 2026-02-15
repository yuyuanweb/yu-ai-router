"""Internal chat API."""

from __future__ import annotations

from fastapi import APIRouter, Depends, Query
from fastapi.responses import StreamingResponse
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.core.constants import API_KEY_STATUS_ACTIVE, ErrorCode
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.middleware.auth import require_login
from app.models.user import User
from app.schemas.chat import ChatRequest, ChatResponse
from app.schemas.common import BaseResponse
from app.services.api_key_service import ApiKeyService
from app.services.chat_service import ChatService

router = APIRouter(prefix="/internal/chat", tags=["internal-chat"])


@router.post("/completions", response_model=BaseResponse[ChatResponse] | None)
async def internal_chat_completions(
    payload: ChatRequest,
    api_key_id: int = Query(alias="apiKeyId"),
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
):
    api_key = await ApiKeyService(db).get_by_id(api_key_id)
    if api_key is None:
        raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "API Key 不存在")
    if api_key.user_id != login_user.id:
        raise BusinessException(ErrorCode.NO_AUTH_ERROR, "无权使用该 API Key")
    if api_key.status != API_KEY_STATUS_ACTIVE:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "API Key 已失效")
    if not payload.messages:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "messages 不能为空")
    if not payload.model:
        payload.model = "qwen-plus"
    service = ChatService(db)
    if payload.stream:
        return StreamingResponse(service.chat_stream(payload, login_user.id, api_key.id), media_type="text/event-stream")
    response = await service.chat(payload, login_user.id, api_key.id)
    return success(response)
