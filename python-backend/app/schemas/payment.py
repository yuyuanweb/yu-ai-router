"""Payment and balance schemas."""

from __future__ import annotations

from datetime import datetime
from decimal import Decimal

from pydantic import Field, field_serializer

from app.schemas.common import CamelBaseModel, LongIdModel


class CreateRechargeRequest(CamelBaseModel):
    amount: Decimal


class CreateRechargeResponse(CamelBaseModel):
    checkout_url: str = Field(alias="checkoutUrl")
    session_id: str = Field(alias="sessionId")


class BalanceVO(CamelBaseModel):
    balance: Decimal
    total_spending: Decimal = Field(alias="totalSpending")
    total_recharge: Decimal = Field(alias="totalRecharge")


class RechargeRecordVO(LongIdModel):
    user_id: int = Field(alias="userId")
    amount: float
    payment_method: str = Field(alias="paymentMethod")
    payment_id: str | None = Field(default=None, alias="paymentId")
    status: str
    description: str | None = None
    create_time: datetime | None = Field(default=None, alias="createTime")
    update_time: datetime | None = Field(default=None, alias="updateTime")

    @field_serializer("user_id", when_used="json")
    def serialize_user_id(self, value: int) -> str:
        return str(value)


class BillingRecordVO(LongIdModel):
    user_id: int = Field(alias="userId")
    request_log_id: int | None = Field(default=None, alias="requestLogId")
    amount: float
    balance_before: float = Field(alias="balanceBefore")
    balance_after: float = Field(alias="balanceAfter")
    description: str | None = None
    billing_type: str = Field(alias="billingType")
    create_time: datetime | None = Field(default=None, alias="createTime")

    @field_serializer("user_id", when_used="json")
    def serialize_user_id(self, value: int) -> str:
        return str(value)

    @field_serializer("request_log_id", when_used="json")
    def serialize_request_log_id(self, value: int | None) -> str | None:
        return str(value) if value is not None else None
