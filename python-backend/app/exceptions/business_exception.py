"""Business exception types."""

from __future__ import annotations


class BusinessException(Exception):
    def __init__(self, code: int, message: str) -> None:
        super().__init__(message)
        self.code = code
        self.message = message
