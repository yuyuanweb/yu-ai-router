"""DashScope adapter."""

from __future__ import annotations

from app.adapter.openai_adapter import OpenAIAdapter


class DashscopeAdapter(OpenAIAdapter):
    supported_providers = {"qwen", "dashscope", "tongyi", "aliyun"}
