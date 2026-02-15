"""User model."""

from __future__ import annotations

from datetime import datetime

from sqlalchemy import BigInteger, DateTime, Index, String, func, text
from sqlalchemy.dialects.mysql import TINYINT
from sqlalchemy.orm import Mapped, mapped_column

from app.db.base import Base


class User(Base):
    __tablename__ = "user"
    __table_args__ = (
        Index("uk_userAccount", "userAccount", unique=True),
        Index("idx_userName", "userName"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    user_account: Mapped[str] = mapped_column("userAccount", String(256), nullable=False)
    user_password: Mapped[str] = mapped_column("userPassword", String(512), nullable=False)
    user_name: Mapped[str | None] = mapped_column("userName", String(256), nullable=True)
    user_avatar: Mapped[str | None] = mapped_column("userAvatar", String(1024), nullable=True)
    user_profile: Mapped[str | None] = mapped_column("userProfile", String(512), nullable=True)
    user_role: Mapped[str] = mapped_column(
        "userRole",
        String(256),
        nullable=False,
        default="user",
        server_default=text("'user'"),
    )
    user_status: Mapped[str] = mapped_column(
        "userStatus",
        String(32),
        nullable=False,
        default="active",
        server_default=text("'active'"),
    )
    token_quota: Mapped[int] = mapped_column(
        "tokenQuota",
        BigInteger,
        nullable=False,
        default=-1,
        server_default=text("-1"),
    )
    used_tokens: Mapped[int] = mapped_column(
        "usedTokens",
        BigInteger,
        nullable=False,
        default=0,
        server_default=text("0"),
    )
    edit_time: Mapped[datetime] = mapped_column(
        "editTime",
        DateTime,
        nullable=False,
        default=datetime.utcnow,
        server_default=func.current_timestamp(),
    )
    create_time: Mapped[datetime] = mapped_column(
        "createTime",
        DateTime,
        nullable=False,
        default=datetime.utcnow,
        server_default=func.current_timestamp(),
    )
    update_time: Mapped[datetime] = mapped_column(
        "updateTime",
        DateTime,
        nullable=False,
        default=datetime.utcnow,
        onupdate=datetime.utcnow,
        server_default=func.current_timestamp(),
    )
    is_delete: Mapped[int] = mapped_column(
        "isDelete",
        TINYINT,
        nullable=False,
        default=0,
        server_default=text("0"),
    )
