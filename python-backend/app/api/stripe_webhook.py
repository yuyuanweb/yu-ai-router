"""Stripe webhook API."""

from __future__ import annotations

from fastapi import APIRouter, Depends, Header, Request
from sqlalchemy.ext.asyncio import AsyncSession

from app.db.session import get_db_session
from app.services.recharge_service import RechargeService
from app.services.stripe_payment_service import StripePaymentService

router = APIRouter(prefix="/webhook/stripe", tags=["stripe-webhook"])


@router.post("", response_model=None)
async def handle_stripe_webhook(
    request: Request,
    stripe_signature: str = Header(alias="Stripe-Signature"),
    db: AsyncSession = Depends(get_db_session),
) -> str:
    payload = (await request.body()).decode("utf-8")
    stripe_service = StripePaymentService(RechargeService(db))
    try:
        event = stripe_service.construct_webhook_event(payload, stripe_signature)
        event_type = event.get("type")
        if event_type == "checkout.session.completed":
            session = ((event.get("data") or {}).get("object")) or {}
            session_id = session.get("id")
            if session_id:
                await stripe_service.handle_payment_success(session_id)
        return "success"
    except Exception:
        return "error"
