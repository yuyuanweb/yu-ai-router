"""Plugin config model."""

from __future__ import annotations

from datetime import datetime

from sqlalchemy import BigInteger, DateTime, Index, Integer, String, Text, func, text
from sqlalchemy.dialects.mysql import TINYINT
from sqlalchemy.orm import Mapped, mapped_column

from app.db.base import Base


class PluginConfig(Base):
    __tablename__ = "plugin_config"
    __table_args__ = (
        Index("uk_pluginKey", "pluginKey", unique=True),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    plugin_key: Mapped[str] = mapped_column("pluginKey", String(64), nullable=False)
    plugin_name: Mapped[str] = mapped_column("pluginName", String(128), nullable=False)
    plugin_type: Mapped[str] = mapped_column(
        "pluginType",
        String(32),
        nullable=False,
        default="builtin",
        server_default=text("'builtin'"),
    )
    description: Mapped[str | None] = mapped_column(String(512), nullable=True)
    config: Mapped[str | None] = mapped_column(Text, nullable=True)
    status: Mapped[str] = mapped_column(
        String(32),
        nullable=False,
        default="active",
        server_default=text("'active'"),
    )
    priority: Mapped[int] = mapped_column(Integer, nullable=False, default=100, server_default=text("100"))
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
    is_delete: Mapped[int] = mapped_column("isDelete", TINYINT, nullable=False, default=0, server_default=text("0"))
