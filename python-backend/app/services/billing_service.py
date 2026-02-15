"""Billing service."""

from __future__ import annotations

from decimal import Decimal
from datetime import datetime, time

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.models.model import Model
from app.models.request_log import RequestLog

TOKENS_PER_UNIT = Decimal("1000")


class BillingService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db

    @staticmethod
    def calculate_cost_from_model(model: Model | None, prompt_tokens: int, completion_tokens: int) -> Decimal:
        if model is None:
            return Decimal("0")
        input_price = Decimal(model.input_price or 0)
        output_price = Decimal(model.output_price or 0)
        input_cost = (input_price * Decimal(prompt_tokens)) / TOKENS_PER_UNIT
        output_cost = (output_price * Decimal(completion_tokens)) / TOKENS_PER_UNIT
        return (input_cost + output_cost).quantize(Decimal("0.000001"))

    async def calculate_cost(self, model_id: int | None, prompt_tokens: int, completion_tokens: int) -> Decimal:
        if model_id is None:
            return Decimal("0")
        model = await self.db.scalar(select(Model).where(Model.id == model_id, Model.is_delete == 0))
        return self.calculate_cost_from_model(model, prompt_tokens, completion_tokens)

    async def get_user_total_cost(self, user_id: int | None) -> Decimal:
        if user_id is None:
            return Decimal("0")
        stmt = select(func.sum(RequestLog.cost)).where(
            RequestLog.user_id == user_id,
            RequestLog.status == "success",
        )
        total = await self.db.scalar(stmt)
        return Decimal(total or 0)

    async def get_user_today_cost(self, user_id: int | None) -> Decimal:
        if user_id is None:
            return Decimal("0")
        now = datetime.now()
        start_at = datetime.combine(now.date(), time.min)
        end_at = datetime.combine(now.date(), time.max)
        stmt = select(func.sum(RequestLog.cost)).where(
            RequestLog.user_id == user_id,
            RequestLog.status == "success",
            RequestLog.create_time >= start_at,
            RequestLog.create_time <= end_at,
        )
        total = await self.db.scalar(stmt)
        return Decimal(total or 0)
