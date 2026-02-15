"""Health check service."""

from __future__ import annotations

import logging
from dataclasses import dataclass
from datetime import datetime, timedelta
from decimal import Decimal, ROUND_HALF_UP

import httpx
from sqlalchemy import case, func, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import (
    HEALTH_CHECK_MAX_HISTORY_SIZE,
    HEALTH_CHECK_TIMEOUT_MS,
    HEALTH_CHECK_STATS_HOURS,
    HEALTH_STATUS_DEGRADED,
    HEALTH_STATUS_HEALTHY,
    HEALTH_STATUS_UNHEALTHY,
    HEALTH_STATUS_UNKNOWN,
    PROVIDER_STATUS_ACTIVE,
    SCORE_COST_WEIGHT,
    SCORE_LATENCY_WEIGHT,
    SCORE_PRIORITY_WEIGHT,
    SCORE_SUCCESS_RATE_WEIGHT,
)
from app.models.request_log import RequestLog
from app.services.model_provider_service import ModelProviderService
from app.services.model_service import ModelService

logger = logging.getLogger("app")


@dataclass
class _ModelStats:
    avg_latency: int = 0
    success_rate: Decimal = Decimal("100")
    total_requests: int = 0


@dataclass
class _NormParams:
    min_cost: Decimal = Decimal("0")
    max_cost: Decimal = Decimal("1")
    min_latency: int = 0
    max_latency: int = 10000
    min_priority: int = 0
    max_priority: int = 100


class HealthCheckService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db
        self.model_provider_service = ModelProviderService(db)
        self.model_service = ModelService(db)
        self.health_history_map: dict[int, list[bool]] = {}

    async def check_provider_health(self, provider) -> bool:
        if provider is None:
            return False
        try:
            ok, latency = await self._send_health_check_request(provider.base_url, provider.api_key)
            self._record_health_check_result(provider.id, ok)
            success_rate = self._calculate_success_rate(provider.id)
            health_status = self._determine_health_status(success_rate)
            await self.model_provider_service.update_health_status(
                provider.id, health_status, latency if ok else 0, success_rate
            )
            return ok
        except Exception:
            self._record_health_check_result(provider.id, False)
            await self.model_provider_service.update_health_status(
                provider.id,
                HEALTH_STATUS_UNHEALTHY,
                0,
                self._calculate_success_rate(provider.id),
            )
            return False

    async def check_all_providers(self) -> None:
        providers = await self.model_provider_service.list_entities()
        if not providers:
            logger.warning("没有找到任何模型提供者")
            return
        logger.info("开始健康检查，共 %s 个提供者", len(providers))
        for provider in providers:
            if provider.status == PROVIDER_STATUS_ACTIVE:
                await self.check_provider_health(provider)
        await self.sync_model_metrics_from_request_log()
        logger.info("健康检查完成")

    async def sync_model_metrics_from_request_log(self) -> None:
        active_models = await self.model_service.list_active_entities()
        if not active_models:
            logger.warning("没有找到任何启用的模型")
            return

        model_stats = await self._query_model_stats_from_db(datetime.utcnow() - timedelta(hours=HEALTH_CHECK_STATS_HOURS))
        norm = self._calculate_normalization_params(active_models, model_stats)

        for model in active_models:
            stats = model_stats.get(model.id, _ModelStats())
            score = self._calculate_score(model, stats, norm)
            health_status = self._determine_health_status_by_success_rate(stats.success_rate)
            await self.model_service.update_model_metrics(
                model.id,
                health_status,
                stats.avg_latency,
                stats.success_rate,
                score,
            )

        await self._sync_provider_metrics_from_model_stats(active_models, model_stats)

    async def _query_model_stats_from_db(self, start_time: datetime) -> dict[int, _ModelStats]:
        success_case = case((RequestLog.status == "success", RequestLog.duration), else_=None)
        success_count_case = case((RequestLog.status == "success", 1), else_=0)
        stmt = (
            select(
                RequestLog.model_id,
                func.avg(success_case).label("avg_latency"),
                func.count().label("total_requests"),
                (func.sum(success_count_case) * 100.0 / func.nullif(func.count(), 0)).label("success_rate"),
            )
            .where(RequestLog.create_time >= start_time, RequestLog.model_id.is_not(None))
            .group_by(RequestLog.model_id)
        )
        rows = (await self.db.execute(stmt)).all()
        result: dict[int, _ModelStats] = {}
        for row in rows:
            model_id = int(row.model_id)
            avg_latency = int(row.avg_latency or 0)
            total_requests = int(row.total_requests or 0)
            success_rate = Decimal(str(row.success_rate if row.success_rate is not None else 100)).quantize(
                Decimal("0.01"), rounding=ROUND_HALF_UP
            )
            result[model_id] = _ModelStats(
                avg_latency=avg_latency,
                success_rate=success_rate,
                total_requests=total_requests,
            )
        return result

    def _calculate_normalization_params(self, models, stats_map: dict[int, _ModelStats]) -> _NormParams:
        params = _NormParams()
        costs = [model.input_price + model.output_price for model in models]
        if costs:
            params.min_cost = min(costs)
            params.max_cost = max(costs) if max(costs) != min(costs) else (min(costs) + Decimal("1"))
        latencies = [v.avg_latency for v in stats_map.values() if v.avg_latency > 0]
        if latencies:
            params.min_latency = min(latencies)
            params.max_latency = max(latencies)
        priorities = [model.priority for model in models]
        if priorities:
            params.min_priority = min(priorities)
            params.max_priority = max(priorities)
        return params

    def _calculate_score(self, model, stats: _ModelStats, params: _NormParams) -> Decimal:
        cost = model.input_price + model.output_price
        cost_score = self._normalize_decimal(cost, params.min_cost, params.max_cost)
        latency = stats.avg_latency if stats.avg_latency > 0 else 5000
        latency_score = self._normalize_int(latency, params.min_latency, params.max_latency)
        success_rate_score = 1.0 - float(stats.success_rate / Decimal("100"))
        priority_score = 1.0 - self._normalize_int(model.priority, params.min_priority, params.max_priority)
        score = (
            cost_score * SCORE_COST_WEIGHT
            + latency_score * SCORE_LATENCY_WEIGHT
            + success_rate_score * SCORE_SUCCESS_RATE_WEIGHT
            + priority_score * SCORE_PRIORITY_WEIGHT
        )
        return Decimal(str(score)).quantize(Decimal("0.0001"), rounding=ROUND_HALF_UP)

    async def _sync_provider_metrics_from_model_stats(self, models, stats_map: dict[int, _ModelStats]) -> None:
        grouped: dict[int, list] = {}
        for model in models:
            grouped.setdefault(model.provider_id, []).append(model)
        for provider_id, model_list in grouped.items():
            latencies = [stats_map[m.id].avg_latency for m in model_list if m.id in stats_map and stats_map[m.id].avg_latency > 0]
            success_rates = [
                stats_map[m.id].success_rate for m in model_list if m.id in stats_map and stats_map[m.id].total_requests > 0
            ]
            avg_latency = int(sum(latencies) / len(latencies)) if latencies else 0
            avg_success_rate = (
                (sum(success_rates, Decimal("0")) / Decimal(len(success_rates))).quantize(Decimal("0.01"), rounding=ROUND_HALF_UP)
                if success_rates
                else Decimal("100")
            )
            health_status = self._determine_health_status_by_success_rate(avg_success_rate)
            await self.model_provider_service.update_health_status(
                provider_id, health_status, avg_latency, avg_success_rate
            )

    @staticmethod
    async def _send_health_check_request(base_url: str, api_key: str) -> tuple[bool, int]:
        test_url = f"{base_url.rstrip('/')}/models"
        headers = {"Authorization": f"Bearer {api_key}"}
        start = datetime.utcnow()
        timeout_s = max(1, int(HEALTH_CHECK_TIMEOUT_MS / 1000))
        try:
            async with httpx.AsyncClient(timeout=timeout_s) as client:
                response = await client.get(test_url, headers=headers)
            latency = int((datetime.utcnow() - start).total_seconds() * 1000)
            ok = response.status_code in (200, 401)
            return ok, latency
        except Exception:
            latency = int((datetime.utcnow() - start).total_seconds() * 1000)
            return False, latency

    def _record_health_check_result(self, provider_id: int, is_healthy: bool) -> None:
        history = self.health_history_map.setdefault(provider_id, [])
        history.append(is_healthy)
        if len(history) > HEALTH_CHECK_MAX_HISTORY_SIZE:
            history.pop(0)

    def _calculate_success_rate(self, provider_id: int) -> Decimal:
        history = self.health_history_map.get(provider_id, [])
        if not history:
            return Decimal("100.00")
        success_count = sum(1 for item in history if item)
        rate = Decimal(success_count * 100) / Decimal(len(history))
        return rate.quantize(Decimal("0.01"), rounding=ROUND_HALF_UP)

    @staticmethod
    def _determine_health_status(success_rate: Decimal) -> str:
        value = float(success_rate)
        if value >= 80.0:
            return HEALTH_STATUS_HEALTHY
        if value >= 50.0:
            return HEALTH_STATUS_DEGRADED
        return HEALTH_STATUS_UNHEALTHY

    @staticmethod
    def _determine_health_status_by_success_rate(success_rate: Decimal) -> str:
        if success_rate is None:
            return HEALTH_STATUS_UNKNOWN
        return HealthCheckService._determine_health_status(success_rate)

    @staticmethod
    def _normalize_decimal(value: Decimal, min_value: Decimal, max_value: Decimal) -> float:
        diff = max_value - min_value
        if diff == 0:
            return 0.0
        return float((value - min_value) / diff)

    @staticmethod
    def _normalize_int(value: int, min_value: int, max_value: int) -> float:
        diff = max_value - min_value
        if diff == 0:
            return 0.0
        return float(value - min_value) / float(diff)
