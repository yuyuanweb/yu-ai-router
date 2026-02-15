"""PDF parser plugin."""

from __future__ import annotations

import io
import json

import httpx

from app.plugin.plugin import Plugin
from app.plugin.plugin_context import PluginContext
from app.plugin.plugin_result import PluginResult


class PdfParserPlugin(Plugin):
    PLUGIN_KEY = "pdf_parser"
    PLUGIN_NAME = "PDF解析"
    DESCRIPTION = "解析PDF文档内容，提取文本信息"

    def __init__(self) -> None:
        self.max_pages = 50
        self.max_text_length = 50000

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
            self.max_pages = int(data.get("maxPages", self.max_pages))
            self.max_text_length = int(data.get("maxTextLength", self.max_text_length))
        except Exception:
            return

    def supports(self, context: PluginContext) -> bool:
        return bool(context.file_url) or bool(context.file_bytes)

    async def execute(self, context: PluginContext) -> PluginResult:
        try:
            pdf_bytes = await self._resolve_pdf_bytes(context)
            if not pdf_bytes:
                return PluginResult.fail("请提供PDF文件或文件URL")
            text = self._parse_pdf_text(pdf_bytes)
            if not text.strip():
                return PluginResult.fail("PDF文档为空或无法提取文本")
            if len(text) > self.max_text_length:
                text = text[: self.max_text_length] + "\n\n[文档内容过长，已截断...]"
            text = self._normalize_text(text)
            content_lines = []
            if context.input:
                content_lines.append(f"用户问题: {context.input}")
                content_lines.append("")
                content_lines.append("以下是PDF文档的内容，请根据内容回答用户问题：")
                content_lines.append("")
            content_lines.append(text)
            return PluginResult.success_result(
                "\n".join(content_lines).strip(),
                {
                    "textLength": len(text),
                    "truncated": len(text) >= self.max_text_length,
                },
            )
        except Exception as exc:
            return PluginResult.fail(f"PDF解析失败: {exc}")

    async def _resolve_pdf_bytes(self, context: PluginContext) -> bytes:
        if context.file_bytes:
            return context.file_bytes
        if not context.file_url:
            return b""
        async with httpx.AsyncClient(timeout=60) as client:
            resp = await client.get(context.file_url)
        if resp.status_code != 200:
            raise RuntimeError(f"下载文件失败，状态码: {resp.status_code}")
        return resp.content

    def _parse_pdf_text(self, pdf_bytes: bytes) -> str:
        parsers = (
            self._parse_with_pypdf,
            self._parse_with_pypdf2,
            self._parse_with_pdfplumber,
        )
        last_error: Exception | None = None
        for parser in parsers:
            try:
                return parser(pdf_bytes)
            except Exception as exc:  # noqa: PERF203
                last_error = exc
        raise RuntimeError(
            "未找到可用的 PDF 解析依赖，请安装 pypdf / PyPDF2 / pdfplumber 之一"
        ) from last_error

    def _parse_with_pypdf(self, pdf_bytes: bytes) -> str:
        from pypdf import PdfReader

        reader = PdfReader(io.BytesIO(pdf_bytes))
        limit = min(len(reader.pages), self.max_pages)
        chunks: list[str] = []
        for idx in range(limit):
            page_text = reader.pages[idx].extract_text() or ""
            chunks.append(page_text)
        return "\n".join(chunks)

    def _parse_with_pypdf2(self, pdf_bytes: bytes) -> str:
        from PyPDF2 import PdfReader

        reader = PdfReader(io.BytesIO(pdf_bytes))
        limit = min(len(reader.pages), self.max_pages)
        chunks: list[str] = []
        for idx in range(limit):
            page_text = reader.pages[idx].extract_text() or ""
            chunks.append(page_text)
        return "\n".join(chunks)

    def _parse_with_pdfplumber(self, pdf_bytes: bytes) -> str:
        import pdfplumber

        chunks: list[str] = []
        with pdfplumber.open(io.BytesIO(pdf_bytes)) as pdf:
            limit = min(len(pdf.pages), self.max_pages)
            for idx in range(limit):
                page_text = pdf.pages[idx].extract_text() or ""
                chunks.append(page_text)
        return "\n".join(chunks)

    @staticmethod
    def _normalize_text(text: str) -> str:
        while "\n\n\n" in text:
            text = text.replace("\n\n\n", "\n\n")
        return text.strip()
