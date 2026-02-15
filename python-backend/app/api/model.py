"""Model API."""

from __future__ import annotations

from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.core.constants import ErrorCode, MODEL_STATUS_ACTIVE, UserRole
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.middleware.auth import require_role
from app.models.model import Model
from app.schemas.common import BaseResponse, DeleteRequest, PageData
from app.schemas.model import ModelAddRequest, ModelQueryRequest, ModelUpdateRequest, ModelVO
from app.services.model_service import ModelService

router = APIRouter(prefix="/model", tags=["model"])


@router.post("/add", response_model=BaseResponse[int])
async def add_model(
    payload: ModelAddRequest,
    _: object = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
):
    if payload is None:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "参数错误")
    model = Model(
        provider_id=payload.provider_id,
        model_key=payload.model_key,
        model_name=payload.model_name,
        model_type=payload.model_type,
        description=payload.description,
        context_length=payload.context_length or 4096,
        input_price=payload.input_price or 0,
        output_price=payload.output_price or 0,
        priority=payload.priority or 100,
        default_timeout=payload.default_timeout or 60000,
        capabilities=payload.capabilities,
        status=MODEL_STATUS_ACTIVE,
    )
    model_id = await ModelService(db).add_model(model)
    return success(model_id)


@router.post("/delete", response_model=BaseResponse[bool])
async def delete_model(
    payload: DeleteRequest,
    _: object = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
):
    if payload.id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "参数错误")
    return success(await ModelService(db).delete_model(payload.id))


@router.post("/update", response_model=BaseResponse[bool])
async def update_model(
    payload: ModelUpdateRequest,
    _: object = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
):
    if payload.id is None or payload.id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "参数错误")
    model = Model(
        id=payload.id,
        model_name=payload.model_name,
        description=payload.description,
        context_length=payload.context_length,
        input_price=payload.input_price,
        output_price=payload.output_price,
        status=payload.status,
        priority=payload.priority,
        default_timeout=payload.default_timeout,
        capabilities=payload.capabilities,
    )
    return success(await ModelService(db).update_model(model))


@router.get("/get/vo", response_model=BaseResponse[ModelVO])
async def get_model_vo(id: int, db: AsyncSession = Depends(get_db_session)):
    if id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "参数错误")
    model_vo = await ModelService(db).get_model_vo(id)
    if model_vo is None:
        raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "数据不存在")
    return success(model_vo)


@router.post("/list/page/vo", response_model=BaseResponse[PageData[ModelVO]])
async def list_model_page_vo(payload: ModelQueryRequest, db: AsyncSession = Depends(get_db_session)):
    return success(await ModelService(db).list_page_vo(payload))


@router.get("/list/vo", response_model=BaseResponse[list[ModelVO]])
async def list_model_vo(db: AsyncSession = Depends(get_db_session)):
    return success(await ModelService(db).list_vo())


@router.get("/list/active", response_model=BaseResponse[list[ModelVO]])
async def list_active_model_vo(db: AsyncSession = Depends(get_db_session)):
    return success(await ModelService(db).list_active())


@router.get("/list/active/provider/{provider_id}", response_model=BaseResponse[list[ModelVO]])
async def list_active_model_by_provider(provider_id: int, db: AsyncSession = Depends(get_db_session)):
    if provider_id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "参数错误")
    return success(await ModelService(db).list_active_by_provider(provider_id))


@router.get("/list/active/type/{model_type}", response_model=BaseResponse[list[ModelVO]])
async def list_active_model_by_type(model_type: str, db: AsyncSession = Depends(get_db_session)):
    if not model_type:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "参数错误")
    return success(await ModelService(db).list_active_by_type(model_type))
