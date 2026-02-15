"""Model adapter factory."""

from __future__ import annotations

from app.adapter.dashscope_adapter import DashscopeAdapter
from app.adapter.deepseek_adapter import DeepSeekAdapter
from app.adapter.model_adapter import ModelAdapter
from app.adapter.openai_adapter import OpenAIAdapter
from app.adapter.zhipu_adapter import ZhipuAdapter


class ModelAdapterFactory:
    def __init__(self) -> None:
        self.adapters: list[ModelAdapter] = [
            DashscopeAdapter(),
            ZhipuAdapter(),
            DeepSeekAdapter(),
            OpenAIAdapter(),
        ]

    def get_adapter(self, provider_name: str) -> ModelAdapter:
        normalized = provider_name.strip().lower()
        for adapter in self.adapters:
            if adapter.supports(normalized):
                return adapter
        return OpenAIAdapter()
