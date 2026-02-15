"""Request log model."""

from __future__ import annotations

from datetime import datetime

from sqlalchemy import BigInteger, DateTime, Index, Integer, String, Text, func
from sqlalchemy.orm import Mapped, mapped_column

from app.db.base import Base


class RequestLog(Base):
    __tablename__ = "request_log"
    __table_args__ = (
        Index("idx_userId", "userId"),
        Index("idx_apiKeyId", "apiKeyId"),
        Index("idx_createTime", "createTime"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    user_id: Mapped[int | None] = mapped_column("userId", BigInteger, nullable=True)
    api_key_id: Mapped[int | None] = mapped_column("apiKeyId", BigInteger, nullable=True)
    model_name: Mapped[str] = mapped_column("modelName", String(128), nullable=False)
    prompt_tokens: Mapped[int] = mapped_column("promptTokens", Integer, nullable=False, default=0)
    completion_tokens: Mapped[int] = mapped_column("completionTokens", Integer, nullable=False, default=0)
    total_tokens: Mapped[int] = mapped_column("totalTokens", Integer, nullable=False, default=0)
    duration: Mapped[int] = mapped_column(Integer, nullable=False, default=0)
    status: Mapped[str] = mapped_column(String(32), nullable=False, default="success")
    error_message: Mapped[str | None] = mapped_column("errorMessage", Text, nullable=True)
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
