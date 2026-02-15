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
    trace_id: Mapped[str | None] = mapped_column("traceId", String(64), nullable=True)
    user_id: Mapped[int | None] = mapped_column("userId", BigInteger, nullable=True)
    api_key_id: Mapped[int | None] = mapped_column("apiKeyId", BigInteger, nullable=True)
    model_id: Mapped[int | None] = mapped_column("modelId", BigInteger, nullable=True)
    request_model: Mapped[str | None] = mapped_column("requestModel", String(128), nullable=True)
    model_name: Mapped[str] = mapped_column("modelName", String(128), nullable=False)
    request_type: Mapped[str] = mapped_column("requestType", String(32), nullable=False, default="chat")
    source: Mapped[str] = mapped_column(String(32), nullable=False, default="web")
    prompt_tokens: Mapped[int] = mapped_column("promptTokens", Integer, nullable=False, default=0)
    completion_tokens: Mapped[int] = mapped_column("completionTokens", Integer, nullable=False, default=0)
    total_tokens: Mapped[int] = mapped_column("totalTokens", Integer, nullable=False, default=0)
    duration: Mapped[int] = mapped_column(Integer, nullable=False, default=0)
    status: Mapped[str] = mapped_column(String(32), nullable=False, default="success")
    error_message: Mapped[str | None] = mapped_column("errorMessage", Text, nullable=True)
    error_code: Mapped[str | None] = mapped_column("errorCode", String(64), nullable=True)
    routing_strategy: Mapped[str | None] = mapped_column("routingStrategy", String(32), nullable=True)
    is_fallback: Mapped[int] = mapped_column("isFallback", Integer, nullable=False, default=0)
    client_ip: Mapped[str | None] = mapped_column("clientIp", String(64), nullable=True)
    user_agent: Mapped[str | None] = mapped_column("userAgent", String(512), nullable=True)
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
