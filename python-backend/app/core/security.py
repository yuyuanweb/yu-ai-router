"""Security helpers."""

from __future__ import annotations

import hashlib

from app.core.constants import PASSWORD_SALT


def encrypt_password(password: str) -> str:
    return hashlib.md5(f"{password}{PASSWORD_SALT}".encode("utf-8")).hexdigest()
