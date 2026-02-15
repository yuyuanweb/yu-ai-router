"""DeepSeek adapter."""

from __future__ import annotations

from app.adapter.openai_adapter import OpenAIAdapter


class DeepSeekAdapter(OpenAIAdapter):
    supported_providers = {"deepseek"}
