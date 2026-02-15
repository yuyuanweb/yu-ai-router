"""SDK data models."""

from __future__ import annotations

from dataclasses import dataclass, field
from typing import Any


@dataclass(slots=True)
class ChatMessage:
    role: str
    content: str

    @classmethod
    def system(cls, content: str) -> "ChatMessage":
        return cls(role="system", content=content)

    @classmethod
    def user(cls, content: str) -> "ChatMessage":
        return cls(role="user", content=content)

    @classmethod
    def assistant(cls, content: str) -> "ChatMessage":
        return cls(role="assistant", content=content)


@dataclass(slots=True)
class ChatRequest:
    model: str | None = None
    messages: list[ChatMessage] = field(default_factory=list)
    stream: bool = False
    temperature: float | None = None
    max_tokens: int | None = None
    enable_reasoning: bool | None = None
    routing_strategy: str | None = None

    @classmethod
    def simple(cls, user_message: str) -> "ChatRequest":
        return cls(messages=[ChatMessage.user(user_message)])

    @classmethod
    def with_model(cls, model: str, user_message: str) -> "ChatRequest":
        return cls(model=model, messages=[ChatMessage.user(user_message)])

    def to_payload(self) -> dict[str, Any]:
        payload: dict[str, Any] = {
            "model": self.model,
            "messages": [{"role": item.role, "content": item.content} for item in self.messages],
            "stream": self.stream,
            "temperature": self.temperature,
            "max_tokens": self.max_tokens,
            "enable_reasoning": self.enable_reasoning,
            "routing_strategy": self.routing_strategy,
        }
        return {k: v for k, v in payload.items() if v is not None}


@dataclass(slots=True)
class ChatResponseMessage:
    role: str | None
    content: str | None


@dataclass(slots=True)
class ChatResponseChoice:
    index: int
    message: ChatResponseMessage
    finish_reason: str | None


@dataclass(slots=True)
class ChatResponseUsage:
    prompt_tokens: int
    completion_tokens: int
    total_tokens: int


@dataclass(slots=True)
class ChatResponse:
    id: str
    object: str
    created: int
    model: str
    choices: list[ChatResponseChoice]
    usage: ChatResponseUsage

    @property
    def content(self) -> str | None:
        if not self.choices:
            return None
        return self.choices[0].message.content

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> "ChatResponse":
        choices = []
        for item in data.get("choices", []) or []:
            message = item.get("message") or {}
            choices.append(
                ChatResponseChoice(
                    index=int(item.get("index", 0)),
                    message=ChatResponseMessage(
                        role=message.get("role"),
                        content=message.get("content"),
                    ),
                    finish_reason=item.get("finishReason") or item.get("finish_reason"),
                )
            )
        usage_data = data.get("usage") or {}
        usage = ChatResponseUsage(
            prompt_tokens=int(usage_data.get("promptTokens", usage_data.get("prompt_tokens", 0)) or 0),
            completion_tokens=int(
                usage_data.get("completionTokens", usage_data.get("completion_tokens", 0)) or 0
            ),
            total_tokens=int(usage_data.get("totalTokens", usage_data.get("total_tokens", 0)) or 0),
        )
        return cls(
            id=str(data.get("id", "")),
            object=str(data.get("object", "chat.completion")),
            created=int(data.get("created", 0) or 0),
            model=str(data.get("model", "")),
            choices=choices,
            usage=usage,
        )


@dataclass(slots=True)
class StreamDelta:
    role: str | None
    content: str | None
    reasoning_content: str | None


@dataclass(slots=True)
class StreamChoice:
    index: int
    delta: StreamDelta
    finish_reason: str | None


@dataclass(slots=True)
class StreamResponse:
    id: str | None
    object: str | None
    created: int | None
    model: str | None
    choices: list[StreamChoice]

    @classmethod
    def from_dict(cls, data: dict[str, Any]) -> "StreamResponse":
        parsed_choices: list[StreamChoice] = []
        for item in data.get("choices", []) or []:
            delta_data = item.get("delta") or {}
            parsed_choices.append(
                StreamChoice(
                    index=int(item.get("index", 0)),
                    delta=StreamDelta(
                        role=delta_data.get("role"),
                        content=delta_data.get("content"),
                        reasoning_content=delta_data.get("reasoningContent")
                        or delta_data.get("reasoning_content"),
                    ),
                    finish_reason=item.get("finishReason") or item.get("finish_reason"),
                )
            )
        created_raw = data.get("created")
        created_val = int(created_raw) if created_raw is not None else None
        return cls(
            id=data.get("id"),
            object=data.get("object"),
            created=created_val,
            model=data.get("model"),
            choices=parsed_choices,
        )


@dataclass(slots=True)
class ChatChunk:
    content: str | None = None
    reasoning_content: str | None = None
    done: bool = False
    model: str | None = None
    prompt_tokens: int | None = None
    completion_tokens: int | None = None
