"""BYOK APIs."""

from __future__ import annotations

from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.core.constants import ErrorCode
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.middleware.auth import require_login
from app.models.user import User
from app.schemas.byok import UserProviderKeyAddRequest, UserProviderKeyUpdateRequest, UserProviderKeyVO
from app.schemas.common import BaseResponse, DeleteRequest
from app.services.user_provider_key_service import UserProviderKeyService

router = APIRouter(prefix="/byok", tags=["byok"])


@router.post("/add", response_model=BaseResponse[bool])
async def add_user_provider_key(
    payload: UserProviderKeyAddRequest,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[bool]:
    ok = await UserProviderKeyService(db).add_user_provider_key(payload, login_user.id)
    return success(ok)


@router.post("/update", response_model=BaseResponse[bool])
async def update_user_provider_key(
    payload: UserProviderKeyUpdateRequest,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[bool]:
    ok = await UserProviderKeyService(db).update_user_provider_key(payload, login_user.id)
    return success(ok)


@router.post("/delete", response_model=BaseResponse[bool])
async def delete_user_provider_key(
    payload: DeleteRequest,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[bool]:
    if payload.id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "请求参数错误")
    ok = await UserProviderKeyService(db).delete_user_provider_key(payload.id, login_user.id)
    return success(ok)


@router.get("/my/list", response_model=BaseResponse[list[UserProviderKeyVO]])
async def list_my_provider_keys(
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[list[UserProviderKeyVO]]:
    rows = await UserProviderKeyService(db).list_user_provider_keys(login_user.id)
    return success(rows)
