"""SDK client config."""

from __future__ import annotations

from dataclasses import dataclass


@dataclass(slots=True)
class ClientConfig:
    api_key: str
    base_url: str = "http://localhost:8123/api"
    connect_timeout: int = 10000
    read_timeout: int = 30000
    write_timeout: int = 30000
    max_retries: int = 3
    retry_delay: int = 1000

    def validate(self) -> None:
        if not self.api_key or not self.api_key.strip():
            msg = "api_key is required"
            raise ValueError(msg)
        if not self.base_url or not self.base_url.strip():
            msg = "base_url is required"
            raise ValueError(msg)
