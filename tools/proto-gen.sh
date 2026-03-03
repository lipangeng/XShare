#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

cd protocol && buf generate

echo "Protobuf code generation complete"
