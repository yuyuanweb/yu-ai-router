"""User provider key service (BYOK)."""

from __future__ import annotations

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import BYOK_STATUS_ACTIVE, ErrorCode
from app.exceptions.business_exception import BusinessException
from app.models.model_provider import ModelProvider
from app.models.user_provider_key import UserProviderKey
from app.schemas.byok import UserProviderKeyAddRequest, UserProviderKeyUpdateRequest, UserProviderKeyVO
from app.utils.encryption import decrypt_text, encrypt_text


class UserProviderKeyService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db

    async def add_user_provider_key(self, request: UserProviderKeyAddRequest, user_id: int) -> bool:
        if request.provider_id is None or not request.api_key:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "请求参数错误")

        provider = await self.db.scalar(
            select(ModelProvider).where(
                ModelProvider.id == request.provider_id,
                ModelProvider.is_delete == 0,
            )
        )
        if provider is None:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "提供者不存在")

        exists = await self.db.scalar(
            select(UserProviderKey).where(
                UserProviderKey.user_id == user_id,
                UserProviderKey.provider_id == request.provider_id,
                UserProviderKey.is_delete == 0,
            )
        )
        if exists is not None:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "已配置过该提供者的密钥，请使用更新功能")

        entity = UserProviderKey(
            user_id=user_id,
            provider_id=request.provider_id,
            provider_name=provider.provider_name,
            api_key=encrypt_text(request.api_key),
            status=BYOK_STATUS_ACTIVE,
        )
        self.db.add(entity)
        await self.db.commit()
        return True

    async def update_user_provider_key(self, request: UserProviderKeyUpdateRequest, user_id: int) -> bool:
        entity = await self.db.scalar(
            select(UserProviderKey).where(
                UserProviderKey.id == request.id,
                UserProviderKey.is_delete == 0,
            )
        )
        if entity is None:
            raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "请求数据不存在")
        if entity.user_id != user_id:
            raise BusinessException(ErrorCode.NO_AUTH_ERROR, "无权限")

        if request.api_key:
            entity.api_key = encrypt_text(request.api_key)
        if request.status:
            entity.status = request.status
        await self.db.commit()
        return True

    async def delete_user_provider_key(self, key_id: int, user_id: int) -> bool:
        entity = await self.db.scalar(
            select(UserProviderKey).where(
                UserProviderKey.id == key_id,
                UserProviderKey.is_delete == 0,
            )
        )
        if entity is None:
            raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "请求数据不存在")
        if entity.user_id != user_id:
            raise BusinessException(ErrorCode.NO_AUTH_ERROR, "无权限")

        entity.is_delete = 1
        await self.db.commit()
        return True

    async def list_user_provider_keys(self, user_id: int) -> list[UserProviderKeyVO]:
        rows = (
            await self.db.scalars(
                select(UserProviderKey)
                .where(
                    UserProviderKey.user_id == user_id,
                    UserProviderKey.is_delete == 0,
                )
                .order_by(UserProviderKey.create_time.desc())
            )
        ).all()
        return [self.to_vo(row) for row in rows]

    async def get_user_provider_api_key(self, user_id: int, provider_id: int) -> str | None:
        entity = await self.db.scalar(
            select(UserProviderKey).where(
                UserProviderKey.user_id == user_id,
                UserProviderKey.provider_id == provider_id,
                UserProviderKey.status == BYOK_STATUS_ACTIVE,
                UserProviderKey.is_delete == 0,
            )
        )
        if entity is None:
            return None
        try:
            return decrypt_text(entity.api_key)
        except Exception:
            return None

    async def has_user_provider_key(self, user_id: int, provider_id: int) -> bool:
        return await self.get_user_provider_api_key(user_id, provider_id) is not None

    @staticmethod
    def to_vo(entity: UserProviderKey) -> UserProviderKeyVO:
        masked = "****"
        if entity.api_key and len(entity.api_key) > 12:
            masked = f"{entity.api_key[:8]}****{entity.api_key[-4:]}"
        return UserProviderKeyVO(
            id=entity.id,
            providerId=entity.provider_id,
            providerName=entity.provider_name,
            apiKey=masked,
            status=entity.status,
            createTime=entity.create_time,
            updateTime=entity.update_time,
        )
