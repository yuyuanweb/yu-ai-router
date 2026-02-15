"""Stats API."""

from __future__ import annotations

from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.db.session import get_db_session
from app.middleware.auth import require_login
from app.models.user import User
from app.schemas.common import BaseResponse
from app.schemas.stats import RequestLogVO, TokenStatsVO
from app.services.request_log_service import RequestLogService

router = APIRouter(prefix="/stats", tags=["stats"])


@router.get("/my/tokens", response_model=BaseResponse[TokenStatsVO])
async def get_my_token_stats(
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[TokenStatsVO]:
    total_tokens = await RequestLogService(db).count_user_tokens(login_user.id)
    return success(TokenStatsVO(totalTokens=total_tokens))


@router.get("/my/logs", response_model=BaseResponse[list[RequestLogVO]])
async def get_my_logs(
    limit: int = 100,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[list[RequestLogVO]]:
    logs = await RequestLogService(db).list_user_logs(login_user.id, limit)
    return success(logs)
