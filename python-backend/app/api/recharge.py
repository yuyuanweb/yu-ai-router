"""Recharge API."""

from __future__ import annotations

from decimal import Decimal

from fastapi import APIRouter, Depends, Query
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.core.config import get_settings
from app.core.constants import ErrorCode
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.middleware.auth import require_login
from app.models.user import User
from app.schemas.common import BaseResponse, PageData
from app.schemas.payment import CreateRechargeRequest, CreateRechargeResponse, RechargeRecordVO
from app.services.recharge_service import RechargeService
from app.services.stripe_payment_service import StripePaymentService

router = APIRouter(prefix="/recharge", tags=["recharge"])


@router.post("/stripe/create", response_model=BaseResponse[CreateRechargeResponse])
async def create_stripe_recharge(
    payload: CreateRechargeRequest,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[CreateRechargeResponse]:
    if payload.amount <= Decimal("0"):
        raise BusinessException(ErrorCode.PARAMS_ERROR, "充值金额必须大于0")
    if payload.amount < Decimal("1") or payload.amount > Decimal("10000"):
        raise BusinessException(ErrorCode.PARAMS_ERROR, "充值金额必须在1-10000元之间")
    settings = get_settings()
    stripe_service = StripePaymentService(RechargeService(db))
    session = await stripe_service.create_checkout_session(
        user_id=login_user.id,
        amount=payload.amount,
        success_url=settings.stripe_success_url,
        cancel_url=settings.stripe_cancel_url,
    )
    return success(CreateRechargeResponse(checkoutUrl=session["url"], sessionId=session["id"]))


@router.get("/stripe/success", response_model=BaseResponse[str])
async def stripe_success(
    session_id: str,
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[str]:
    ok = await StripePaymentService(RechargeService(db)).handle_payment_success(session_id)
    if not ok:
        raise BusinessException(ErrorCode.OPERATION_ERROR, "充值处理失败")
    return success("充值成功！")


@router.get("/stripe/cancel", response_model=BaseResponse[str])
async def stripe_cancel() -> BaseResponse[str]:
    return success("您取消了充值")


@router.get("/list/my", response_model=BaseResponse[PageData[RechargeRecordVO]])
async def get_my_recharge_records(
    page_num: int = Query(default=1, alias="pageNum"),
    page_size: int = Query(default=10, alias="pageSize"),
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[PageData[RechargeRecordVO]]:
    data = await RechargeService(db).list_user_recharge_records(login_user.id, page_num, page_size)
    return success(data)
