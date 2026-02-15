"""Plugin API."""

from __future__ import annotations

from fastapi import APIRouter, Depends, Query
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.core.constants import ErrorCode, UserRole
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.middleware.auth import require_login, require_role
from app.models.user import User
from app.schemas.common import BaseResponse
from app.schemas.plugin import PluginConfigVO, PluginExecuteRequest, PluginExecuteVO, PluginUpdateRequest
from app.services.plugin_service import PluginService

router = APIRouter(prefix="/plugin", tags=["plugin"])


@router.get("/list", response_model=BaseResponse[list[PluginConfigVO]])
async def list_plugins(
    _: object = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[list[PluginConfigVO]]:
    return success(await PluginService(db).list_all_plugins())


@router.get("/list/enabled", response_model=BaseResponse[list[PluginConfigVO]])
async def list_enabled_plugins(
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[list[PluginConfigVO]]:
    return success(await PluginService(db).list_enabled_plugins())


@router.get("/get", response_model=BaseResponse[PluginConfigVO])
async def get_plugin(
    plugin_key: str = Query(alias="pluginKey"),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[PluginConfigVO]:
    if not plugin_key:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "插件标识不能为空")
    plugin = await PluginService(db).get_plugin_by_key(plugin_key)
    if plugin is None:
        raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "插件不存在")
    return success(plugin)


@router.post("/update", response_model=BaseResponse[bool])
async def update_plugin(
    payload: PluginUpdateRequest,
    _: object = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[bool]:
    if payload is None:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "参数错误")
    return success(await PluginService(db).update_plugin(payload))


@router.post("/enable", response_model=BaseResponse[bool])
async def enable_plugin(
    plugin_key: str = Query(alias="pluginKey"),
    _: object = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[bool]:
    if not plugin_key:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "插件标识不能为空")
    return success(await PluginService(db).enable_plugin(plugin_key))


@router.post("/disable", response_model=BaseResponse[bool])
async def disable_plugin(
    plugin_key: str = Query(alias="pluginKey"),
    _: object = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[bool]:
    if not plugin_key:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "插件标识不能为空")
    return success(await PluginService(db).disable_plugin(plugin_key))


@router.post("/execute", response_model=BaseResponse[PluginExecuteVO])
async def execute_plugin(
    payload: PluginExecuteRequest,
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[PluginExecuteVO]:
    if payload is None:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "参数错误")
    data = await PluginService(db).execute_plugin(payload, login_user.id if login_user else None)
    return success(data)


@router.post("/reload", response_model=BaseResponse[bool])
async def reload_plugin(
    plugin_key: str = Query(alias="pluginKey"),
    _: object = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[bool]:
    if not plugin_key:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "插件标识不能为空")
    await PluginService(db).reload_plugin(plugin_key)
    return success(True)


@router.post("/reload/all", response_model=BaseResponse[bool])
async def reload_all_plugins(
    _: object = Depends(require_role(UserRole.ADMIN)),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[bool]:
    await PluginService(db).init_plugins()
    return success(True)
