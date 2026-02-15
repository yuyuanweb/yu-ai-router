"""User provider key (BYOK) model."""

from __future__ import annotations

from datetime import datetime

from sqlalchemy import BigInteger, DateTime, Index, String, func, text
from sqlalchemy.dialects.mysql import TINYINT
from sqlalchemy.orm import Mapped, mapped_column

from app.db.base import Base


class UserProviderKey(Base):
    __tablename__ = "user_provider_key"
    __table_args__ = (
        Index("uk_user_provider", "userId", "providerId"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    user_id: Mapped[int] = mapped_column("userId", BigInteger, nullable=False)
    provider_id: Mapped[int] = mapped_column("providerId", BigInteger, nullable=False)
    provider_name: Mapped[str] = mapped_column("providerName", String(64), nullable=False)
    api_key: Mapped[str] = mapped_column("apiKey", String(512), nullable=False)
    status: Mapped[str] = mapped_column(
        String(32), nullable=False, default="active", server_default=text("'active'")
    )
    create_time: Mapped[datetime] = mapped_column(
        "createTime", DateTime, nullable=False, default=datetime.utcnow, server_default=func.current_timestamp()
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
        "isDelete", TINYINT, nullable=False, default=0, server_default=text("0")
    )
