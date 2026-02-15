"""Zhipu adapter."""

from __future__ import annotations

from app.adapter.openai_adapter import OpenAIAdapter


class ZhipuAdapter(OpenAIAdapter):
    supported_providers = {"zhipu", "zhipuai", "glm"}
