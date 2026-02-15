"""Image generation API."""

from __future__ import annotations

from fastapi import APIRouter, Depends, Header, Query, Request
from redis.asyncio import Redis
from sqlalchemy.ext.asyncio import AsyncSession

from app.common.result_utils import success
from app.core.constants import ErrorCode
from app.db.session import get_db_session
from app.exceptions.business_exception import BusinessException
from app.infra.redis_client import get_redis_client
from app.middleware.auth import get_login_user, require_login
from app.models.user import User
from app.schemas.common import BaseResponse, PageData
from app.schemas.image import ImageGenerationRecordVO, ImageGenerationRequest, ImageGenerationResponse
from app.services.api_key_service import ApiKeyService
from app.services.image_generation_service import ImageGenerationService
from app.utils.request import get_client_ip

router = APIRouter(prefix="/v1/images", tags=["image"])


@router.post("/generations", response_model=ImageGenerationResponse)
async def generate_image(
    payload: ImageGenerationRequest,
    request: Request,
    authorization: str | None = Header(default=None, alias="Authorization"),
    db: AsyncSession = Depends(get_db_session),
    redis: Redis = Depends(get_redis_client),
) -> ImageGenerationResponse:
    user_id: int | None = None
    api_key_id: int | None = None

    if authorization and authorization.startswith("Bearer "):
        api_key_value = authorization[7:]
        api_key = await ApiKeyService(db).get_by_key_value(api_key_value)
        if api_key is None:
            raise BusinessException(ErrorCode.NO_AUTH_ERROR, "API Key 无效或已失效")
        user_id = api_key.user_id
        api_key_id = api_key.id
    else:
        try:
            login_user = await get_login_user(request, db, redis)
            user_id = login_user.id
        except BusinessException as exc:
            raise BusinessException(ErrorCode.NOT_LOGIN_ERROR, "请先登录或提供有效的 API Key") from exc

    if user_id is None:
        raise BusinessException(ErrorCode.NOT_LOGIN_ERROR, "请先登录或提供有效的 API Key")

    return await ImageGenerationService(db).generate_image(
        request=payload,
        user_id=user_id,
        api_key_id=api_key_id,
        client_ip=get_client_ip(request),
    )


@router.get("/my/records", response_model=BaseResponse[PageData[ImageGenerationRecordVO]])
async def get_my_records(
    page_num: int = Query(default=1, alias="pageNum"),
    page_size: int = Query(default=10, alias="pageSize"),
    login_user: User = Depends(require_login),
    db: AsyncSession = Depends(get_db_session),
) -> BaseResponse[PageData[ImageGenerationRecordVO]]:
    data = await ImageGenerationService(db).list_user_records(login_user.id, page_num, page_size)
    return success(data)
