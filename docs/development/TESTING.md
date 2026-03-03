# XShare Testing Guide

## Test Status

| Component | Test Type | Status | Notes |
|-----------|-----------|--------|-------|
| Go Core | Unit Tests | ✅ Pass | 9 packages, all tests pass |
| Go Core | Integration Tests | ✅ Pass | E2E loopback tests pass |
| Go Core | Build | ✅ Pass | `go build ./...` succeeds |
| Android App | Unit Tests | ⏸️ Skipped | Requires Android SDK |
| Android App | Build | ⏸️ Skipped | Requires Gradle/Android SDK |
| ESP32 Firmware | Build | ⏸️ Skipped | Requires ESP-IDF |

## Running Tests

### Go Core (Available Now)

```bash
# Run all tests
cd core/go && go test ./... -v

# Run with coverage
go test ./... -cover

# Run specific package
go test ./pkg/controller -v
go test ./integration -v
```

### Android (Requires Setup)

```bash
# Prerequisites:
# - Android SDK installed
# - ANDROID_HOME environment variable set
# - Gradle or ./gradlew wrapper

cd android
./gradlew :app:testDebugUnitTest
./gradlew :corebridge:build
./gradlew assembleDebug
```

### ESP32 Firmware (Requires Setup)

```bash
# Prerequisites:
# - ESP-IDF v5.x installed
# - idf.py in PATH

cd firmware/esp32
idf.py build
idf.py flash monitor
```

## Verification Script

Run the full verification suite:

```bash
bash tools/verify-mvp.sh
```

This script:
1. Generates protobuf code (if buf available)
2. Runs Go tests
3. Builds Go core
4. Builds ESP32 firmware (if ESP-IDF available)
5. Runs Android tests (if Gradle available)

## Test Coverage

### Go Core Packages

| Package | Tests | Coverage |
|---------|-------|----------|
| `cmd/xshare-daemon` | 1 | Version string validation |
| `integration` | 2 | E2E loopback, start/stop |
| `pkg/api` | 2 | Protocol constants |
| `pkg/controller` | 4 | State machine, stats |
| `pkg/dataplane/netstack` | 2 | Engine lifecycle |
| `pkg/diag` | 3 | Stats counters |
| `pkg/forwarder` | 3 | Session management |
| `pkg/protocol/mux` | 4 | Frame codec |
| `pkg/transport/usb` | 1 | Link read/write |

### Android Tests

| Class | Tests | Coverage |
|-------|-------|----------|
| `ForwardViewModelTest` | 3 | State transitions |

### ESP32 Tests

| Component | Tests | Coverage |
|-----------|-------|----------|
| `softap_mgr` | Stub | Returns ESP_OK |
| `packet_io` | Stub | Returns ESP_OK |
| `usb_mux` | Stub | Returns ESP_OK |
| `ctrl_agent` | Stub | Returns ESP_OK |

## CI/CD Integration

Example GitHub Actions workflow:

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: jdx/mise-action@v2
      
      - name: Run Go tests
        run: cd core/go && go test ./...
      
      - name: Build Go core
        run: cd core/go && go build ./...
      
      - name: Run verification script
        run: bash tools/verify-mvp.sh
```

## Troubleshooting

### Go Tests Fail

```bash
# Clean build cache
go clean -testcache

# Re-run tests
go test ./... -count=1
```

### Android Build Fails

```bash
# Clean Gradle cache
./gradlew clean

# Check SDK installation
echo $ANDROID_HOME
adb devices
```

### ESP32 Build Fails

```bash
# Check ESP-IDF installation
idf.py --version

# Clean build
idf.py fullclean
idf.py build
```
