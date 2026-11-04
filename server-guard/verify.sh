#!/usr/bin/env bash
# Standalone authorization adapter verification. No NeoForge runtime required.
set -euo pipefail
cd "$(dirname "$0")"
out="$(mktemp -d)"
trap 'rm -rf "$out"' EXIT
mapfile -d '' sources < <(find src/main/java src/test/java -type f -name '*.java' -print0 | sort -z)
if (( ${#sources[@]} == 0 )); then
  echo "No Java sources found" >&2
  exit 1
fi
javac -Xlint:all -Werror -d "$out" "${sources[@]}"
mapfile -d '' tests < <(find src/test/java -type f -name '*Test.java' -print0 | sort -z)
if (( ${#tests[@]} == 0 )); then
  echo "No Java tests found" >&2
  exit 1
fi
for source in "${tests[@]}"; do
  name="${source#src/test/java/}"
  name="${name%.java}"
  name="${name//\//.}"
  echo "Running $name"
  java -cp "$out" "$name"
done
set +e
java -cp "$out" dev.minerux.guard.ServerGuardPreflight > "$out/preflight.log" 2>&1
status=$?
set -e
if [[ "$status" != 3 ]]; then
  echo "Preflight returned unexpected exit code: $status" >&2
  exit 1
fi
if grep -E 'Bearer |https://' "$out/preflight.log"; then
  echo "Preflight leaked sensitive configuration" >&2
  exit 1
fi
echo "All standalone server guard checks passed"
