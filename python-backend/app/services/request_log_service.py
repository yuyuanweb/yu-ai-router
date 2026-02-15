"""Request log service."""

from __future__ import annotations

from datetime import datetime

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import DEFAULT_LOG_LIMIT, REQUEST_STATUS_SUCCESS
from app.models.request_log import RequestLog
from app.schemas.stats import RequestLogVO
from app.services.api_key_service import ApiKeyService


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
    ) -> None:
        entity = RequestLog(
            user_id=user_id,
            api_key_id=api_key_id,
            model_name=model_name,
            prompt_tokens=prompt_tokens,
            completion_tokens=completion_tokens,
            total_tokens=total_tokens,
            duration=duration,
            status=status,
            error_message=error_message,
            create_time=datetime.utcnow(),
            update_time=datetime.utcnow(),
        )
        self.db.add(entity)
        await self.db.commit()
        if status == REQUEST_STATUS_SUCCESS and api_key_id and total_tokens > 0:
            await ApiKeyService(self.db).update_usage_stats(api_key_id, total_tokens)

    async def list_user_logs(self, user_id: int, limit: int | None) -> list[RequestLogVO]:
        safe_limit = limit if limit and limit > 0 else DEFAULT_LOG_LIMIT
        stmt = (
            select(RequestLog)
            .where(RequestLog.user_id == user_id)
            .order_by(RequestLog.create_time.desc())
            .limit(safe_limit)
        )
        rows = (await self.db.scalars(stmt)).all()
        return [RequestLogVO.model_validate(item) for item in rows]

    async def count_user_tokens(self, user_id: int) -> int:
        stmt = select(func.sum(RequestLog.total_tokens)).where(
            RequestLog.user_id == user_id,
            RequestLog.status == REQUEST_STATUS_SUCCESS,
        )
        total = await self.db.scalar(stmt)
        return int(total or 0)
