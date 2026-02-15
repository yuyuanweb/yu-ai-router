"""Stripe payment service."""

from __future__ import annotations

from decimal import Decimal
from typing import Any

from app.core.config import get_settings
from app.core.constants import ErrorCode, PAYMENT_METHOD_STRIPE
from app.exceptions.business_exception import BusinessException
from app.services.recharge_service import RechargeService

try:
    import stripe
except Exception:  # pragma: no cover
    stripe = None


class StripePaymentService:
    def __init__(self, recharge_service: RechargeService) -> None:
        self.recharge_service = recharge_service
        self.settings = get_settings()
        self._ensure_client_ready()

    def _ensure_client_ready(self) -> None:
        if stripe is None:
            raise BusinessException(ErrorCode.OPERATION_ERROR, "未安装 stripe SDK")
        if not self.settings.stripe_api_key:
            raise BusinessException(ErrorCode.OPERATION_ERROR, "Stripe API Key 未配置")
        stripe.api_key = self.settings.stripe_api_key

    async def create_checkout_session(
        self,
        user_id: int,
        amount: Decimal,
        success_url: str,
        cancel_url: str,
    ) -> dict[str, Any]:
        record = await self.recharge_service.create_recharge_record(
            user_id=user_id,
            amount=amount,
            payment_method=PAYMENT_METHOD_STRIPE,
        )
        amount_in_cents = int((amount * Decimal("100")).to_integral_value())
        try:
            session = stripe.checkout.Session.create(
                mode="payment",
                success_url=f"{success_url}?session_id={{CHECKOUT_SESSION_ID}}",
                cancel_url=cancel_url,
                line_items=[
                    {
                        "price_data": {
                            "currency": "cny",
                            "unit_amount": amount_in_cents,
                            "product_data": {
                                "name": "Yu AI Router 账户充值",
                                "description": f"充值金额：¥{amount}",
                            },
                        },
                        "quantity": 1,
                    }
                ],
                metadata={
                    "userId": str(user_id),
                    "recordId": str(record.id),
                    "amount": str(amount),
                },
            )
            return {"id": session.id, "url": session.url}
        except Exception as exc:
            raise BusinessException(ErrorCode.OPERATION_ERROR, f"创建支付会话失败：{exc}") from exc

    async def handle_payment_success(self, session_id: str) -> bool:
        try:
            session = stripe.checkout.Session.retrieve(session_id)
        except Exception as exc:
            raise BusinessException(ErrorCode.OPERATION_ERROR, "获取支付会话失败") from exc
        if session.get("payment_status") != "paid":
            return False
        metadata = session.get("metadata") or {}
        record_id = metadata.get("recordId")
        if not record_id:
            return False
        return await self.recharge_service.complete_recharge(int(record_id), session_id)

    def construct_webhook_event(self, payload: str, sig_header: str) -> dict[str, Any]:
        if not self.settings.stripe_webhook_secret:
            raise BusinessException(ErrorCode.OPERATION_ERROR, "Stripe Webhook Secret 未配置")
        try:
            event = stripe.Webhook.construct_event(payload, sig_header, self.settings.stripe_webhook_secret)
            return dict(event)
        except Exception as exc:
            raise BusinessException(ErrorCode.FORBIDDEN_ERROR, f"Webhook 签名验证失败: {exc}") from exc
