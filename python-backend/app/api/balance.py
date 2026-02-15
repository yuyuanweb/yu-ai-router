"""Balance API."""

from __future__ import annotations

from fastapi import APIRouter, Depends, Query
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.db.session import get_db_session
from app.middleware.auth import require_login
from app.models.user import User
from app.schemas.common import BaseResponse, PageData
from app.schemas.payment import BalanceVO, BillingRecordVO
from app.services.balance_service import BalanceService
from app.services.billing_record_service import BillingRecordService

router = APIRouter(prefix="/balance", tags=["balance"])


@router.get("/my", response_model=BaseResponse[BalanceVO])
async def get_my_balance(
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[BalanceVO]:
    balance_service = BalanceService(db)
    billing_record_service = BillingRecordService(db)
    balance = await balance_service.get_user_balance(login_user.id)
    total_spending = await billing_record_service.get_user_total_spending(login_user.id)
    total_recharge = await billing_record_service.get_user_total_recharge(login_user.id)
    return success(
        BalanceVO(
            balance=balance,
            totalSpending=total_spending,
            totalRecharge=total_recharge,
        )
    )


@router.get("/billing/my", response_model=BaseResponse[PageData[BillingRecordVO]])
async def get_my_billing_records(
    page_num: int = Query(default=1, alias="pageNum"),
    page_size: int = Query(default=10, alias="pageSize"),
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[PageData[BillingRecordVO]]:
    data = await BillingRecordService(db).list_user_billing_records(login_user.id, page_num, page_size)
    return success(data)
