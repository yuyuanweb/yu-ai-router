"""SDK exceptions."""

from __future__ import annotations


class YuAIException(Exception):
    """Base exception for SDK errors."""

    def __init__(self, message: str, status_code: int | None = None) -> None:
        super().__init__(message)
        self.status_code = status_code


class AuthException(YuAIException):
    """Raised when auth fails."""


class RateLimitException(YuAIException):
    """Raised when request is rate limited."""
