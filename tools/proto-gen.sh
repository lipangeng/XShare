#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/../protocol"

if ! awk '
  /^plugins:[[:space:]]*\[[[:space:]]*\]/ { next }
  /^plugins:[[:space:]]*\[[^]]+\]/ { found=1; next }
  /^plugins:[[:space:]]*$/ { in_plugins=1; next }
  in_plugins && /^[^[:space:]-]/ { in_plugins=0 }
  in_plugins && /^[[:space:]]*-[[:space:]]/ { found=1 }
  END { exit(found ? 0 : 1) }
' buf.gen.yaml; then
  echo "error: protocol/buf.gen.yaml has no plugins configured. Add at least one plugin entry under 'plugins:'." >&2
  exit 1
fi

buf generate
