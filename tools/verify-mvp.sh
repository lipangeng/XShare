#!/usr/bin/env bash
set -euo pipefail

echo "=== XShare MVP Verification ==="
echo ""

# Step 1: Protobuf
echo "Step 1: Generate protobuf code..."
cd protocol
if command -v buf &> /dev/null; then
    buf generate && echo "  ✓ Protobuf code generated"
else
    echo "  ⊘ buf not available, skipping proto generation"
fi
cd ..
echo ""

# Step 2: Go tests
echo "Step 2: Run Go tests..."
cd core/go
if command -v go &> /dev/null; then
    go test ./... -count=1 && echo "  ✓ All Go tests passed"
else
    echo "  ⊘ Go not available, skipping tests"
fi
cd ../..
echo ""

# Step 3: Go build
echo "Step 3: Go build check..."
cd core/go
if command -v go &> /dev/null; then
    go build ./... && echo "  ✓ Go build successful"
else
    echo "  ⊘ Go not available, skipping build"
fi
cd ../..
echo ""

# Step 4: ESP32 firmware
echo "Step 4: ESP32 firmware build..."
cd firmware/esp32
if command -v idf.py &> /dev/null; then
    idf.py build && echo "  ✓ ESP32 firmware build successful"
else
    echo "  ⊘ ESP-IDF not available, skipping firmware build"
fi
cd ../..
echo ""

# Step 5: Android
echo "Step 5: Android build and tests..."
cd android
if command -v gradle &> /dev/null || [ -f "./gradlew" ]; then
    ./gradlew :app:testDebugUnitTest :corebridge:build && echo "  ✓ Android build and tests successful"
else
    echo "  ⊘ Gradle not available, skipping Android build"
fi
cd ..
echo ""

echo "=== Verification Complete ==="
