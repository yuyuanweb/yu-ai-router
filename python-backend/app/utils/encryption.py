"""AES encryption helpers used by BYOK."""

from __future__ import annotations

import base64

from app.core.config import get_settings


def _get_key_bytes(secret_key: str) -> bytes:
    src = secret_key.encode("utf-8")
    return src[:32].ljust(32, b"\0")


def encrypt_text(plain_text: str) -> str:
    try:
        from Crypto.Cipher import AES
        from Crypto.Util.Padding import pad
    except ImportError as exc:  # pragma: no cover
        raise RuntimeError('缺少依赖 "pycryptodome"，无法加密 BYOK 密钥') from exc

    key = _get_key_bytes(get_settings().encryption_secret_key)
    cipher = AES.new(key, AES.MODE_ECB)
    encrypted = cipher.encrypt(pad(plain_text.encode("utf-8"), AES.block_size))
    return base64.b64encode(encrypted).decode("utf-8")


def decrypt_text(cipher_text: str) -> str:
    try:
        from Crypto.Cipher import AES
        from Crypto.Util.Padding import unpad
    except ImportError as exc:  # pragma: no cover
        raise RuntimeError('缺少依赖 "pycryptodome"，无法解密 BYOK 密钥') from exc

    key = _get_key_bytes(get_settings().encryption_secret_key)
    cipher = AES.new(key, AES.MODE_ECB)
    raw = base64.b64decode(cipher_text)
    plain = unpad(cipher.decrypt(raw), AES.block_size)
    return plain.decode("utf-8")
