"""Billing record service."""

from __future__ import annotations

from decimal import Decimal
from math import ceil

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import (
    BILLING_TYPE_API_CALL,
    BILLING_TYPE_RECHARGE,
    DEFAULT_PAGE_NUM,
    DEFAULT_PAGE_SIZE,
    MAX_PAGE_SIZE,
)
from app.models.billing_record import BillingRecord
from app.schemas.common import PageData
from app.schemas.payment import BillingRecordVO


class BillingRecordService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db

    async def list_user_billing_records(self, user_id: int, page_num: int, page_size: int) -> PageData[BillingRecordVO]:
        safe_page_num = page_num if page_num > 0 else DEFAULT_PAGE_NUM
        safe_page_size = page_size if page_size > 0 else DEFAULT_PAGE_SIZE
        safe_page_size = min(safe_page_size, MAX_PAGE_SIZE)
        base_stmt = select(BillingRecord).where(BillingRecord.user_id == user_id)
        total_row = await self.db.scalar(select(func.count()).select_from(base_stmt.subquery())) or 0
        rows = (
            await self.db.scalars(
                base_stmt.order_by(BillingRecord.create_time.desc())
                .offset((safe_page_num - 1) * safe_page_size)
                .limit(safe_page_size)
            )
        ).all()
        return PageData[BillingRecordVO](
            records=[BillingRecordVO.model_validate(item) for item in rows],
            pageNumber=safe_page_num,
            pageSize=safe_page_size,
            totalPage=ceil(total_row / safe_page_size) if total_row else 0,
            totalRow=total_row,
            optimizeCountQuery=True,
        )

    async def get_user_total_spending(self, user_id: int) -> Decimal:
        total = await self.db.scalar(
            select(func.sum(BillingRecord.amount)).where(
                BillingRecord.user_id == user_id,
                BillingRecord.billing_type == BILLING_TYPE_API_CALL,
            )
        )
        return Decimal(total or 0)

    async def get_user_total_recharge(self, user_id: int) -> Decimal:
        total = await self.db.scalar(
            select(func.sum(BillingRecord.amount)).where(
                BillingRecord.user_id == user_id,
                BillingRecord.billing_type == BILLING_TYPE_RECHARGE,
            )
        )
        return Decimal(total or 0)
