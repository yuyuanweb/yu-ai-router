"""Model provider model."""

from __future__ import annotations

from datetime import datetime
from decimal import Decimal

from sqlalchemy import BigInteger, DateTime, Index, Numeric, String, Text, func, text
from sqlalchemy.dialects.mysql import TINYINT
from sqlalchemy.orm import Mapped, mapped_column

from app.db.base import Base


class ModelProvider(Base):
    __tablename__ = "model_provider"
    __table_args__ = (
        Index("uk_providerName", "providerName", unique=True),
        Index("idx_status", "status"),
        Index("idx_healthStatus", "healthStatus"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    provider_name: Mapped[str] = mapped_column("providerName", String(64), nullable=False)
    display_name: Mapped[str] = mapped_column("displayName", String(128), nullable=False)
    base_url: Mapped[str] = mapped_column("baseUrl", String(512), nullable=False)
    api_key: Mapped[str] = mapped_column("apiKey", String(512), nullable=False)
    status: Mapped[str] = mapped_column(
        String(32), nullable=False, default="active", server_default=text("'active'")
    )
    health_status: Mapped[str] = mapped_column(
        "healthStatus",
        String(32),
        nullable=False,
        default="unknown",
        server_default=text("'unknown'"),
    )
    avg_latency: Mapped[int] = mapped_column("avgLatency", nullable=False, default=0, server_default=text("0"))
    success_rate: Mapped[Decimal] = mapped_column(
        "successRate", Numeric(5, 2), nullable=False, default=Decimal("100.00"), server_default=text("100.00")
    )
    priority: Mapped[int] = mapped_column(nullable=False, default=100, server_default=text("100"))
    config: Mapped[str | None] = mapped_column(Text, nullable=True)
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
