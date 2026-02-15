"""User schemas."""

from __future__ import annotations

from datetime import datetime

from pydantic import Field

from app.schemas.common import CamelBaseModel, LongIdModel, PageRequest


class UserRegisterRequest(CamelBaseModel):
    user_account: str = Field(alias="userAccount")
    user_password: str = Field(alias="userPassword")
    check_password: str = Field(alias="checkPassword")


class UserLoginRequest(CamelBaseModel):
    user_account: str = Field(alias="userAccount")
    user_password: str = Field(alias="userPassword")


class UserAddRequest(CamelBaseModel):
    user_name: str | None = Field(default=None, alias="userName")
    user_account: str
    user_avatar: str | None = Field(default=None, alias="userAvatar")
    user_profile: str | None = Field(default=None, alias="userProfile")
    user_role: str | None = Field(default=None, alias="userRole")


class UserUpdateRequest(CamelBaseModel):
    id: int
    user_name: str | None = Field(default=None, alias="userName")
    user_avatar: str | None = Field(default=None, alias="userAvatar")
    user_profile: str | None = Field(default=None, alias="userProfile")
    user_role: str | None = Field(default=None, alias="userRole")


class UserQueryRequest(PageRequest):
    id: int | None = None
    user_name: str | None = Field(default=None, alias="userName")
    user_account: str | None = Field(default=None, alias="userAccount")
    user_profile: str | None = Field(default=None, alias="userProfile")
    user_role: str | None = Field(default=None, alias="userRole")


class LoginUserVO(LongIdModel):
    user_account: str | None = Field(default=None, alias="userAccount")
    user_name: str | None = Field(default=None, alias="userName")
    user_avatar: str | None = Field(default=None, alias="userAvatar")
    user_profile: str | None = Field(default=None, alias="userProfile")
    user_role: str | None = Field(default=None, alias="userRole")
    create_time: datetime | None = Field(default=None, alias="createTime")
    update_time: datetime | None = Field(default=None, alias="updateTime")


class UserVO(LongIdModel):
    user_account: str | None = Field(default=None, alias="userAccount")
    user_name: str | None = Field(default=None, alias="userName")
    user_avatar: str | None = Field(default=None, alias="userAvatar")
    user_profile: str | None = Field(default=None, alias="userProfile")
    user_role: str | None = Field(default=None, alias="userRole")
    create_time: datetime | None = Field(default=None, alias="createTime")


class UserRawVO(LongIdModel):
    user_account: str | None = Field(default=None, alias="userAccount")
    user_password: str | None = Field(default=None, alias="userPassword")
    user_name: str | None = Field(default=None, alias="userName")
    user_avatar: str | None = Field(default=None, alias="userAvatar")
    user_profile: str | None = Field(default=None, alias="userProfile")
    user_role: str | None = Field(default=None, alias="userRole")
    edit_time: datetime | None = Field(default=None, alias="editTime")
    create_time: datetime | None = Field(default=None, alias="createTime")
    update_time: datetime | None = Field(default=None, alias="updateTime")
    is_delete: int | None = Field(default=None, alias="isDelete")
