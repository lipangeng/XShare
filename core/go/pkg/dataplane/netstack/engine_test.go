package netstack

import (
	"errors"
	"testing"
)

func TestEngineInjectInboundRejectsShortPacket(t *testing.T) {
	eng := NewEngine()

	if err := eng.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer func() {
		_ = eng.Stop()
	}()

	err := eng.InjectInbound([]byte{0x45, 0x00, 0x00})
	if !errors.Is(err, ErrShortPacket) {
		t.Fatalf("expected ErrShortPacket, got %v", err)
	}
}
