"""Image generation service."""

from __future__ import annotations

import asyncio
import logging
import time
from decimal import Decimal
from math import ceil

import httpx
from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import (
    DEFAULT_IMAGE_COUNT,
    DEFAULT_IMAGE_MODEL,
    DEFAULT_IMAGE_SIZE,
    DEFAULT_PAGE_NUM,
    DEFAULT_PAGE_SIZE,
    ErrorCode,
    IMAGE_RECORD_STATUS_FAILED,
    IMAGE_RECORD_STATUS_SUCCESS,
    IMAGE_TOKENS_PER_IMAGE,
    MAX_PAGE_SIZE,
    MODEL_TYPE_IMAGE,
)
from app.exceptions.business_exception import BusinessException
from app.models.image_generation_record import ImageGenerationRecord
from app.schemas.common import PageData
from app.schemas.image import (
    ImageData,
    ImageGenerationRecordVO,
    ImageGenerationRequest,
    ImageGenerationResponse,
)
from app.services.balance_service import BalanceService
from app.services.model_provider_service import ModelProviderService
from app.services.model_service import ModelService
from app.services.quota_service import QuotaService
from app.services.user_service import UserService

logger = logging.getLogger("app")


class ImageGenerationService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db
        self.model_service = ModelService(db)
        self.model_provider_service = ModelProviderService(db)
        self.quota_service = QuotaService(db)
        self.balance_service = BalanceService(db)
        self.user_service = UserService(db)

    async def generate_image(
        self,
        request: ImageGenerationRequest,
        user_id: int,
        api_key_id: int | None,
        client_ip: str | None,
    ) -> ImageGenerationResponse:
        start_ms = int(time.time() * 1000)
        if not request.prompt or not request.prompt.strip():
            raise BusinessException(ErrorCode.PARAMS_ERROR, "提示词不能为空")

        if await self.user_service.is_user_disabled(user_id):
            raise BusinessException(ErrorCode.FORBIDDEN_ERROR, "账号已被禁用，无法使用服务")
        if not await self.quota_service.check_quota(user_id):
            raise BusinessException(ErrorCode.OPERATION_ERROR, "Token配额已用尽，请联系管理员增加配额")

        model_key = request.model or DEFAULT_IMAGE_MODEL
        size = request.size or DEFAULT_IMAGE_SIZE
        image_count = DEFAULT_IMAGE_COUNT

        model = await self.model_service.get_by_model_key(model_key)
        if model is None or model.model_type != MODEL_TYPE_IMAGE:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "模型不存在或不是绘图模型")
        provider = await self.model_provider_service.get_by_id(model.provider_id)
        if provider is None:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "模型提供者不存在")

        estimated_cost = Decimal(model.input_price or 0) * Decimal(image_count)
        if not await self.balance_service.check_balance(user_id, estimated_cost):
            raise BusinessException(
                ErrorCode.OPERATION_ERROR,
                f"账户余额不足，生成{image_count}张图片预计需要¥{estimated_cost}，请先充值",
            )

        try:
            response = await self._call_dashscope_image_api(
                provider_base_url=provider.base_url,
                provider_api_key=provider.api_key,
                model_key=model.model_key,
                prompt=request.prompt,
                size=size,
                n=image_count,
            )
            duration = int(time.time() * 1000) - start_ms
            actual_count = len(response.data) if response.data else image_count
            actual_cost = Decimal(model.input_price or 0) * Decimal(actual_count)

            for item in response.data:
                self.db.add(
                    ImageGenerationRecord(
                        user_id=user_id,
                        api_key_id=api_key_id,
                        model_id=model.id,
                        model_key=model_key,
                        prompt=request.prompt,
                        revised_prompt=item.revised_prompt,
                        image_url=item.url,
                        image_data=item.b64_json,
                        size=size,
                        quality=request.quality,
                        status=IMAGE_RECORD_STATUS_SUCCESS,
                        cost=Decimal(model.input_price or 0),
                        duration=duration,
                        client_ip=client_ip,
                    )
                )
            await self.db.commit()

            await self.quota_service.deduct_tokens(user_id, actual_count * IMAGE_TOKENS_PER_IMAGE)
            if actual_cost > 0:
                source = "API图片生成" if api_key_id is not None else "网页图片生成"
                await self.balance_service.deduct_balance(
                    user_id,
                    actual_cost,
                    request_log_id=None,
                    description=f"{source} - {model_key} x{actual_count}",
                )

            logger.info("图片生成成功：用户 %s，模型 %s，数量 %s", user_id, model_key, actual_count)
            return response
        except BusinessException as exc:
            await self._save_failed_record(
                user_id=user_id,
                api_key_id=api_key_id,
                model_id=model.id,
                model_key=model_key,
                request=request,
                size=size,
                client_ip=client_ip,
                start_ms=start_ms,
                error_message=str(exc),
            )
            raise
        except Exception as exc:
            await self._save_failed_record(
                user_id=user_id,
                api_key_id=api_key_id,
                model_id=model.id,
                model_key=model_key,
                request=request,
                size=size,
                client_ip=client_ip,
                start_ms=start_ms,
                error_message=str(exc),
            )
            raise BusinessException(ErrorCode.SYSTEM_ERROR, f"图片生成失败: {exc}") from exc

    async def list_user_records(self, user_id: int, page_num: int, page_size: int) -> PageData[ImageGenerationRecordVO]:
        safe_page_num = page_num if page_num > 0 else DEFAULT_PAGE_NUM
        safe_page_size = page_size if page_size > 0 else DEFAULT_PAGE_SIZE
        safe_page_size = min(safe_page_size, MAX_PAGE_SIZE)

        base_stmt = select(ImageGenerationRecord).where(ImageGenerationRecord.user_id == user_id)
        total_row = await self.db.scalar(select(func.count()).select_from(base_stmt.subquery())) or 0
        rows = (
            await self.db.scalars(
                base_stmt.order_by(ImageGenerationRecord.create_time.desc())
                .offset((safe_page_num - 1) * safe_page_size)
                .limit(safe_page_size)
            )
        ).all()
        return PageData[ImageGenerationRecordVO](
            records=[ImageGenerationRecordVO.model_validate(item) for item in rows],
            pageNumber=safe_page_num,
            pageSize=safe_page_size,
            totalPage=ceil(total_row / safe_page_size) if total_row else 0,
            totalRow=total_row,
            optimizeCountQuery=True,
        )

    async def _save_failed_record(
        self,
        user_id: int,
        api_key_id: int | None,
        model_id: int,
        model_key: str,
        request: ImageGenerationRequest,
        size: str,
        client_ip: str | None,
        start_ms: int,
        error_message: str,
    ) -> None:
        duration = int(time.time() * 1000) - start_ms
        self.db.add(
            ImageGenerationRecord(
                user_id=user_id,
                api_key_id=api_key_id,
                model_id=model_id,
                model_key=model_key,
                prompt=request.prompt,
                size=size,
                quality=request.quality,
                status=IMAGE_RECORD_STATUS_FAILED,
                cost=Decimal("0"),
                duration=duration,
                error_message=error_message[:512],
                client_ip=client_ip,
            )
        )
        await self.db.commit()

    async def _call_dashscope_image_api(
        self,
        provider_base_url: str,
        provider_api_key: str,
        model_key: str,
        prompt: str,
        size: str,
        n: int,
    ) -> ImageGenerationResponse:
        normalized_base_url = provider_base_url.rstrip("/").replace("/compatible-mode", "")
        create_task_url = f"{normalized_base_url}/api/v1/services/aigc/text2image/image-synthesis"
        headers = {
            "Authorization": f"Bearer {provider_api_key}",
            "Content-Type": "application/json",
            "X-DashScope-Async": "enable",
        }
        payload = {
            "model": model_key,
            "input": {"prompt": prompt},
            "parameters": {"size": size, "n": n},
        }

        async with httpx.AsyncClient(timeout=30) as client:
            create_resp = await client.post(create_task_url, headers=headers, json=payload)
            create_data = create_resp.json()

            task_id = (((create_data.get("output") or {}).get("task_id")) or "").strip()
            if not task_id:
                message = create_data.get("message") or "未返回任务ID"
                raise BusinessException(ErrorCode.SYSTEM_ERROR, f"创建图片生成任务失败：{message}")

            task_url = f"{normalized_base_url}/api/v1/tasks/{task_id}"
            max_retries = 60
            for _ in range(max_retries):
                await asyncio.sleep(2)
                task_resp = await client.get(task_url, headers={"Authorization": f"Bearer {provider_api_key}"}, timeout=10)
                task_data = task_resp.json()
                output = task_data.get("output") or {}
                task_status = output.get("task_status")
                if task_status == "SUCCEEDED":
                    results = output.get("results") or []
                    image_data = [
                        ImageData(url=item.get("url"), revisedPrompt=prompt)
                        for item in results
                        if item.get("url")
                    ]
                    if not image_data:
                        raise BusinessException(ErrorCode.SYSTEM_ERROR, "任务成功但未返回图片")
                    return ImageGenerationResponse(created=int(time.time()), data=image_data)
                if task_status == "FAILED":
                    message = output.get("message") or "任务失败"
                    raise BusinessException(ErrorCode.SYSTEM_ERROR, f"图片生成失败：{message}")

            raise BusinessException(ErrorCode.SYSTEM_ERROR, "图片生成超时")
