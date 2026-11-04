#!/usr/bin/env python3
"""Fail CI if a development NeoForge JAR is missing required guard classes.

This checks packaging only. It does not prove login denial or mod compatibility.
"""
import sys
import zipfile
from pathlib import Path

REQUIRED = {
    "META-INF/neoforge.mods.toml",
    "dev/minerux/guard/neoforge/MineruxBetaGuardMod.class",
    "dev/minerux/guard/AuthorizationService.class",
    "dev/minerux/guard/LoginAccessGate.class",
    "dev/minerux/guard/BetaAccessClient.class",
    "dev/minerux/guard/LoginCheckWorkers.class",
    "dev/minerux/guard/ServerGuardConfig.class",
}
FORBIDDEN_NAMES = {"state.json", "codes.json", ".env", "secrets.json"}

def check(path: Path) -> None:
    with zipfile.ZipFile(path) as jar:
        names = set(jar.namelist())
        missing = REQUIRED - names
        if missing:
            raise ValueError("Missing required JAR entries: " + ", ".join(sorted(missing)))
        for name in names:
            if Path(name).name.lower() in FORBIDDEN_NAMES:
                raise ValueError("Private file accidentally packaged: " + name)
        metadata = jar.read("META-INF/neoforge.mods.toml").decode("utf-8")
        if 'modId="mineruxbetaguard"' not in metadata:
            raise ValueError("Incorrect NeoForge mod ID")
        if len(names) != len(jar.namelist()):
            raise ValueError("Duplicate JAR entries")
    print("JAR structure verified:", path.name)
    print("Runtime login enforcement still requires real server testing.")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        sys.exit("Usage: verify_jar.py <jar> [<jar>...]")
    for arg in sys.argv[1:]:
        check(Path(arg))
