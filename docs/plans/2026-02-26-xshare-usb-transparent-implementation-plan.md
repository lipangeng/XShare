# XShare USB Transparent Gateway Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build an MVP that transparently forwards terminal TCP/UDP traffic from ESP32 SoftAP through Android (non-root) over USB custom protocol, with Go core powered by gVisor netstack.

**Architecture:** ESP32 handles Wi-Fi access (SoftAP/DHCP) and packet ingress/egress only. Android app owns lifecycle/control UI and delegates forwarding to Go core through JNI. Go core multiplexes USB channels (control/data/ota), injects IP packets into gVisor netstack, bridges sessions to outbound sockets, and sends return packets back to ESP32.

**Tech Stack:** ESP-IDF (C), Go 1.22+, gVisor netstack, Protobuf + buf, Android (Kotlin + JNI/NDK), Gradle.

---

### Task 1: Initialize Monorepo Skeleton

**Files:**
- Create: `firmware/esp32/main/CMakeLists.txt`
- Create: `firmware/esp32/main/main.c`
- Create: `firmware/esp32/CMakeLists.txt`
- Create: `firmware/esp32/sdkconfig.defaults`
- Create: `core/go/go.mod`
- Create: `core/go/cmd/xshare-daemon/main.go`
- Create: `android/settings.gradle.kts`
- Create: `android/build.gradle.kts`
- Create: `android/app/build.gradle.kts`
- Create: `protocol/buf.yaml`
- Create: `protocol/buf.gen.yaml`
- Create: `tools/proto-gen.sh`

**Step 1: Write the failing test**

Create `core/go/cmd/xshare-daemon/main_test.go`:
```go
package main

import "testing"

func TestVersionStringNonEmpty(t *testing.T) {
	if version == "" {
		t.Fatal("version must not be empty")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `cd core/go && go test ./cmd/xshare-daemon -v`
Expected: FAIL with `undefined: version`.

**Step 3: Write minimal implementation**

In `core/go/cmd/xshare-daemon/main.go` add:
```go
package main

import "fmt"

var version = "dev"

func main() {
	fmt.Println("xshare-daemon", version)
}
```

**Step 4: Run test to verify it passes**

Run: `cd core/go && go test ./cmd/xshare-daemon -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add firmware core android protocol tools
git commit -m "chore: scaffold xshare monorepo skeleton"
```

### Task 2: Define Protobuf Contracts (Control/Data/OTA)

**Files:**
- Create: `protocol/proto/control/v1/control.proto`
- Create: `protocol/proto/data/v1/data.proto`
- Create: `protocol/proto/ota/v1/ota.proto`
- Modify: `protocol/buf.yaml`
- Modify: `protocol/buf.gen.yaml`

**Step 1: Write the failing test**

Create `core/go/pkg/api/proto_contract_test.go`:
```go
package api

import "testing"

func TestControlMethodForwardStartExists(t *testing.T) {
	if MethodForwardStart == "" {
		t.Fatal("MethodForwardStart must be defined from generated protobuf")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `cd core/go && go test ./pkg/api -v`
Expected: FAIL with `undefined: MethodForwardStart`.

**Step 3: Write minimal implementation**

- Define protobuf enums/messages with required fields: `version`, `request_id`, `correlation_id`, error model.
- Add a generated constants shim in `core/go/pkg/api/constants.go`:
```go
package api

const MethodForwardStart = "forward.start"
```
- Generate code: `cd protocol && buf generate`.

**Step 4: Run test to verify it passes**

Run: `cd core/go && go test ./pkg/api -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add protocol core/go/pkg/api
git commit -m "feat(protocol): add protobuf contracts for control data ota"
```

### Task 3: Implement USB MUX Frame Codec in Go

**Files:**
- Create: `core/go/pkg/protocol/mux/frame.go`
- Create: `core/go/pkg/protocol/mux/codec.go`
- Create: `core/go/pkg/protocol/mux/codec_test.go`

**Step 1: Write the failing test**

In `codec_test.go`:
```go
func TestEncodeDecodeRoundTrip(t *testing.T) {
	in := &Frame{Version: 1, Channel: 2, StreamID: 7, Seq: 1, Payload: []byte{1,2,3}}
	b, err := Encode(in)
	if err != nil { t.Fatal(err) }
	out, err := Decode(b)
	if err != nil { t.Fatal(err) }
	if out.Channel != in.Channel || len(out.Payload) != 3 { t.Fatal("roundtrip mismatch") }
}
```

**Step 2: Run test to verify it fails**

Run: `cd core/go && go test ./pkg/protocol/mux -v`
Expected: FAIL with missing `Encode/Decode`.

**Step 3: Write minimal implementation**

Implement frame struct + binary codec with CRC32 verification.

**Step 4: Run test to verify it passes**

Run: `cd core/go && go test ./pkg/protocol/mux -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add core/go/pkg/protocol/mux
git commit -m "feat(core): add usb mux frame codec"
```

### Task 4: Build Go USB Transport Abstraction

**Files:**
- Create: `core/go/pkg/transport/usb/link.go`
- Create: `core/go/pkg/transport/usb/link_test.go`
- Create: `core/go/pkg/transport/usb/mockio.go`

**Step 1: Write the failing test**

`link_test.go`:
```go
func TestLinkReadWriteFrame(t *testing.T) {
	io := newMockIO()
	l := NewLink(io)
	f := &mux.Frame{Version:1, Channel:2, Payload:[]byte("abc")}
	if err := l.WriteFrame(context.Background(), f); err != nil { t.Fatal(err) }
	got, err := l.ReadFrame(context.Background())
	if err != nil { t.Fatal(err) }
	if string(got.Payload) != "abc" { t.Fatal("payload mismatch") }
}
```

**Step 2: Run test to verify it fails**

Run: `cd core/go && go test ./pkg/transport/usb -v`
Expected: FAIL with undefined `NewLink`.

**Step 3: Write minimal implementation**

Implement `Link` interface and a blocking read/write loop using mux codec.

**Step 4: Run test to verify it passes**

Run: `cd core/go && go test ./pkg/transport/usb -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add core/go/pkg/transport/usb
git commit -m "feat(core): add usb transport abstraction"
```

### Task 5: Integrate gVisor Netstack Engine Skeleton

**Files:**
- Create: `core/go/pkg/dataplane/netstack/engine.go`
- Create: `core/go/pkg/dataplane/netstack/engine_test.go`

**Step 1: Write the failing test**

`engine_test.go`:
```go
func TestEngineInjectAndRead(t *testing.T) {
	e := NewEngine()
	if err := e.Start(); err != nil { t.Fatal(err) }
	defer e.Stop()
	if err := e.InjectInbound([]byte{0x45,0,0,20}); err == nil {
		t.Fatal("expected short packet error")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `cd core/go && go test ./pkg/dataplane/netstack -v`
Expected: FAIL with undefined `NewEngine`.

**Step 3: Write minimal implementation**

Implement engine scaffold and validation; wire placeholder channels for packet path.

**Step 4: Run test to verify it passes**

Run: `cd core/go && go test ./pkg/dataplane/netstack -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add core/go/pkg/dataplane/netstack
git commit -m "feat(core): add netstack engine scaffold"
```

### Task 6: Implement Forwarder (TCP/UDP session bridge)

**Files:**
- Create: `core/go/pkg/forwarder/forwarder.go`
- Create: `core/go/pkg/forwarder/forwarder_test.go`

**Step 1: Write the failing test**

`forwarder_test.go`:
```go
func TestForwarderCreateUDPMapping(t *testing.T) {
	f := New()
	id := FiveTuple{SrcIP:"1.1.1.1", SrcPort:1234, DstIP:"8.8.8.8", DstPort:53, Proto:17}
	_, created := f.GetOrCreate(id)
	if !created { t.Fatal("expected new mapping") }
}
```

**Step 2: Run test to verify it fails**

Run: `cd core/go && go test ./pkg/forwarder -v`
Expected: FAIL with undefined `New`.

**Step 3: Write minimal implementation**

Implement in-memory session table with timeout metadata; no optimization yet.

**Step 4: Run test to verify it passes**

Run: `cd core/go && go test ./pkg/forwarder -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add core/go/pkg/forwarder
git commit -m "feat(core): add tcp udp session forwarder table"
```

### Task 7: Add Controller State Machine and Control APIs

**Files:**
- Create: `core/go/pkg/controller/controller.go`
- Create: `core/go/pkg/controller/controller_test.go`

**Step 1: Write the failing test**

`controller_test.go`:
```go
func TestStartStopStateTransition(t *testing.T) {
	c := NewController()
	if err := c.StartForward(); err != nil { t.Fatal(err) }
	if err := c.StopForward(); err != nil { t.Fatal(err) }
}
```

**Step 2: Run test to verify it fails**

Run: `cd core/go && go test ./pkg/controller -v`
Expected: FAIL with undefined `NewController`.

**Step 3: Write minimal implementation**

Implement state machine: `Idle -> Running -> Idle`, with guard checks and error codes.

**Step 4: Run test to verify it passes**

Run: `cd core/go && go test ./pkg/controller -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add core/go/pkg/controller
git commit -m "feat(core): add forwarding controller state machine"
```

### Task 8: ESP32 SoftAP + Packet IO Stubs

**Files:**
- Create: `firmware/esp32/components/softap_mgr/include/softap_mgr.h`
- Create: `firmware/esp32/components/softap_mgr/softap_mgr.c`
- Create: `firmware/esp32/components/packet_io/include/packet_io.h`
- Create: `firmware/esp32/components/packet_io/packet_io.c`
- Modify: `firmware/esp32/main/main.c`

**Step 1: Write the failing test**

Create `firmware/esp32/main/test_smoke.c`:
```c
#include "softap_mgr.h"

void test_softap_start_returns_ok(void) {
    TEST_ASSERT_EQUAL(ESP_OK, softap_mgr_start());
}
```

**Step 2: Run test to verify it fails**

Run: `cd firmware/esp32 && idf.py build`
Expected: FAIL with missing symbol/component wiring.

**Step 3: Write minimal implementation**

Add stub components returning `ESP_OK`; wire component registration in CMake.

**Step 4: Run test to verify it passes**

Run: `cd firmware/esp32 && idf.py build`
Expected: PASS.

**Step 5: Commit**

```bash
git add firmware/esp32
git commit -m "feat(firmware): add softap and packet io component stubs"
```

### Task 9: Android App Skeleton + JNI Bridge Contract

**Files:**
- Create: `android/app/src/main/java/com/xshare/app/MainActivity.kt`
- Create: `android/app/src/main/java/com/xshare/app/ForwardViewModel.kt`
- Create: `android/corebridge/src/main/cpp/native_bridge.cpp`
- Create: `android/corebridge/src/main/java/com/xshare/corebridge/CoreBridge.kt`

**Step 1: Write the failing test**

Create `android/app/src/test/java/com/xshare/app/ForwardViewModelTest.kt`:
```kotlin
@Test
fun startForward_updatesStateToRunning() {
    val vm = ForwardViewModel(FakeBridge())
    vm.startForward()
    assertEquals(State.Running, vm.state.value)
}
```

**Step 2: Run test to verify it fails**

Run: `cd android && ./gradlew :app:testDebugUnitTest`
Expected: FAIL due to missing ViewModel/bridge implementation.

**Step 3: Write minimal implementation**

Implement ViewModel state transition and JNI bridge interface with placeholder native calls.

**Step 4: Run test to verify it passes**

Run: `cd android && ./gradlew :app:testDebugUnitTest`
Expected: PASS.

**Step 5: Commit**

```bash
git add android
git commit -m "feat(android): add app skeleton and core jni bridge"
```

### Task 10: End-to-End Loopback Integration Test (Host-side)

**Files:**
- Create: `core/go/integration/loopback_e2e_test.go`
- Create: `core/go/integration/testdata/sample_ipv4_udp.bin`

**Step 1: Write the failing test**

`loopback_e2e_test.go`:
```go
func TestE2E_UplinkToDownlinkLoopback(t *testing.T) {
	ctx := context.Background()
	s := NewHarness(t)
	defer s.Close()
	pkt := loadPacket(t, "testdata/sample_ipv4_udp.bin")
	if err := s.InjectFromEsp(pkt); err != nil { t.Fatal(err) }
	resp, err := s.ReadToEsp(ctx)
	if err != nil { t.Fatal(err) }
	if len(resp) == 0 { t.Fatal("expected response packet") }
}
```

**Step 2: Run test to verify it fails**

Run: `cd core/go && go test ./integration -v`
Expected: FAIL with missing harness.

**Step 3: Write minimal implementation**

Implement harness with mock USB transport and fake outbound responder.

**Step 4: Run test to verify it passes**

Run: `cd core/go && go test ./integration -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add core/go/integration
git commit -m "test(core): add e2e loopback integration harness"
```

### Task 11: Diagnostics and Runtime Stats API

**Files:**
- Create: `core/go/pkg/diag/stats.go`
- Modify: `core/go/pkg/controller/controller.go`
- Create: `core/go/pkg/diag/stats_test.go`

**Step 1: Write the failing test**

`stats_test.go`:
```go
func TestStatsCounterIncrement(t *testing.T) {
	s := NewStats()
	s.IncUplinkPackets(3)
	if s.Snapshot().UplinkPackets != 3 { t.Fatal("bad uplink counter") }
}
```

**Step 2: Run test to verify it fails**

Run: `cd core/go && go test ./pkg/diag -v`
Expected: FAIL with undefined `NewStats`.

**Step 3: Write minimal implementation**

Implement atomic counters and controller hook for `forward.get_stats`.

**Step 4: Run test to verify it passes**

Run: `cd core/go && go test ./pkg/diag -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add core/go/pkg/diag core/go/pkg/controller
git commit -m "feat(core): add runtime stats and diagnostics api"
```

### Task 12: Verification Gate and Developer Docs

**Files:**
- Create: `docs/protocol/README.md`
- Create: `docs/development/mvp-bringup.md`
- Create: `tools/verify-mvp.sh`

**Step 1: Write the failing test**

Create `core/go/pkg/controller/verify_contract_test.go`:
```go
func TestForwardStartThenStatsNonZeroContract(t *testing.T) {
	c := NewController()
	_ = c.StartForward()
	stats := c.Stats()
	if stats == nil { t.Fatal("stats should be available") }
}
```

**Step 2: Run test to verify it fails**

Run: `cd core/go && go test ./pkg/controller -v`
Expected: FAIL if `Stats()` missing.

**Step 3: Write minimal implementation**

Add missing controller `Stats()` access, write verification script:
```bash
#!/usr/bin/env bash
set -euo pipefail
cd protocol && buf generate
cd ../core/go && go test ./...
cd ../../android && ./gradlew :app:testDebugUnitTest
```

**Step 4: Run test to verify it passes**

Run:
- `cd core/go && go test ./pkg/controller -v`
- `bash tools/verify-mvp.sh`

Expected: PASS on contract test and verification script.

**Step 5: Commit**

```bash
git add docs tools core/go/pkg/controller
git commit -m "docs: add mvp bringup and verification workflow"
```

## Cross-cutting execution rules
- Use `@superpowers:test-driven-development` before each implementation change.
- Use `@superpowers:systematic-debugging` immediately for any test failure not explained by current step.
- Use `@superpowers:verification-before-completion` before claiming milestone complete.
- Keep commits small and aligned to one task each.
- Do not add features beyond approved design (YAGNI).
