# XShare MVP Bringup Guide

## Prerequisites

- Go 1.22+
- ESP-IDF (for firmware)
- Android SDK (for Android app)
- buf CLI (for protobuf)

## Quick Start

### 1. Setup Development Environment

```bash
# Install Go with mise
mise use go@latest

# Install buf
go install github.com/bufbuild/buf/cmd/buf@latest
```

### 2. Run Tests

```bash
# Run all Go tests
cd core/go && go test ./...

# Run verification script
cd ../.. && bash tools/verify-mvp.sh
```

### 3. Build Components

#### Go Core
```bash
cd core/go
go build ./...
```

#### ESP32 Firmware
```bash
cd firmware/esp32
idf.py build
```

#### Android App
```bash
cd android
./gradlew assembleDebug
```

## Project Structure

```
XShare/
├── core/go/           # Go core (gVisor netstack, forwarding)
├── firmware/esp32/    # ESP32 firmware
├── android/          # Android app
├── protocol/         # Protobuf definitions
└── tools/            # Build/verification scripts
```

## Key Interfaces

### USB MUX Frame
- Protocol: custom binary framing
- Channels: control (1), data (2), ota (3)

### Control API
- `device.hello` - handshake
- `forward.start` / `forward.stop` - control forwarding
- `forward.get_stats` - get statistics
- `ota.begin/chunk/commit` - firmware updates

## Testing

- Unit tests: `go test ./...`
- Integration tests: `go test ./integration/...`
- View test coverage: `go test -cover ./...`
