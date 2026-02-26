package integration

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoopbackE2EUDPIPv4(t *testing.T) {
	t.Parallel()

	packetPath := filepath.Join("testdata", "sample_ipv4_udp.bin")
	packet, err := os.ReadFile(packetPath)
	if err != nil {
		t.Fatalf("read sample packet: %v", err)
	}

	h := newLoopbackHarness(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	got, err := h.SendAndReceive(ctx, packet)
	if err != nil {
		t.Fatalf("SendAndReceive returned error: %v", err)
	}

	if !bytes.Equal(got, packet) {
		t.Fatalf("loopback mismatch: got %d bytes want %d", len(got), len(packet))
	}
}
