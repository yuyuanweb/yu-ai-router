"""Periodic health check task."""

from __future__ import annotations

import asyncio
import logging

from app.core.constants import HEALTH_CHECK_INTERVAL_SECONDS
from app.db.session import session_maker
from app.services.health_check_service import HealthCheckService

logger = logging.getLogger("app")


class HealthCheckTask:
    async def execute_health_check(self) -> None:
        logger.debug("执行定时健康检查任务")
        try:
            async with session_maker() as db:
                await HealthCheckService(db).check_all_providers()
        except Exception:
            logger.exception("健康检查任务执行失败")

    async def run_loop(self, stop_event: asyncio.Event) -> None:
        while not stop_event.is_set():
            await self.execute_health_check()
            try:
                await asyncio.wait_for(stop_event.wait(), timeout=HEALTH_CHECK_INTERVAL_SECONDS)
            except TimeoutError:
                continue
