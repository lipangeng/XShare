#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

fail() {
  echo "error: $*" >&2
  exit 1
}

require_cmd() {
  local cmd="$1"
  command -v "$cmd" >/dev/null 2>&1 || fail "required command '$cmd' is not installed or not on PATH"
}

echo "==> checking prerequisites"
require_cmd buf
require_cmd go
require_cmd java

GRADLEW="${ROOT_DIR}/android/gradlew"
[[ -f "${GRADLEW}" ]] || fail "missing Android Gradle wrapper at android/gradlew"
[[ -x "${GRADLEW}" ]] || fail "android/gradlew is not executable (run: chmod +x android/gradlew)"

echo "==> generating protocol code"
(
  cd "${ROOT_DIR}/protocol"
  buf generate
)

echo "==> running Go tests"
(
  cd "${ROOT_DIR}/core/go"
  go test ./...
)

echo "==> running Android unit tests"
(
  cd "${ROOT_DIR}/android"
  ./gradlew :app:testDebugUnitTest
)

echo "MVP verification complete."
