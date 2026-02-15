"""Web search plugin."""

from __future__ import annotations

import json
from urllib.parse import quote_plus

import httpx

from app.core.config import get_settings
from app.plugin.plugin import Plugin
from app.plugin.plugin_context import PluginContext
from app.plugin.plugin_result import PluginResult


class WebSearchPlugin(Plugin):
    PLUGIN_KEY = "web_search"
    PLUGIN_NAME = "Web搜索"
    DESCRIPTION = "实时联网搜索，获取最新信息（SerpApi）"
    SERPAPI_ENDPOINT = "https://serpapi.com/search.json"

    def __init__(self) -> None:
        self.settings = get_settings()
        self.max_results = 5
        self.search_engine = "google"
        self.timeout = 15000

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
            self.max_results = int(data.get("maxResults", self.max_results))
            self.search_engine = str(data.get("searchEngine", self.search_engine))
            self.timeout = int(data.get("timeout", self.timeout))
        except Exception:
            return

    def supports(self, context: PluginContext) -> bool:
        return bool(context.input and context.input.strip())

    async def execute(self, context: PluginContext) -> PluginResult:
        query = (context.input or "").strip()
        if not query:
            return PluginResult.fail("搜索关键词不能为空")
        if not self.settings.plugin_serpapi_api_key:
            return PluginResult.fail("SerpApi API Key 未配置，请在 .env 中设置 PLUGIN_SERPAPI_API_KEY")
        url = (
            f"{self.SERPAPI_ENDPOINT}?api_key={self.settings.plugin_serpapi_api_key}"
            f"&q={quote_plus(query)}&engine={quote_plus(self.search_engine)}"
            f"&num={self.max_results}&hl=zh-CN&gl=cn"
        )
        try:
            async with httpx.AsyncClient(timeout=self.timeout / 1000) as client:
                resp = await client.get(url)
            if resp.status_code != 200:
                return PluginResult.fail(f"SerpApi 调用失败，状态码: {resp.status_code}")
            body = resp.json()
            if body.get("error"):
                return PluginResult.fail(f"SerpApi 错误: {body.get('error')}")
            organic_results = body.get("organic_results") or []
            results: list[dict[str, str]] = []
            for item in organic_results[: self.max_results]:
                title = str(item.get("title") or "")
                link = str(item.get("link") or "")
                snippet = str(item.get("snippet") or "")
                if title and link:
                    results.append({"title": title, "url": link, "snippet": snippet})
            if not results:
                return PluginResult.fail("未找到相关搜索结果")
            content_lines = [f"以下是关于「{query}」的搜索结果：", ""]
            for idx, item in enumerate(results, start=1):
                content_lines.append(f"{idx}. {item['title']}")
                content_lines.append(item["snippet"])
                content_lines.append(f"来源: {item['url']}")
                content_lines.append("")
            return PluginResult.success_result(
                "\n".join(content_lines).strip(),
                {
                    "query": query,
                    "resultCount": len(results),
                    "searchEngine": self.search_engine,
                    "results": results,
                },
            )
        except Exception as exc:
            return PluginResult.fail(f"搜索失败: {exc}")
