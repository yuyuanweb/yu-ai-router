"""API key API."""

from __future__ import annotations

from fastapi import APIRouter, Depends, Request
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.core.constants import ErrorCode
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.middleware.auth import require_login
from app.models.user import User
from app.schemas.apikey import ApiKeyCreateRequest, ApiKeyVO
from app.schemas.common import BaseResponse, DeleteRequest, PageData
from app.services.api_key_service import ApiKeyService

router = APIRouter(prefix="/api/key", tags=["api-key"])


@router.post("/create", response_model=BaseResponse[ApiKeyVO])
async def create_api_key(
    payload: ApiKeyCreateRequest,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[ApiKeyVO]:
    entity = await ApiKeyService(db).create_api_key(payload.key_name, login_user.id)
    return success(ApiKeyService.to_api_key_vo(entity, mask=False))


@router.get("/list/my", response_model=BaseResponse[PageData[ApiKeyVO]])
async def list_my_api_keys(
    request: Request,
    page_num: int = 1,
    page_size: int = 10,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[PageData[ApiKeyVO]]:
    _ = request
    data = await ApiKeyService(db).list_user_api_key_page(login_user.id, page_num, page_size)
    return success(data)


@router.post("/revoke", response_model=BaseResponse[bool])
async def revoke_api_key(
    payload: DeleteRequest,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[bool]:
    if payload.id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "请求参数错误")
    ok = await ApiKeyService(db).revoke_api_key(payload.id, login_user.id)
    return success(ok)
