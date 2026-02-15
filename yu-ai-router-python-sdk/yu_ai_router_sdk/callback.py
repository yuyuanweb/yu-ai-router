"""SDK stream callback protocol."""

from __future__ import annotations

from typing import Protocol

from yu_ai_router_sdk.models import ChatChunk


class StreamCallback(Protocol):
    """Callback for streaming responses."""

    def on_message(self, chunk: ChatChunk) -> None:
        """Called when a stream chunk arrives."""

    def on_complete(self) -> None:
        """Called when stream completes."""

    def on_error(self, error: Exception) -> None:
        """Called when stream fails."""
