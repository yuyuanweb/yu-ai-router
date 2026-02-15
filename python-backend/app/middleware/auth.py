"""Auth helpers based on Redis session."""

from __future__ import annotations

import json
import uuid

from fastapi import Depends, Request, Response
from redis.asyncio import Redis
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import (
    ErrorCode,
    SESSION_COOKIE_NAME,
    SESSION_EXPIRE_SECONDS,
    SESSION_KEY_PREFIX,
    UserRole,
)
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.infra.redis_client import get_redis_client
from app.models.user import User
from app.services.user_service import UserService


def _session_key(session_id: str) -> str:
    return f"{SESSION_KEY_PREFIX}{session_id}"


async def save_login_session(response: Response, redis: Redis, user: User) -> None:
    session_id = str(uuid.uuid4())
    payload = {"userId": user.id}
    await redis.set(_session_key(session_id), json.dumps(payload), ex=SESSION_EXPIRE_SECONDS)
    response.set_cookie(
        key=SESSION_COOKIE_NAME,
        value=session_id,
        max_age=SESSION_EXPIRE_SECONDS,
        httponly=True,
        samesite="lax",
        path="/",
    )


async def clear_login_session(request: Request, response: Response, redis: Redis) -> None:
    session_id = request.cookies.get(SESSION_COOKIE_NAME)
    if not session_id:
        raise BusinessException(ErrorCode.OPERATION_ERROR, "用户未登录")
    await redis.delete(_session_key(session_id))
    response.delete_cookie(SESSION_COOKIE_NAME, path="/")


async def get_login_user(
    request: Request,
    db: AsyncSession,
    redis: Redis,
) -> User:
    session_id = request.cookies.get(SESSION_COOKIE_NAME)
    if not session_id:
        raise BusinessException(ErrorCode.NOT_LOGIN_ERROR, "未登录")
    payload = await redis.get(_session_key(session_id))
    if not payload:
        raise BusinessException(ErrorCode.NOT_LOGIN_ERROR, "未登录")
    session = json.loads(payload)
    user_id = session.get("userId")
    if not user_id:
        raise BusinessException(ErrorCode.NOT_LOGIN_ERROR, "未登录")
    user = await UserService(db).get_by_id(int(user_id))
    if user is None:
        raise BusinessException(ErrorCode.NOT_LOGIN_ERROR, "未登录")
    await redis.expire(_session_key(session_id), SESSION_EXPIRE_SECONDS)
    request.state.login_user_id = user.id
    return user


async def require_login(
    request: Request,
    db: AsyncSession = Depends(get_db_session),
    redis: Redis = Depends(get_redis_client),
) -> User:
    return await get_login_user(request, db, redis)


def require_role(must_role: UserRole):
    async def checker(
        request: Request,
        db: AsyncSession = Depends(get_db_session),
        redis: Redis = Depends(get_redis_client),
    ) -> User:
        login_user = await get_login_user(request, db, redis)
        if must_role == UserRole.ADMIN and login_user.user_role != UserRole.ADMIN.value:
            raise BusinessException(ErrorCode.NO_AUTH_ERROR, "无权限")
        return login_user

    return checker
