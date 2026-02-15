"""Stats API."""

from __future__ import annotations

from datetime import date, datetime, time, timedelta

from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.core.constants import ErrorCode, UserRole
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.middleware.auth import require_login, require_role
from app.models.request_log import RequestLog
from app.models.user import User
from app.schemas.common import BaseResponse, PageData
from app.schemas.stats import (
    CostStatsVO,
    RequestLogQueryRequest,
    RequestLogVO,
    TokenStatsVO,
    UserSummaryStatsVO,
)
from app.services.billing_service import BillingService
from app.services.quota_service import QuotaService
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
    return success([RequestLogVO.model_validate(item) for item in logs])


@router.get("/my/cost", response_model=BaseResponse[CostStatsVO])
async def get_my_cost_stats(
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[CostStatsVO]:
    billing = BillingService(db)
    return success(
        CostStatsVO(
            totalCost=await billing.get_user_total_cost(login_user.id),
            todayCost=await billing.get_user_today_cost(login_user.id),
        )
    )


@router.get("/my/summary", response_model=BaseResponse[UserSummaryStatsVO])
async def get_my_summary_stats(
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[UserSummaryStatsVO]:
    request_log_service = RequestLogService(db)
    billing_service = BillingService(db)
    quota_service = QuotaService(db)
    return success(
        UserSummaryStatsVO(
            totalTokens=await request_log_service.count_user_tokens(login_user.id),
            tokenQuota=login_user.token_quota,
            usedTokens=login_user.used_tokens or 0,
            remainingQuota=await quota_service.get_remaining_quota(login_user.id),
            totalCost=await billing_service.get_user_total_cost(login_user.id),
            todayCost=await billing_service.get_user_today_cost(login_user.id),
            totalRequests=await request_log_service.count_user_requests(login_user.id),
            successRequests=await request_log_service.count_user_success_requests(login_user.id),
        )
    )


@router.get("/my/daily", response_model=BaseResponse[list[dict[str, object]]])
async def get_my_daily_stats(
    start_date: str | None = None,
    end_date: str | None = None,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[list[dict[str, object]]]:
    end = date.fromisoformat(end_date) if end_date else date.today()
    start = date.fromisoformat(start_date) if start_date else end - timedelta(days=6)
    data = await RequestLogService(db).get_user_daily_stats(login_user.id, start, end)
    return success(data)


@router.post("/history/my/page", response_model=BaseResponse[PageData[RequestLogVO]])
async def page_my_history(
    payload: RequestLogQueryRequest,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[PageData[RequestLogVO]]:
    logs, total = await RequestLogService(db).page_by_query(
        page_num=payload.page_num,
        page_size=payload.page_size,
        user_id=login_user.id,
        request_model=payload.request_model,
        request_type=payload.request_type,
        source=payload.source,
        status=payload.status,
        start_at=datetime.combine(date.fromisoformat(payload.start_date), time.min)
        if payload.start_date
        else None,
        end_at=datetime.combine(date.fromisoformat(payload.end_date), time.max) if payload.end_date else None,
    )
    total_page = (total + payload.page_size - 1) // payload.page_size if total else 0
    return success(
        PageData[RequestLogVO](
            records=[RequestLogVO.model_validate(item) for item in logs],
            pageNumber=payload.page_num,
            pageSize=payload.page_size,
            totalPage=total_page,
            totalRow=total,
            optimizeCountQuery=True,
        )
    )


@router.get("/history/detail", response_model=BaseResponse[RequestLogVO])
async def get_history_detail(
    id: int,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[RequestLogVO]:
    if id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "请求参数错误")
    request_log = await RequestLogService(db).get_by_id(id)
    if request_log is None:
        raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "请求数据不存在")
    if login_user.user_role != UserRole.ADMIN.value and request_log.user_id != login_user.id:
        raise BusinessException(ErrorCode.NO_AUTH_ERROR, "只能查看自己的调用历史")
    return success(RequestLogVO.model_validate(request_log))


@router.post("/history/page", response_model=BaseResponse[PageData[RequestLogVO]])
async def page_history(
    payload: RequestLogQueryRequest,
    _: User = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[PageData[RequestLogVO]]:
    logs, total = await RequestLogService(db).page_by_query(
        page_num=payload.page_num,
        page_size=payload.page_size,
        user_id=payload.user_id,
        request_model=payload.request_model,
        request_type=payload.request_type,
        source=payload.source,
        status=payload.status,
        start_at=datetime.combine(date.fromisoformat(payload.start_date), time.min)
        if payload.start_date
        else None,
        end_at=datetime.combine(date.fromisoformat(payload.end_date), time.max) if payload.end_date else None,
    )
    total_page = (total + payload.page_size - 1) // payload.page_size if total else 0
    return success(
        PageData[RequestLogVO](
            records=[RequestLogVO.model_validate(item) for item in logs],
            pageNumber=payload.page_num,
            pageSize=payload.page_size,
            totalPage=total_page,
            totalRow=total,
            optimizeCountQuery=True,
        )
    )
