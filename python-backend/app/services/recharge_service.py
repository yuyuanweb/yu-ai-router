"""Recharge service."""

from __future__ import annotations

from datetime import datetime
from decimal import Decimal
from math import ceil

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import (
    DEFAULT_PAGE_NUM,
    DEFAULT_PAGE_SIZE,
    MAX_PAGE_SIZE,
    RECHARGE_STATUS_PENDING,
    RECHARGE_STATUS_SUCCESS,
)
from app.exceptions.business_exception import BusinessException
from app.core.constants import ErrorCode
from app.models.recharge_record import RechargeRecord
from app.schemas.common import PageData
from app.schemas.payment import RechargeRecordVO
from app.services.balance_service import BalanceService


class RechargeService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db
        self.balance_service = BalanceService(db)

    async def create_recharge_record(self, user_id: int, amount: Decimal, payment_method: str) -> RechargeRecord:
        entity = RechargeRecord(
            user_id=user_id,
            amount=amount,
            payment_method=payment_method,
            status=RECHARGE_STATUS_PENDING,
            description="账户充值",
            create_time=datetime.utcnow(),
            update_time=datetime.utcnow(),
        )
        self.db.add(entity)
        await self.db.commit()
        await self.db.refresh(entity)
        return entity

    async def update_recharge_status(self, record_id: int, status: str, payment_id: str | None) -> bool:
        entity = await self.db.scalar(select(RechargeRecord).where(RechargeRecord.id == record_id))
        if entity is None:
            raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "充值记录不存在")
        entity.status = status
        entity.payment_id = payment_id
        entity.update_time = datetime.utcnow()
        await self.db.commit()
        return True

    async def get_by_payment_id(self, payment_id: str) -> RechargeRecord | None:
        return await self.db.scalar(select(RechargeRecord).where(RechargeRecord.payment_id == payment_id))

    async def list_user_recharge_records(self, user_id: int, page_num: int, page_size: int) -> PageData[RechargeRecordVO]:
        safe_page_num = page_num if page_num > 0 else DEFAULT_PAGE_NUM
        safe_page_size = page_size if page_size > 0 else DEFAULT_PAGE_SIZE
        safe_page_size = min(safe_page_size, MAX_PAGE_SIZE)
        base_stmt = select(RechargeRecord).where(RechargeRecord.user_id == user_id)
        total_row = await self.db.scalar(select(func.count()).select_from(base_stmt.subquery())) or 0
        rows = (
            await self.db.scalars(
                base_stmt.order_by(RechargeRecord.create_time.desc())
                .offset((safe_page_num - 1) * safe_page_size)
                .limit(safe_page_size)
            )
        ).all()
        return PageData[RechargeRecordVO](
            records=[RechargeRecordVO.model_validate(item) for item in rows],
            pageNumber=safe_page_num,
            pageSize=safe_page_size,
            totalPage=ceil(total_row / safe_page_size) if total_row else 0,
            totalRow=total_row,
            optimizeCountQuery=True,
        )

    async def complete_recharge(self, record_id: int, payment_id: str) -> bool:
        entity = await self.db.scalar(select(RechargeRecord).where(RechargeRecord.id == record_id))
        if entity is None:
            raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "充值记录不存在")
        if entity.status == RECHARGE_STATUS_SUCCESS:
            return True
        entity.status = RECHARGE_STATUS_SUCCESS
        entity.payment_id = payment_id
        entity.update_time = datetime.utcnow()
        await self.db.commit()
        return await self.balance_service.add_balance(entity.user_id, entity.amount, f"充值：{payment_id}")
