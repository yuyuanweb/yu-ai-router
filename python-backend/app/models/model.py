"""Model entity."""

from __future__ import annotations

from datetime import datetime
from decimal import Decimal

from sqlalchemy import BigInteger, DateTime, Index, Integer, Numeric, String, func, text
from sqlalchemy.dialects.mysql import TINYINT
from sqlalchemy.orm import Mapped, mapped_column

from app.db.base import Base


class Model(Base):
    __tablename__ = "model"
    __table_args__ = (
        Index("uk_modelKey", "modelKey", unique=True),
        Index("idx_providerId", "providerId"),
        Index("idx_modelType", "modelType"),
        Index("idx_status", "status"),
        Index("idx_healthStatus", "healthStatus"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    provider_id: Mapped[int] = mapped_column("providerId", BigInteger, nullable=False)
    model_key: Mapped[str] = mapped_column("modelKey", String(128), nullable=False)
    model_name: Mapped[str] = mapped_column("modelName", String(128), nullable=False)
    model_type: Mapped[str] = mapped_column(
        "modelType", String(32), nullable=False, default="chat", server_default=text("'chat'")
    )
    description: Mapped[str | None] = mapped_column(String(512), nullable=True)
    context_length: Mapped[int] = mapped_column("contextLength", Integer, nullable=False, default=4096, server_default=text("4096"))
    input_price: Mapped[Decimal] = mapped_column(
        "inputPrice", Numeric(10, 6), nullable=False, default=Decimal("0"), server_default=text("0")
    )
    output_price: Mapped[Decimal] = mapped_column(
        "outputPrice", Numeric(10, 6), nullable=False, default=Decimal("0"), server_default=text("0")
    )
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
    avg_latency: Mapped[int] = mapped_column("avgLatency", Integer, nullable=False, default=0, server_default=text("0"))
    success_rate: Mapped[Decimal] = mapped_column(
        "successRate", Numeric(5, 2), nullable=False, default=Decimal("100.00"), server_default=text("100.00")
    )
    score: Mapped[Decimal] = mapped_column(
        Numeric(10, 4), nullable=False, default=Decimal("0"), server_default=text("0")
    )
    priority: Mapped[int] = mapped_column(Integer, nullable=False, default=100, server_default=text("100"))
    default_timeout: Mapped[int] = mapped_column(
        "defaultTimeout", Integer, nullable=False, default=60000, server_default=text("60000")
    )
    support_reasoning: Mapped[int] = mapped_column(
        "supportReasoning", TINYINT, nullable=False, default=0, server_default=text("0")
    )
    capabilities: Mapped[str | None] = mapped_column(String(512), nullable=True)
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
