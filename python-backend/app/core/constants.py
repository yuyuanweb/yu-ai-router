"""Project constants."""

from __future__ import annotations

from enum import StrEnum


class UserRole(StrEnum):
    USER = "user"
    ADMIN = "admin"


class ErrorCode:
    SUCCESS = 0
    PARAMS_ERROR = 40000
    NOT_LOGIN_ERROR = 40100
    NO_AUTH_ERROR = 40101
    FORBIDDEN_ERROR = 40300
    NOT_FOUND_ERROR = 40400
    TOO_MANY_REQUEST = 42900
    SYSTEM_ERROR = 50000
    OPERATION_ERROR = 50001


ERROR_MESSAGE_MAP: dict[int, str] = {
    ErrorCode.SUCCESS: "ok",
    ErrorCode.PARAMS_ERROR: "请求参数错误",
    ErrorCode.NOT_LOGIN_ERROR: "未登录",
    ErrorCode.NO_AUTH_ERROR: "无权限",
    ErrorCode.FORBIDDEN_ERROR: "禁止访问",
    ErrorCode.NOT_FOUND_ERROR: "请求数据不存在",
    ErrorCode.TOO_MANY_REQUEST: "请求过于频繁",
    ErrorCode.SYSTEM_ERROR: "系统内部异常",
    ErrorCode.OPERATION_ERROR: "操作失败",
}

USER_LOGIN_STATE = "user_login"
SESSION_COOKIE_NAME = "SESSION"
SESSION_KEY_PREFIX = "session:"
SESSION_EXPIRE_SECONDS = 30 * 24 * 60 * 60
PASSWORD_SALT = "yupi"
DEFAULT_USER_NAME = "无名"
DEFAULT_USER_PASSWORD = "12345678"
DEFAULT_PAGE_NUM = 1
DEFAULT_PAGE_SIZE = 10
MAX_PAGE_SIZE = 100
