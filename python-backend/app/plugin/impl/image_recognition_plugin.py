"""Image recognition plugin."""

from __future__ import annotations

import base64
import json

import httpx
from sqlalchemy import select

from app.db.session import session_maker
from app.models.model_provider import ModelProvider
from app.plugin.plugin import Plugin
from app.plugin.plugin_context import PluginContext
from app.plugin.plugin_result import PluginResult


class ImageRecognitionPlugin(Plugin):
    PLUGIN_KEY = "image_recognition"
    PLUGIN_NAME = "图片识别"
    DESCRIPTION = "识别图片内容，返回图片描述"

    def __init__(self) -> None:
        self.model = "qwen-vl-plus"
        self.max_image_size = 4 * 1024 * 1024
        self.default_prompt = "请描述这张图片的内容，包括主要对象、场景、颜色、文字等信息。"

    def get_plugin_key(self) -> str:
        return self.PLUGIN_KEY

    def get_plugin_name(self) -> str:
        return self.PLUGIN_NAME

    def get_description(self) -> str:
        return self.DESCRIPTION

    def init(self, config: str | None) -> None:
        if not config:
            return
        try:
            data = json.loads(config)
            self.model = str(data.get("model", self.model))
            self.max_image_size = int(data.get("maxImageSize", self.max_image_size))
            self.default_prompt = str(data.get("defaultPrompt", self.default_prompt))
        except Exception:
            return

    def supports(self, context: PluginContext) -> bool:
        return bool(context.file_url) or bool(context.file_bytes)

    async def execute(self, context: PluginContext) -> PluginResult:
        try:
            image_bytes, mime_type = await self._resolve_image(context)
            if len(image_bytes) > self.max_image_size:
                return PluginResult.fail(f"图片大小超过限制（最大 {self.max_image_size // 1024 // 1024}MB）")
            prompt = (context.input or "").strip() or self.default_prompt
            content = await self._recognize(image_bytes, mime_type, prompt)
            text = content
            if context.input:
                text = f"用户问题: {context.input}\n\n{content}"
            return PluginResult.success_result(
                text,
                {"model": self.model, "imageSize": len(image_bytes)},
            )
        except Exception as exc:
            return PluginResult.fail(f"图片识别失败: {exc}")

    async def _resolve_image(self, context: PluginContext) -> tuple[bytes, str]:
        if context.file_bytes:
            return context.file_bytes, context.file_type or "image/jpeg"
        if not context.file_url:
            raise RuntimeError("请提供图片文件或图片URL")
        async with httpx.AsyncClient(timeout=30) as client:
            resp = await client.get(context.file_url)
        if resp.status_code != 200:
            raise RuntimeError(f"下载图片失败，状态码: {resp.status_code}")
        mime = self._guess_mime(context.file_url)
        return resp.content, mime

    async def _recognize(self, image_bytes: bytes, mime_type: str, prompt: str) -> str:
        provider = await self._get_qwen_provider()
        if provider is None:
            raise RuntimeError("未找到通义千问提供者配置")
        image_url = f"data:{mime_type};base64,{base64.b64encode(image_bytes).decode('utf-8')}"
        payload = {
            "model": self.model,
            "messages": [
                {
                    "role": "user",
                    "content": [
                        {"type": "image_url", "image_url": {"url": image_url}},
                        {"type": "text", "text": prompt},
                    ],
                }
            ],
        }
        url = f"{provider.base_url.rstrip('/')}/v1/chat/completions"
        headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {provider.api_key}",
        }
        async with httpx.AsyncClient(timeout=60) as client:
            resp = await client.post(url, headers=headers, json=payload)
        if resp.status_code != 200:
            raise RuntimeError(f"视觉模型调用失败，状态码: {resp.status_code}")
        data = resp.json()
        choices = data.get("choices") or []
        if not choices:
            raise RuntimeError("模型返回结果格式错误")
        message = choices[0].get("message") or {}
        return str(message.get("content") or "无法识别图片内容")

    @staticmethod
    async def _get_qwen_provider() -> ModelProvider | None:
        async with session_maker() as db:
            stmt = select(ModelProvider).where(
                ModelProvider.provider_name == "qwen",
                ModelProvider.is_delete == 0,
            )
            return await db.scalar(stmt)

    @staticmethod
    def _guess_mime(url: str) -> str:
        lowered = url.lower()
        if ".png" in lowered:
            return "image/png"
        if ".gif" in lowered:
            return "image/gif"
        if ".webp" in lowered:
            return "image/webp"
        return "image/jpeg"
