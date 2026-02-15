"""Request log service."""

from __future__ import annotations

from datetime import date, datetime, time, timedelta
from decimal import Decimal

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import DEFAULT_LOG_LIMIT, REQUEST_STATUS_SUCCESS
from app.models.request_log import RequestLog
from app.services.api_key_service import ApiKeyService
from app.services.billing_service import BillingService


class RequestLogService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db

    async def log_request(
        self,
        user_id: int | None,
        api_key_id: int | None,
        model_name: str,
        prompt_tokens: int,
        completion_tokens: int,
        total_tokens: int,
        duration: int,
        status: str,
        error_message: str | None,
        trace_id: str | None = None,
        model_id: int | None = None,
        request_model: str | None = None,
        request_type: str = "chat",
        source: str = "web",
        error_code: str | None = None,
        routing_strategy: str | None = None,
        is_fallback: bool = False,
        client_ip: str | None = None,
        user_agent: str | None = None,
        cost: Decimal | None = None,
    ) -> RequestLog:
        final_cost = cost
        if final_cost is None and status == REQUEST_STATUS_SUCCESS and model_id is not None:
            final_cost = await BillingService(self.db).calculate_cost(
                model_id=model_id,
                prompt_tokens=prompt_tokens,
                completion_tokens=completion_tokens,
            )
        entity = RequestLog(
            trace_id=trace_id,
            user_id=user_id,
            api_key_id=api_key_id,
            model_id=model_id,
            request_model=request_model,
            model_name=model_name,
            request_type=request_type,
            source=source,
            prompt_tokens=prompt_tokens,
            completion_tokens=completion_tokens,
            total_tokens=total_tokens,
            cost=final_cost or Decimal("0"),
            duration=duration,
            status=status,
            error_message=error_message,
            error_code=error_code,
            routing_strategy=routing_strategy,
            is_fallback=1 if is_fallback else 0,
            client_ip=client_ip,
            user_agent=user_agent,
            create_time=datetime.utcnow(),
            update_time=datetime.utcnow(),
        )
        self.db.add(entity)
        await self.db.commit()
        await self.db.refresh(entity)
        if status == REQUEST_STATUS_SUCCESS and api_key_id and total_tokens > 0:
            await ApiKeyService(self.db).update_usage_stats(api_key_id, total_tokens)
        return entity

    async def list_user_logs(self, user_id: int, limit: int | None) -> list[RequestLog]:
        safe_limit = limit if limit and limit > 0 else DEFAULT_LOG_LIMIT
        stmt = (
            select(RequestLog)
            .where(RequestLog.user_id == user_id)
            .order_by(RequestLog.create_time.desc())
            .limit(safe_limit)
        )
        return (await self.db.scalars(stmt)).all()

    async def count_user_tokens(self, user_id: int) -> int:
        stmt = select(func.sum(RequestLog.total_tokens)).where(
            RequestLog.user_id == user_id,
            RequestLog.status == REQUEST_STATUS_SUCCESS,
        )
        total = await self.db.scalar(stmt)
        return int(total or 0)

    async def count_user_requests(self, user_id: int) -> int:
        stmt = select(func.count()).select_from(RequestLog).where(RequestLog.user_id == user_id)
        total = await self.db.scalar(stmt)
        return int(total or 0)

    async def count_user_success_requests(self, user_id: int) -> int:
        stmt = select(func.count()).select_from(RequestLog).where(
            RequestLog.user_id == user_id,
            RequestLog.status == REQUEST_STATUS_SUCCESS,
        )
        total = await self.db.scalar(stmt)
        return int(total or 0)

    async def get_user_daily_stats(self, user_id: int, start_date: date, end_date: date) -> list[dict[str, object]]:
        result: list[dict[str, object]] = []
        current = start_date
        while current <= end_date:
            start_at = datetime.combine(current, time.min)
            end_at = datetime.combine(current, time.max)
            stmt = select(RequestLog).where(
                RequestLog.user_id == user_id,
                RequestLog.create_time >= start_at,
                RequestLog.create_time <= end_at,
            )
            logs = (await self.db.scalars(stmt)).all()
            request_count = len(logs)
            success_logs = [item for item in logs if item.status == REQUEST_STATUS_SUCCESS]
            result.append(
                {
                    "date": current.isoformat(),
                    "totalTokens": sum(item.total_tokens or 0 for item in success_logs),
                    "requestCount": request_count,
                    "successCount": len(success_logs),
                    "totalCost": sum((item.cost or Decimal("0")) for item in success_logs),
                }
            )
            current = current + timedelta(days=1)
        return result

    async def page_by_query(
        self,
        *,
        page_num: int,
        page_size: int,
        user_id: int | None = None,
        request_model: str | None = None,
        request_type: str | None = None,
        source: str | None = None,
        status: str | None = None,
        start_at: datetime | None = None,
        end_at: datetime | None = None,
    ) -> tuple[list[RequestLog], int]:
        stmt = select(RequestLog)
        if user_id is not None:
            stmt = stmt.where(RequestLog.user_id == user_id)
        if request_model:
            stmt = stmt.where(RequestLog.request_model.like(f"%{request_model}%"))
        if request_type:
            stmt = stmt.where(RequestLog.request_type == request_type)
        if source:
            stmt = stmt.where(RequestLog.source == source)
        if status:
            stmt = stmt.where(RequestLog.status == status)
        if start_at:
            stmt = stmt.where(RequestLog.create_time >= start_at)
        if end_at:
            stmt = stmt.where(RequestLog.create_time <= end_at)
        stmt = stmt.order_by(RequestLog.create_time.desc())
        total = await self.db.scalar(select(func.count()).select_from(stmt.subquery())) or 0
        logs = (await self.db.scalars(stmt.offset((page_num - 1) * page_size).limit(page_size))).all()
        return logs, int(total)

    async def get_by_id(self, log_id: int) -> RequestLog | None:
        return await self.db.scalar(select(RequestLog).where(RequestLog.id == log_id))
