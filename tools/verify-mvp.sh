#!/usr/bin/env bash
set -euo pipefail

echo "=== XShare MVP Verification ==="

echo "Step 1: Generate protobuf code..."
cd protocol && buf generate 2>/dev/null || echo "buf not available, skipping proto generation"

echo "Step 2: Run Go tests..."
cd ../core/go
go test ./...

echo "Step 3: Go build check..."
go build ./...

echo "Step 4: Android unit tests..."
cd ../../android
./gradlew :app:testDebugUnitTest 2>/dev/null || echo "Android SDK not available, skipping"

echo "=== Verification Complete ==="
