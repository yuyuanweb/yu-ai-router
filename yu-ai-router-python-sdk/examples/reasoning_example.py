"""Reasoning mode example."""

from __future__ import annotations

import os

from yu_ai_router_sdk import ChatChunk, ChatMessage, ChatRequest, StreamCallback, YuAIClient


class ReasoningCallback(StreamCallback):
    def on_message(self, chunk: ChatChunk) -> None:
        if chunk.reasoning_content:
            print(f"\n💭 [思考] {chunk.reasoning_content}")
        if chunk.content:
            print(chunk.content, end="", flush=True)

    def on_complete(self) -> None:
        print("\n\n✅ 完成")

    def on_error(self, error: Exception) -> None:
        print(f"\n❌ 错误: {error}")


def main() -> None:
    api_key = os.getenv("YU_AI_API_KEY", "")
    base_url = os.getenv("YU_AI_BASE_URL", "http://localhost:8123/api")
    client = YuAIClient.builder().api_key(api_key).base_url(base_url).build()
    try:
        print("=== 深度思考模式示例 ===\n")
        request = ChatRequest(
            model="qwen-plus",
            messages=[ChatMessage.user("请详细解释量子纠缠现象，并给出实际应用")],
            enable_reasoning=True,
            temperature=0.7,
        )
        client.chat_stream(request, ReasoningCallback())
    finally:
        client.close()


if __name__ == "__main__":
    main()
