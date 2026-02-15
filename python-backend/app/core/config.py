"""Application settings."""

from __future__ import annotations

from functools import lru_cache

from dotenv import load_dotenv
from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict

load_dotenv()


class Settings(BaseSettings):
    model_config = SettingsConfigDict(extra="ignore", case_sensitive=False)

    app_name: str = "yu-ai-router-python-backend"
    app_env: str = Field(default="local")
    app_host: str = "0.0.0.0"
    app_port: int = 8123
    app_base_path: str = "/api"
    app_timezone: str = "Asia/Shanghai"

    mysql_host: str = "127.0.0.1"
    mysql_port: int = 3306
    mysql_user: str = "root"
    mysql_password: str = ""
    mysql_db: str = "yu_ai_router"
    mysql_charset: str = "utf8mb4"

    redis_host: str = "127.0.0.1"
    redis_port: int = 6379
    redis_db: int = 0
    redis_password: str = ""

    cors_allow_origin_patterns: str = "*"
    log_level: str = "INFO"
    ai_base_url: str = "https://dashscope.aliyuncs.com/compatible-mode"
    ai_api_key: str = ""
    ai_model: str = "qwen-plus"
    ai_chat_completions_path: str = "/v1/chat/completions"
    ai_timeout_seconds: int = 120
    ai_cache_enabled: bool = True
    ai_cache_ttl: int = 3600

    stripe_api_key: str = ""
    stripe_webhook_secret: str = ""
    stripe_success_url: str = "http://localhost:5173/recharge/success"
    stripe_cancel_url: str = "http://localhost:5173/recharge/cancel"
    encryption_secret_key: str = "yupi-ai-router-secret-key-256"
    plugin_serpapi_api_key: str = ""

    @property
    def mysql_dsn(self) -> str:
        return (
            "mysql+aiomysql://"
            f"{self.mysql_user}:{self.mysql_password}"
            f"@{self.mysql_host}:{self.mysql_port}/{self.mysql_db}"
            f"?charset={self.mysql_charset}"
        )

    @property
    def redis_dsn(self) -> str:
        auth = f":{self.redis_password}@" if self.redis_password else ""
        return f"redis://{auth}{self.redis_host}:{self.redis_port}/{self.redis_db}"

    @property
    def ai_chat_completions_url(self) -> str:
        return f"{self.ai_base_url.rstrip('/')}{self.ai_chat_completions_path}"


@lru_cache
def get_settings() -> Settings:
    return Settings()
