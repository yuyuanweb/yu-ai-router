"""Model provider API."""

from __future__ import annotations

from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.core.constants import ErrorCode, PROVIDER_STATUS_ACTIVE, UserRole
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.middleware.auth import require_role
from app.models.model_provider import ModelProvider
from app.schemas.common import BaseResponse, DeleteRequest, PageData
from app.schemas.model_provider import (
    ProviderAddRequest,
    ProviderQueryRequest,
    ProviderUpdateRequest,
    ProviderVO,
)
from app.services.model_provider_service import ModelProviderService

router = APIRouter(prefix="/provider", tags=["provider"])


@router.post("/add", response_model=BaseResponse[int])
async def add_provider(
    payload: ProviderAddRequest,
    _: object = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
):
    if payload is None:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "参数错误")
    provider = ModelProvider(
        provider_name=payload.provider_name,
        display_name=payload.display_name,
        base_url=payload.base_url,
        api_key=payload.api_key,
        priority=payload.priority or 100,
        config=payload.config,
        status=PROVIDER_STATUS_ACTIVE,
    )
    provider_id = await ModelProviderService(db).add_provider(provider)
    return success(provider_id)


@router.post("/delete", response_model=BaseResponse[bool])
async def delete_provider(
    payload: DeleteRequest,
    _: object = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
):
    if payload.id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "参数错误")
    return success(await ModelProviderService(db).delete_provider(payload.id))


@router.post("/update", response_model=BaseResponse[bool])
async def update_provider(
    payload: ProviderUpdateRequest,
    _: object = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
):
    if payload.id is None or payload.id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "参数错误")
    provider = ModelProvider(
        id=payload.id,
        display_name=payload.display_name,
        base_url=payload.base_url,
        api_key=payload.api_key,
        status=payload.status,
        priority=payload.priority,
        config=payload.config,
    )
    return success(await ModelProviderService(db).update_provider(provider))


@router.get("/get/vo", response_model=BaseResponse[ProviderVO])
async def get_provider_vo(id: int, db: AsyncSession = Depends(get_db_session)):
    if id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "参数错误")
    entity = await ModelProviderService(db).get_by_id(id)
    if entity is None:
        raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "数据不存在")
    return success(ProviderVO.model_validate(entity))


@router.post("/list/page/vo", response_model=BaseResponse[PageData[ProviderVO]])
async def list_provider_page_vo(payload: ProviderQueryRequest, db: AsyncSession = Depends(get_db_session)):
    page_data = await ModelProviderService(db).list_page_vo(payload)
    return success(page_data)


@router.get("/list/vo", response_model=BaseResponse[list[ProviderVO]])
async def list_provider_vo(db: AsyncSession = Depends(get_db_session)):
    return success(await ModelProviderService(db).list_vo())


@router.get("/list/healthy", response_model=BaseResponse[list[ProviderVO]])
async def list_healthy_provider_vo(db: AsyncSession = Depends(get_db_session)):
    return success(await ModelProviderService(db).list_healthy())
