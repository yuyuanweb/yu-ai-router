"""Response helpers."""

from __future__ import annotations

from typing import TypeVar

from app.core.constants import ERROR_MESSAGE_MAP, ErrorCode
from app.schemas.common import BaseResponse

T = TypeVar("T")


def success(data: T) -> BaseResponse[T]:
    return BaseResponse(code=ErrorCode.SUCCESS, data=data, message="ok")


def error(code: int, message: str | None = None) -> BaseResponse[None]:
    return BaseResponse(
        code=code,
        data=None,
        message=message or ERROR_MESSAGE_MAP.get(code, "系统内部异常"),
    )
