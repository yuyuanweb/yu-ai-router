"""Health API."""

from __future__ import annotations

from fastapi import APIRouter

from app.common.result_utils import success
from app.schemas.common import BaseResponse

router = APIRouter(prefix="/health", tags=["health"])


@router.get("/", response_model=BaseResponse[str])
async def health_check() -> BaseResponse[str]:
    return success("ok")
