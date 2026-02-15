"""Internal chat API."""

from __future__ import annotations

from fastapi import APIRouter, Depends, Request
from fastapi.responses import StreamingResponse
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.core.constants import ErrorCode
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.middleware.auth import require_login
from app.models.user import User
from app.schemas.chat import ChatRequest, ChatResponse
from app.schemas.common import BaseResponse
from app.services.chat_service import ChatService

router = APIRouter(prefix="/internal/chat", tags=["internal-chat"])


@router.post("/completions", response_model=BaseResponse[ChatResponse] | None)
async def internal_chat_completions(
    payload: ChatRequest,
    request: Request,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
):
    if not payload.messages:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "messages 不能为空")
    if not payload.model:
        payload.model = "qwen-plus"
    service = ChatService(db)
    client_ip = request.client.host if request.client else None
    user_agent = request.headers.get("user-agent")
    if payload.stream:
        return StreamingResponse(
            service.chat_stream(payload, login_user.id, None, client_ip, user_agent),
            media_type="text/event-stream",
        )
    response = await service.chat(payload, login_user.id, None, client_ip, user_agent)
    return success(response)
