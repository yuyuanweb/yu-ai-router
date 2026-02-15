"""FastAPI entrypoint."""

from __future__ import annotations

import logging
import traceback
import uuid
import asyncio

from fastapi import FastAPI, Request
from fastapi.exceptions import RequestValidationError
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse
from starlette.middleware.base import BaseHTTPMiddleware

from app.api.health import router as health_router
from app.api.apikey import router as apikey_router
from app.api.blacklist import router as blacklist_router
from app.api.chat import router as chat_router
from app.api.image import router as image_router
from app.api.internal_chat import router as internal_chat_router
from app.api.model import router as model_router
from app.api.model_provider import router as model_provider_router
from app.api.balance import router as balance_router
from app.api.recharge import router as recharge_router
from app.api.stripe_webhook import router as stripe_webhook_router
from app.api.stats import router as stats_router
from app.api.user import router as user_router
from app.common.result_utils import error
from app.core.config import get_settings
from app.core.constants import ErrorCode
from app.core.logging_config import setup_logging
from app.exceptions.business_exception import BusinessException
from app.middleware.ip_blacklist import IpBlacklistMiddleware
from app.task.health_check_task import HealthCheckTask

settings = get_settings()
setup_logging(settings.log_level)
logger = logging.getLogger("app")


class TraceIdMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        request_id = request.headers.get("X-Trace-Id") or str(uuid.uuid4())
        request.state.trace_id = request_id
        request.state.login_user_id = None
        response = await call_next(request)
        response.headers["X-Trace-Id"] = request_id
        return response


def create_app() -> FastAPI:
    app = FastAPI(title=settings.app_name)
    app.add_middleware(IpBlacklistMiddleware)
    app.add_middleware(TraceIdMiddleware)
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_origin_regex=".*",
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    @app.exception_handler(BusinessException)
    async def business_exception_handler(request: Request, exc: BusinessException):
        logger.error(
            "BusinessException path=%s userId=%s code=%s message=%s traceId=%s\n%s",
            request.url.path,
            getattr(request.state, "login_user_id", None),
            exc.code,
            exc.message,
            getattr(request.state, "trace_id", None),
            traceback.format_exc(),
        )
        return JSONResponse(
            status_code=200,
            content=error(exc.code, exc.message).model_dump(by_alias=True),
        )

    @app.exception_handler(RequestValidationError)
    async def validation_exception_handler(request: Request, exc: RequestValidationError):
        logger.error(
            "ValidationException path=%s userId=%s code=%s message=%s traceId=%s details=%s",
            request.url.path,
            getattr(request.state, "login_user_id", None),
            ErrorCode.PARAMS_ERROR,
            "请求参数错误",
            getattr(request.state, "trace_id", None),
            exc.errors(),
        )
        return JSONResponse(
            status_code=200,
            content=error(ErrorCode.PARAMS_ERROR, "请求参数错误").model_dump(by_alias=True),
        )

    @app.exception_handler(Exception)
    async def runtime_exception_handler(request: Request, exc: Exception):
        logger.error(
            "RuntimeException path=%s userId=%s code=%s message=%s traceId=%s\n%s",
            request.url.path,
            getattr(request.state, "login_user_id", None),
            ErrorCode.SYSTEM_ERROR,
            str(exc),
            getattr(request.state, "trace_id", None),
            traceback.format_exc(),
        )
        return JSONResponse(
            status_code=200,
            content=error(ErrorCode.SYSTEM_ERROR, "系统错误").model_dump(by_alias=True),
        )

    app.include_router(health_router, prefix=settings.app_base_path)
    app.include_router(user_router, prefix=settings.app_base_path)
    app.include_router(apikey_router, prefix=settings.app_base_path)
    app.include_router(chat_router, prefix=settings.app_base_path)
    app.include_router(image_router, prefix=settings.app_base_path)
    app.include_router(internal_chat_router, prefix=settings.app_base_path)
    app.include_router(stats_router, prefix=settings.app_base_path)
    app.include_router(model_router, prefix=settings.app_base_path)
    app.include_router(model_provider_router, prefix=settings.app_base_path)
    app.include_router(blacklist_router, prefix=settings.app_base_path)
    app.include_router(balance_router, prefix=settings.app_base_path)
    app.include_router(recharge_router, prefix=settings.app_base_path)
    app.include_router(stripe_webhook_router, prefix=settings.app_base_path)

    @app.on_event("startup")
    async def startup_health_check_task() -> None:
        stop_event = asyncio.Event()
        app.state.health_check_stop_event = stop_event
        app.state.health_check_runner = asyncio.create_task(HealthCheckTask().run_loop(stop_event))

    @app.on_event("shutdown")
    async def shutdown_health_check_task() -> None:
        stop_event = getattr(app.state, "health_check_stop_event", None)
        runner = getattr(app.state, "health_check_runner", None)
        if stop_event is not None:
            stop_event.set()
        if runner is not None:
            await runner
    return app


app = create_app()
