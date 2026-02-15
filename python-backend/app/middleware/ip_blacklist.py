"""IP blacklist middleware."""

from __future__ import annotations

from fastapi import Request
from fastapi.responses import JSONResponse
from starlette.middleware.base import BaseHTTPMiddleware

from app.common.result_utils import error
from app.core.constants import ErrorCode
from app.infra.redis_client import redis_client
from app.services.blacklist_service import BlacklistService
from app.utils.request import get_client_ip


class IpBlacklistMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        client_ip = get_client_ip(request)
        blocked = await BlacklistService(redis_client).is_blocked(client_ip)
        if blocked:
            payload = error(ErrorCode.FORBIDDEN_ERROR, "您的 IP 已被封禁")
            return JSONResponse(
                status_code=403,
                content=payload.model_dump(by_alias=True),
            )
        return await call_next(request)
