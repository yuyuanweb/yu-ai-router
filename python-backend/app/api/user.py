"""User API."""

from __future__ import annotations

from fastapi import APIRouter, Depends, Request, Response
from redis.asyncio import Redis
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.core.constants import DEFAULT_USER_PASSWORD, ErrorCode, UserRole
from app.core.security import encrypt_password
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.infra.redis_client import get_redis_client
from app.middleware.auth import clear_login_session, require_login, require_role, save_login_session
from app.models.user import User
from app.schemas.common import BaseResponse, DeleteRequest, PageData
from app.schemas.user import (
    LoginUserVO,
    QuotaUpdateRequest,
    QuotaVO,
    UserAddRequest,
    UserAnalysisVO,
    UserLoginRequest,
    UserQueryRequest,
    UserRawVO,
    UserRegisterRequest,
    UserUpdateRequest,
    UserVO,
)
from app.services.billing_service import BillingService
from app.services.quota_service import QuotaService
from app.services.request_log_service import RequestLogService
from app.services.user_service import UserService

router = APIRouter(prefix="/user", tags=["user"])


@router.post("/register", response_model=BaseResponse[str])
async def user_register(
    payload: UserRegisterRequest,
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[str]:
    user_id = await UserService(db).register(
        payload.user_account,
        payload.user_password,
        payload.check_password,
    )
    return success(str(user_id))


@router.post("/login", response_model=BaseResponse[LoginUserVO])
async def user_login(
    payload: UserLoginRequest,
    response: Response,
    db: AsyncSession = Depends(get_db_session),
    redis: Redis = Depends(get_redis_client),
) -> BaseResponse[LoginUserVO]:
    user = await UserService(db).login(payload.user_account, payload.user_password)
    await save_login_session(response, redis, user)
    return success(UserService.to_login_user_vo(user))


@router.get("/get/login", response_model=BaseResponse[LoginUserVO])
async def get_login(
    login_user: User = Depends(require_login),
) -> BaseResponse[LoginUserVO]:
    return success(UserService.to_login_user_vo(login_user))


@router.post("/logout", response_model=BaseResponse[bool])
async def user_logout(
    request: Request,
    response: Response,
    redis: Redis = Depends(get_redis_client),
) -> BaseResponse[bool]:
    await clear_login_session(request, response, redis)
    return success(True)


@router.post("/add", response_model=BaseResponse[str])
async def add_user(
    payload: UserAddRequest,
    db: AsyncSession = Depends(get_db_session),
    _: User = Depends(require_role(UserRole.ADMIN)),
) -> BaseResponse[str]:
    user_id = await UserService(db).add_user(
        user_account=payload.user_account,
        user_password=encrypt_password(DEFAULT_USER_PASSWORD),
        user_name=payload.user_name,
        user_avatar=payload.user_avatar,
        user_profile=payload.user_profile,
        user_role=payload.user_role,
    )
    return success(str(user_id))


@router.get("/get", response_model=BaseResponse[UserRawVO])
async def get_user_by_id(
    id: int,
    db: AsyncSession = Depends(get_db_session),
    _: User = Depends(require_role(UserRole.ADMIN)),
) -> BaseResponse[UserRawVO]:
    if id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "请求参数错误")
    user = await UserService(db).get_by_id(id)
    if user is None:
        raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "请求数据不存在")
    return success(UserRawVO.model_validate(user))


@router.get("/get/vo", response_model=BaseResponse[UserVO])
async def get_user_vo_by_id(
    id: int,
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[UserVO]:
    if id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "请求参数错误")
    user = await UserService(db).get_by_id(id)
    if user is None:
        raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "请求数据不存在")
    return success(UserService.to_user_vo(user))


@router.post("/delete", response_model=BaseResponse[bool])
async def delete_user(
    payload: DeleteRequest,
    db: AsyncSession = Depends(get_db_session),
    _: User = Depends(require_role(UserRole.ADMIN)),
) -> BaseResponse[bool]:
    if payload.id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "请求参数错误")
    ok = await UserService(db).delete_user(payload.id)
    return success(ok)


@router.post("/update", response_model=BaseResponse[bool])
async def update_user(
    payload: UserUpdateRequest,
    db: AsyncSession = Depends(get_db_session),
    _: User = Depends(require_role(UserRole.ADMIN)),
) -> BaseResponse[bool]:
    ok = await UserService(db).update_user(
        user_id=payload.id,
        user_name=payload.user_name,
        user_avatar=payload.user_avatar,
        user_profile=payload.user_profile,
        user_role=payload.user_role,
    )
    return success(ok)


@router.post("/list/page/vo", response_model=BaseResponse[PageData[UserVO]])
async def list_user_vo_by_page(
    payload: UserQueryRequest,
    db: AsyncSession = Depends(get_db_session),
    _: User = Depends(require_role(UserRole.ADMIN)),
) -> BaseResponse[PageData[UserVO]]:
    data = await UserService(db).list_user_vo_page(
        page_num=payload.page_num,
        page_size=payload.page_size,
        user_id=payload.id,
        user_account=payload.user_account,
        user_name=payload.user_name,
        user_profile=payload.user_profile,
        user_role=payload.user_role,
        sort_field=payload.sort_field,
        sort_order=payload.sort_order,
    )
    return success(data)


@router.get("/quota/my", response_model=BaseResponse[QuotaVO])
async def get_my_quota(
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[QuotaVO]:
    remaining_quota = await QuotaService(db).get_remaining_quota(login_user.id)
    return success(
        QuotaVO(
            tokenQuota=login_user.token_quota,
            usedTokens=login_user.used_tokens or 0,
            remainingQuota=remaining_quota,
        )
    )


@router.post("/quota/set", response_model=BaseResponse[bool])
async def set_user_quota(
    payload: QuotaUpdateRequest,
    db: AsyncSession = Depends(get_db_session),
    _: User = Depends(require_role(UserRole.ADMIN)),
) -> BaseResponse[bool]:
    await UserService(db).set_user_quota(payload.user_id, payload.token_quota)
    return success(True)


@router.post("/quota/reset", response_model=BaseResponse[bool])
async def reset_user_quota(
    user_id: int,
    db: AsyncSession = Depends(get_db_session),
    _: User = Depends(require_role(UserRole.ADMIN)),
) -> BaseResponse[bool]:
    if user_id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "用户ID不合法")
    await UserService(db).reset_user_used_tokens(user_id)
    return success(True)


@router.post("/disable", response_model=BaseResponse[bool])
async def disable_user(
    user_id: int,
    db: AsyncSession = Depends(get_db_session),
    _: User = Depends(require_role(UserRole.ADMIN)),
) -> BaseResponse[bool]:
    if user_id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "用户ID不合法")
    ok = await UserService(db).disable_user(user_id)
    if not ok:
        raise BusinessException(ErrorCode.OPERATION_ERROR, "操作失败")
    return success(True)


@router.post("/enable", response_model=BaseResponse[bool])
async def enable_user(
    user_id: int,
    db: AsyncSession = Depends(get_db_session),
    _: User = Depends(require_role(UserRole.ADMIN)),
) -> BaseResponse[bool]:
    if user_id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "用户ID不合法")
    ok = await UserService(db).enable_user(user_id)
    if not ok:
        raise BusinessException(ErrorCode.OPERATION_ERROR, "操作失败")
    return success(True)


@router.get("/analysis", response_model=BaseResponse[UserAnalysisVO])
async def get_user_analysis(
    user_id: int,
    db: AsyncSession = Depends(get_db_session),
    _: User = Depends(require_role(UserRole.ADMIN)),
) -> BaseResponse[UserAnalysisVO]:
    if user_id <= 0:
        raise BusinessException(ErrorCode.PARAMS_ERROR, "用户ID不合法")
    user = await UserService(db).get_by_id(user_id)
    if user is None:
        raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "用户不存在")

    request_log_service = RequestLogService(db)
    analysis = UserAnalysisVO(
        userId=str(user_id),
        userAccount=user.user_account,
        userName=user.user_name,
        userStatus=user.user_status,
        userRole=user.user_role,
        tokenQuota=user.token_quota,
        usedTokens=user.used_tokens or 0,
        remainingQuota=await QuotaService(db).get_remaining_quota(user_id),
        totalRequests=await request_log_service.count_user_requests(user_id),
        successRequests=await request_log_service.count_user_success_requests(user_id),
        totalTokens=await request_log_service.count_user_tokens(user_id),
        totalCost=await BillingService(db).get_user_total_cost(user_id),
        todayCost=await BillingService(db).get_user_today_cost(user_id),
    )
    return success(analysis)
