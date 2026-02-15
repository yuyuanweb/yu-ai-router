"""Advanced SDK usage example."""

from __future__ import annotations

import os

from yu_ai_router_sdk import ChatMessage, ChatRequest, YuAIClient


def main() -> None:
    api_key = os.getenv("YU_AI_API_KEY", "")
    base_url = os.getenv("YU_AI_BASE_URL", "http://localhost:8123/api")
    client = (
        YuAIClient.builder()
        .api_key(api_key)
        .base_url(base_url)
        .connect_timeout(15000)
        .read_timeout(60000)
        .max_retries(5)
        .build()
    )
    try:
        print("=== 多轮对话示例 ===\n")
        request = ChatRequest(
            messages=[
                ChatMessage.system("你是一个编程助手"),
                ChatMessage.user("什么是 Java？"),
                ChatMessage.assistant("Java 是一种面向对象的编程语言..."),
                ChatMessage.user("它的主要特点是什么？"),
            ],
            model="qwen-turbo",
            temperature=0.7,
        )
        response = client.chat(request)
        print(f"回答: {response.content}")
    finally:
        client.close()


if __name__ == "__main__":
    main()
