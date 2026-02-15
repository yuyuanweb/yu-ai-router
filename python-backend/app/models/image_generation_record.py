"""Image generation record model."""

from __future__ import annotations

from datetime import datetime
from decimal import Decimal

from sqlalchemy import BigInteger, DateTime, Index, Integer, Numeric, String, Text, func, text
from sqlalchemy.orm import Mapped, mapped_column

from app.db.base import Base


class ImageGenerationRecord(Base):
    __tablename__ = "image_generation_record"
    __table_args__ = (
        Index("idx_userId", "userId"),
        Index("idx_modelId", "modelId"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    user_id: Mapped[int] = mapped_column("userId", BigInteger, nullable=False)
    api_key_id: Mapped[int | None] = mapped_column("apiKeyId", BigInteger, nullable=True)
    model_id: Mapped[int] = mapped_column("modelId", BigInteger, nullable=False)
    model_key: Mapped[str] = mapped_column("modelKey", String(128), nullable=False)
    prompt: Mapped[str] = mapped_column(Text, nullable=False)
    revised_prompt: Mapped[str | None] = mapped_column("revisedPrompt", Text, nullable=True)
    image_url: Mapped[str | None] = mapped_column("imageUrl", String(1024), nullable=True)
    image_data: Mapped[str | None] = mapped_column("imageData", Text, nullable=True)
    size: Mapped[str | None] = mapped_column(String(32), nullable=True)
    quality: Mapped[str | None] = mapped_column(String(32), nullable=True)
    status: Mapped[str] = mapped_column(
        String(32),
        nullable=False,
        default="success",
        server_default=text("'success'"),
    )
    cost: Mapped[Decimal] = mapped_column(
        Numeric(12, 4),
        nullable=False,
        default=Decimal("0"),
        server_default=text("0.0000"),
    )
    duration: Mapped[int | None] = mapped_column(Integer, nullable=True)
    error_message: Mapped[str | None] = mapped_column("errorMessage", String(512), nullable=True)
    client_ip: Mapped[str | None] = mapped_column("clientIp", String(128), nullable=True)
    create_time: Mapped[datetime] = mapped_column(
        "createTime",
        DateTime,
        nullable=False,
        default=datetime.utcnow,
        server_default=func.current_timestamp(),
    )
