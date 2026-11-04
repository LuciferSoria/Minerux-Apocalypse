#!/usr/bin/env python3
"""Pre-deployment checks for a private Minerux NeoForge beta server.

Read-only: never modifies server.properties or prints secrets.
"""
import argparse
import os
from pathlib import Path
from urllib.parse import urlsplit


def inspect(properties: str, environment: dict[str, str]) -> list[str]:
    errors = []
    config = {}
    for raw in properties.splitlines():
        line = raw.strip()
        if not line or line.startswith(("#", "!")) or "=" not in line:
            continue
        key, value = line.split("=", 1)
        config[key.strip()] = value.strip()

    if config.get("online-mode", "").lower() != "true":
        errors.append("server.properties must set online-mode=true")
    if config.get("white-list", "").lower() != "true":
        errors.append("server.properties must set white-list=true until the beta gate is verified")
    if config.get("enforce-whitelist", "").lower() != "true":
        errors.append("server.properties must set enforce-whitelist=true until the beta gate is verified")

    url = environment.get("MINERUX_API_URL", "")
    try:
        parsed = urlsplit(url)
        if (parsed.scheme.lower() != "https" or not parsed.hostname
                or parsed.username or parsed.password or parsed.query or parsed.fragment
                or parsed.path not in ("", "/") or parsed.port == 0):
            errors.append("MINERUX_API_URL must be an HTTPS origin")
    except ValueError:
        errors.append("MINERUX_API_URL is invalid")

    token = environment.get("MINERUX_SERVER_TOKEN", "")
    if not (32 <= len(token) <= 4096 and all(c.isascii() and (c.isalnum() or c in "._~+/-") for c in token)):
        errors.append("MINERUX_SERVER_TOKEN is absent or invalid (value hidden)")
    return errors


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("server_properties", type=Path)
    args = parser.parse_args()
    try:
        contents = args.server_properties.read_text(encoding="utf-8")
    except OSError as error:
        print(f"FAIL: unable to read server.properties: {type(error).__name__}")
        return 2
    errors = inspect(contents, dict(os.environ))
    for error in errors:
        print("FAIL:", error)
    if errors:
        return 1
    print("PASS: basic beta server configuration checks (not a live authorization test)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
