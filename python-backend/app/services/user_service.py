"""User service."""

from __future__ import annotations

from math import ceil

from sqlalchemy import Select, func, select
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.constants import (
    DEFAULT_PAGE_NUM,
    DEFAULT_PAGE_SIZE,
    DEFAULT_USER_NAME,
    ErrorCode,
    MAX_PAGE_SIZE,
    UserRole,
)
from app.core.security import encrypt_password
from app.exceptions.business_exception import BusinessException
from app.models.user import User
from app.schemas.common import PageData
from app.schemas.user import LoginUserVO, UserVO


class UserService:
    def __init__(self, db: AsyncSession) -> None:
        self.db = db

    async def register(self, user_account: str, user_password: str, check_password: str) -> int:
        self._validate_register_params(user_account, user_password, check_password)
        exists_stmt = self._base_query().where(User.user_account == user_account)
        exists = await self.db.scalar(select(func.count()).select_from(exists_stmt.subquery()))
        if exists and exists > 0:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "账号重复")

        user = User(
            user_account=user_account,
            user_password=encrypt_password(user_password),
            user_name=DEFAULT_USER_NAME,
            user_role=UserRole.USER.value,
        )
        self.db.add(user)
        await self.db.commit()
        await self.db.refresh(user)
        return user.id

    async def login(self, user_account: str, user_password: str) -> User:
        if not user_account or not user_password:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "参数为空")
        if len(user_account) < 4:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "账号长度过短")
        if len(user_password) < 8:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "密码长度过短")

        stmt = self._base_query().where(
            User.user_account == user_account,
            User.user_password == encrypt_password(user_password),
        )
        user = await self.db.scalar(stmt)
        if user is None:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "用户不存在或密码错误")
        return user

    async def get_by_id(self, user_id: int) -> User | None:
        stmt = self._base_query().where(User.id == user_id)
        return await self.db.scalar(stmt)

    async def add_user(
        self,
        user_account: str,
        user_password: str,
        user_name: str | None,
        user_avatar: str | None,
        user_profile: str | None,
        user_role: str | None,
    ) -> int:
        exists_stmt = self._base_query().where(User.user_account == user_account)
        exists = await self.db.scalar(select(func.count()).select_from(exists_stmt.subquery()))
        if exists and exists > 0:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "账号重复")
        user = User(
            user_account=user_account,
            user_password=user_password,
            user_name=user_name or DEFAULT_USER_NAME,
            user_avatar=user_avatar,
            user_profile=user_profile,
            user_role=user_role or UserRole.USER.value,
        )
        self.db.add(user)
        await self.db.commit()
        await self.db.refresh(user)
        return user.id

    async def update_user(
        self,
        user_id: int,
        user_name: str | None,
        user_avatar: str | None,
        user_profile: str | None,
        user_role: str | None,
    ) -> bool:
        user = await self.get_by_id(user_id)
        if user is None:
            raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "用户不存在")
        user.user_name = user_name
        user.user_avatar = user_avatar
        user.user_profile = user_profile
        if user_role:
            user.user_role = user_role
        await self.db.commit()
        return True

    async def delete_user(self, user_id: int) -> bool:
        user = await self.get_by_id(user_id)
        if user is None:
            raise BusinessException(ErrorCode.NOT_FOUND_ERROR, "用户不存在")
        user.is_delete = 1
        await self.db.commit()
        return True

    async def list_user_vo_page(
        self,
        page_num: int,
        page_size: int,
        user_id: int | None = None,
        user_account: str | None = None,
        user_name: str | None = None,
        user_profile: str | None = None,
        user_role: str | None = None,
        sort_field: str | None = None,
        sort_order: str | None = None,
    ) -> PageData[UserVO]:
        page_num = page_num if page_num > 0 else DEFAULT_PAGE_NUM
        page_size = page_size if page_size > 0 else DEFAULT_PAGE_SIZE
        page_size = min(page_size, MAX_PAGE_SIZE)
        stmt = self._base_query()
        if user_id is not None:
            stmt = stmt.where(User.id == user_id)
        if user_role:
            stmt = stmt.where(User.user_role == user_role)
        if user_account:
            stmt = stmt.where(User.user_account.like(f"%{user_account}%"))
        if user_name:
            stmt = stmt.where(User.user_name.like(f"%{user_name}%"))
        if user_profile:
            stmt = stmt.where(User.user_profile.like(f"%{user_profile}%"))
        stmt = self._apply_sort(stmt, sort_field, sort_order)

        count = await self.db.scalar(select(func.count()).select_from(stmt.subquery())) or 0
        total_page = ceil(count / page_size) if count else 0
        paged_stmt = stmt.offset((page_num - 1) * page_size).limit(page_size)
        rows = (await self.db.scalars(paged_stmt)).all()
        records = [self.to_user_vo(item) for item in rows]
        return PageData[UserVO](
            records=records,
            pageNumber=page_num,
            pageSize=page_size,
            totalPage=total_page,
            totalRow=count,
            optimizeCountQuery=True,
        )

    @staticmethod
    def to_login_user_vo(user: User) -> LoginUserVO:
        return LoginUserVO.model_validate(user)

    @staticmethod
    def to_user_vo(user: User) -> UserVO:
        return UserVO.model_validate(user)

    def _base_query(self) -> Select[tuple[User]]:
        return select(User).where(User.is_delete == 0)

    @staticmethod
    def _apply_sort(stmt: Select[tuple[User]], sort_field: str | None, sort_order: str | None):
        sort_map = {
            "id": User.id,
            "userAccount": User.user_account,
            "userName": User.user_name,
            "userProfile": User.user_profile,
            "userRole": User.user_role,
            "createTime": User.create_time,
            "updateTime": User.update_time,
        }
        order_field = sort_map.get(sort_field or "")
        if order_field is None:
            return stmt
        return stmt.order_by(order_field.asc() if sort_order == "ascend" else order_field.desc())

    @staticmethod
    def _validate_register_params(user_account: str, user_password: str, check_password: str) -> None:
        if not user_account or not user_password or not check_password:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "参数为空")
        if len(user_account) < 4:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "账号长度过短")
        if len(user_password) < 8 or len(check_password) < 8:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "密码长度过短")
        if user_password != check_password:
            raise BusinessException(ErrorCode.PARAMS_ERROR, "两次输入的密码不一致")
