# XShare Android USB Transparent Gateway Design

Date: 2026-02-26
Status: Approved
Scope: MVP design for ESP32-S3 + Android (non-root), USB custom protocol, transparent TCP/UDP forwarding

## 1. End-to-end architecture and data flows

MVP focus: make terminal devices access the Internet transparently by connecting to ESP32 SoftAP, while Android app handles forwarding and outbound traffic goes through system network path (including an already-enabled VPN if present on phone).

```text
Terminal Device
   | Wi-Fi (SoftAP)
   v
ESP32-S3 Gateway
   |-- Control Plane: USB ctrl channel <--------------------> Android App Controller
   |-- Data Plane:    USB data channel (IP packet stream) --> Go Core (gVisor netstack)
   |                                                    |--> Session mgmt / NAT mapping
   |                                                    |--> Outbound sockets (system network path)
   |                                                    '--> Return packets back to ESP32 via USB
   '-- OTA Plane:     USB ota channel <---------------------> App OTA Manager
                                                          (chunking / verification / rollback)
Android System Network
   '-- Active VPN (system-level or third-party) as egress if configured by OS
```

### Design constraints captured
- App does not implement `VpnService`; app is forwarding-only.
- Android target is non-root.
- ESP32 uses USB custom protocol (not USB virtual NIC mode).
- Terminal experience is transparent (no proxy config).
- MVP supports both TCP and UDP traffic.
- Priority: availability first, power second, performance third.

## 2. Components and detailed flow

### 2.1 ESP32 firmware (`firmware/esp32`)
- `softap_mgr`: SoftAP, DHCP, gateway policy, client lifecycle.
- `packet_io`: capture uplink IP packets from Wi-Fi side and inject downlink packets.
- `usb_mux`: physical USB framing and channel multiplexing.
- `ctrl_agent`: config/state/log/start-stop forwarding.
- `ota_agent`: OTA chunk receive, integrity verify, partition switch, rollback.

### 2.2 Go core (`core/go`)
- `transport/usb`: Android USB host read/write, frame reassembly, flow control.
- `protocol`: protobuf message handling and version/error conventions.
- `dataplane/netstack`: gVisor netstack integration.
- `forwarder`: bridge netstack TCP/UDP sessions to outbound sockets; feed responses back.
- `controller`: device state machine, diagnostics, metrics.

### 2.3 Android app (`android/app`)
- Control panel UI: config, status, start/stop, logs, diagnostics.
- Native bridge: JNI bridge to Go core.
- USB/permissions/service: USB authorization, foreground service keepalive, reconnection.

### 2.4 Data path
Uplink:
`Terminal IP packet -> ESP32 packet capture -> USB(data) -> Go transport -> gVisor netstack -> outbound socket -> Android system network`

Downlink:
`Network response -> outbound socket -> gVisor netstack -> USB(data) -> ESP32 inject -> terminal`

## 3. Protocol design (protobuf)

Two-layer protocol model:

### 3.1 USB MUX frame layer
Header fields:
- `magic(2)`
- `version(1)`
- `channel(1)` where `1=control, 2=data, 3=ota`
- `flags(1)` (`SYN/ACK/FIN/RST/FRAG`)
- `stream_id(4)`
- `seq(4)`
- `length(4)`
- `crc32(4)`

Purpose:
- multi-channel multiplexing
- segmentation/reassembly
- flow control
- optional retransmit for data/ota channel

### 3.2 Control channel (`xshare.control.v1`)
Encoding: protobuf binary.
Common envelope fields:
- `version`
- `request_id`
- `correlation_id`
- `timestamp_ms`
- `auth_context`
- `error { code, message, details }`

### 3.3 Data channel (`xshare.data.v1`)
`PacketEnvelope` with:
- `version`
- `packet_id`
- `direction` (`uplink|downlink`)
- `if_name`
- `bytes ip_packet` (raw IPv4 first; IPv6 reserved)

### 3.4 OTA channel (`xshare.ota.v1`)
Messages:
- `OtaBegin`
- `OtaChunk`
- `OtaCommit`
- `OtaAbort`

### 3.5 Compatibility and security extension points
- Versioning by `version + capability bitmap` negotiation.
- Auth plug-in: `none | psk-hmac | noise`.
- Optional data-plane AEAD wrapper in later milestones.

### 3.6 Error code convention
- `0` OK
- `1000-1099` protocol/parameter
- `2000-2099` device state
- `3000-3099` forwarding/network runtime
- `4000-4099` OTA
- `5000-5099` internal

### 3.7 Core API set (12)
1. `device.hello`
2. `device.get_info`
3. `wifi.set_ap_config`
4. `wifi.get_clients`
5. `forward.start`
6. `forward.stop`
7. `forward.get_stats`
8. `log.subscribe`
9. `diag.ping`
10. `ota.begin`
11. `ota.chunk`
12. `ota.commit`

## 4. Monorepo structure and build/run commands

```text
XShare/
  firmware/
    esp32/
      main/
      components/
        softap_mgr/
        packet_io/
        usb_mux/
        ctrl_agent/
        ota_agent/
      partitions.csv
      sdkconfig.defaults
      CMakeLists.txt
  core/
    go/
      cmd/
        xshare-daemon/
      pkg/
        api/
        controller/
        transport/usb/
        protocol/mux/
        dataplane/netstack/
        forwarder/
        diag/
      internal/
      go.mod
  android/
    app/
    corebridge/
    build.gradle.kts
  protocol/
    proto/
      control/v1/control.proto
      data/v1/data.proto
      ota/v1/ota.proto
    buf.yaml
    buf.gen.yaml
  tools/
    proto-gen.sh
    dev-env.sh
  docs/
    plans/
    protocol/
```

Build/run baseline:
1. `cd protocol && buf generate`
2. `cd core/go && go test ./...`
3. `cd firmware/esp32 && idf.py build flash monitor`
4. `cd android && ./gradlew assembleDebug`
5. Open app, grant USB access, click Start Forwarding, verify counters/logs.

## 5. Minimal runnable demo (MVP)

### 5.1 ESP32 steps
- Start SoftAP + DHCP with ESP32 as gateway.
- Capture all terminal uplink IP packets (TCP/UDP).
- Send packets through USB data channel.
- Receive downlink packets from app and inject back to Wi-Fi interface.
- Expose control APIs: `forward.start/stop`, `get_stats`, `log.subscribe`.

### 5.2 Android + Go core steps
- Create USB host bulk endpoint communication.
- Parse USB MUX frames, dispatch data channel packets into gVisor netstack.
- Bridge netstack sessions to outbound sockets.
- Feed socket responses back to netstack, then back to ESP32 over USB.
- Keep Android side focused on control/diagnostics and lifecycle management.

### 5.3 Key interface skeletons

`core/go/pkg/transport/usb/link.go`
```go
type Link interface {
    ReadFrame(ctx context.Context) (*mux.Frame, error)
    WriteFrame(ctx context.Context, f *mux.Frame) error
}
```

`core/go/pkg/dataplane/netstack/engine.go`
```go
type Engine interface {
    InjectInbound(pkt []byte) error
    ReadOutbound(ctx context.Context) ([]byte, error)
    Start() error
    Stop() error
}
```

`core/go/pkg/forwarder/forwarder.go`
```go
type Forwarder struct { /* tcp/udp session table */ }
func (f *Forwarder) HandleFromNetstack(pkt []byte) error { /* open/reuse socket */ }
func (f *Forwarder) HandleFromSocket(pkt []byte) error   { /* write back to stack */ }
```

`firmware/esp32/components/packet_io/packet_io.h`
```c
typedef void (*uplink_cb_t)(const uint8_t* pkt, size_t len);
esp_err_t packet_io_start(uplink_cb_t cb);
esp_err_t packet_io_inject_downlink(const uint8_t* pkt, size_t len);
```

### 5.4 MVP acceptance criteria
- Terminal gets Internet access with zero manual proxy settings.
- TCP + UDP traffic works end-to-end (at minimum DNS+HTTPS+one UDP scenario).
- App can start/stop forwarding and show sessions/throughput/errors.

## 6. Risks and roadmap

### 6.1 Risks and constraints
1. Non-root Android transparent forwarding boundary:
   - App egress path depends on system network policy and active VPN behavior.
2. USB link stability:
   - Device vendor differences in power/background policy impact reliability.
3. ESP32 resource limits:
   - Keep heavy stateful transport logic on phone side.
4. gVisor integration complexity:
   - Strong correctness, but needs CPU/memory/power tuning.
5. Compatibility matrix:
   - Android version/vendor and USB controller diversity require testing matrix.

### 6.2 Milestones

M1: Transparent chain MVP
- Goal: non-root Android + USB custom protocol + transparent TCP/UDP forwarding.
- Deliverables: ESP32 SoftAP/packet bridge, Go core with gVisor path, minimal Android control app.
- Acceptance: no proxy config, DNS/HTTPS/UDP sample pass, stable 30-minute run.
- Key points: USB MUX, netstack injection/egress return, baseline diagnostics.

M2: Stability and observability
- Goal: production-like robustness across devices.
- Deliverables: reconnect state machine, flow control, structured logs, metrics panel, error taxonomy.
- Acceptance: 10+ models, 24h run, >95% auto-recovery after cable disconnect/reconnect.
- Key points: lifecycle hardening, fault-injection tests, power baseline optimization.

M3: OTA and security enhancements
- Goal: safe upgrades and pluggable auth/encryption.
- Deliverables: OTA chunk/verify/rollback, auth plugin, optional data-channel encryption.
- Acceptance: OTA success >98% with rollback on failure, compatibility checks pass.
- Key points: transactional upgrade safety, key management, version negotiation.
