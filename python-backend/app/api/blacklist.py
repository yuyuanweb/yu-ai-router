"""Blacklist management API."""

from __future__ import annotations

from fastapi import APIRouter, Depends, Query
from redis.asyncio import Redis

from app.common.result_utils import success
from app.core.constants import UserRole
from app.infra.redis_client import get_redis_client
from app.middleware.auth import require_role
from app.schemas.blacklist import BlacklistRequest
from app.schemas.common import BaseResponse
from app.services.blacklist_service import BlacklistService

router = APIRouter(prefix="/admin/blacklist", tags=["blacklist"])


@router.get("/list", response_model=BaseResponse[set[str]])
async def get_blacklist(
    _: object = Depends(require_role(UserRole.ADMIN)),
    redis: Redis = Depends(get_redis_client),
) -> BaseResponse[set[str]]:
    data = await BlacklistService(redis).get_all_blacklist()
    return success(data)


@router.post("/add", response_model=BaseResponse[bool])
async def add_to_blacklist(
    payload: BlacklistRequest,
    _: object = Depends(require_role(UserRole.ADMIN)),
    redis: Redis = Depends(get_redis_client),
) -> BaseResponse[bool]:
    await BlacklistService(redis).add_to_blacklist(payload.ip, payload.reason)
    return success(True)


@router.post("/remove", response_model=BaseResponse[bool])
async def remove_from_blacklist(
    payload: BlacklistRequest,
    _: object = Depends(require_role(UserRole.ADMIN)),
    redis: Redis = Depends(get_redis_client),
) -> BaseResponse[bool]:
    await BlacklistService(redis).remove_from_blacklist(payload.ip)
    return success(True)


@router.get("/check", response_model=BaseResponse[bool])
async def check_blacklist(
    ip: str = Query(...),
    _: object = Depends(require_role(UserRole.ADMIN)),
    redis: Redis = Depends(get_redis_client),
) -> BaseResponse[bool]:
    blocked = await BlacklistService(redis).is_blocked(ip)
    return success(blocked)


@router.get("/count", response_model=BaseResponse[int])
async def get_blacklist_count(
    _: object = Depends(require_role(UserRole.ADMIN)),
    redis: Redis = Depends(get_redis_client),
) -> BaseResponse[int]:
    count = await BlacklistService(redis).get_blacklist_count()
    return success(count)
