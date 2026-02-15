"""Public chat API."""

from __future__ import annotations

from fastapi import APIRouter, Depends, Header, Request
from fastapi.responses import StreamingResponse
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import ErrorCode
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.schemas.chat import ChatRequest, ChatResponse
from app.services.api_key_service import ApiKeyService
from app.services.chat_service import ChatService

router = APIRouter(prefix="/v1/chat", tags=["chat"])


@router.post("/completions", response_model=ChatResponse | None)
async def chat_completions(
    payload: ChatRequest,
    request: Request,
    authorization: str | None = Header(default=None, alias="Authorization"),
    db: AsyncSession = Depends(get_db_session),
):
    if not authorization or not authorization.startswith("Bearer "):
        raise BusinessException(ErrorCode.NO_AUTH_ERROR, "缺少或无效的 Authorization Header")
    api_key_value = authorization[7:]
    api_key = await ApiKeyService(db).get_by_key_value(api_key_value)
    if api_key is None:
        raise BusinessException(ErrorCode.NO_AUTH_ERROR, "API Key 无效或已失效")
    if not payload.messages:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "messages 不能为空")
    if not payload.model:
        payload.model = "qwen-plus"
    service = ChatService(db)
    client_ip = request.client.host if request.client else None
    user_agent = request.headers.get("user-agent")
    if payload.stream:
        return StreamingResponse(
            service.chat_stream(payload, api_key.user_id, api_key.id, client_ip, user_agent),
            media_type="text/event-stream",
        )
    return await service.chat(payload, api_key.user_id, api_key.id, client_ip, user_agent)
