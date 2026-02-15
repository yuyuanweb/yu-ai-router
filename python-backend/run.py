"""Unified startup entry for PyCharm / CLI / Docker."""

from __future__ import annotations

import argparse
from pathlib import Path

import uvicorn
from dotenv import load_dotenv

PROJECT_ROOT = Path(__file__).resolve().parent
ENV_FILE_MAP: dict[str, str] = {
    "dev": ".env",
    "local": ".env.local",
    "prod": ".env.prod",
}


def _parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Start yu-ai-router python backend")
    parser.add_argument("--env", choices=tuple(ENV_FILE_MAP.keys()), default="local")
    parser.add_argument("--host", default="0.0.0.0")
    parser.add_argument("--port", type=int, default=8123)
    parser.add_argument("--reload", action="store_true")
    return parser.parse_args()


def _load_env(env_name: str) -> None:
    env_filename = ENV_FILE_MAP.get(env_name, ".env")
    env_path = PROJECT_ROOT / env_filename
    default_env_path = PROJECT_ROOT / ".env"
    if default_env_path.exists():
        load_dotenv(default_env_path, override=False)
    if env_path.exists():
        load_dotenv(env_path, override=True)


def main() -> None:
    args = _parse_args()
    _load_env(args.env)
    uvicorn.run(
        "app.main:app",
        host=args.host,
        port=args.port,
        reload=args.reload,
        factory=False,
    )


if __name__ == "__main__":
    main()
