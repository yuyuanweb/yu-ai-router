"""Billing record model."""

from __future__ import annotations

from datetime import datetime
from decimal import Decimal

from sqlalchemy import BigInteger, DateTime, Index, Numeric, String, func, text
from sqlalchemy.orm import Mapped, mapped_column

from app.db.base import Base


class BillingRecord(Base):
    __tablename__ = "billing_record"
    __table_args__ = (
        Index("idx_billingType", "billingType"),
        Index("idx_createTime", "createTime"),
        Index("idx_requestLogId", "requestLogId"),
        Index("idx_userId", "userId"),
    )

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True, autoincrement=True)
    user_id: Mapped[int] = mapped_column("userId", BigInteger, nullable=False)
    request_log_id: Mapped[int | None] = mapped_column("requestLogId", BigInteger, nullable=True)
    amount: Mapped[Decimal] = mapped_column(Numeric(12, 4), nullable=False)
    balance_before: Mapped[Decimal] = mapped_column("balanceBefore", Numeric(12, 4), nullable=False)
    balance_after: Mapped[Decimal] = mapped_column("balanceAfter", Numeric(12, 4), nullable=False)
    description: Mapped[str | None] = mapped_column(String(512), nullable=True)
    billing_type: Mapped[str] = mapped_column(
        "billingType",
        String(32),
        nullable=False,
        default="api_call",
        server_default=text("'api_call'"),
    )
    create_time: Mapped[datetime] = mapped_column(
        "createTime",
        DateTime,
        nullable=False,
        default=datetime.utcnow,
        server_default=func.current_timestamp(),
    )
