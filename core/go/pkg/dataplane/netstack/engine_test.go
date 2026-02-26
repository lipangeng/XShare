package netstack

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestEngineStartStopIdempotent(t *testing.T) {
	eng := NewEngine()

	if err := eng.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	if err := eng.Start(); err != nil {
		t.Fatalf("second start: %v", err)
	}

	if err := eng.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}

	if err := eng.Stop(); err != nil {
		t.Fatalf("second stop: %v", err)
	}
}

func TestEngineReadOutboundStopped(t *testing.T) {
	eng := NewEngine()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := eng.ReadOutbound(ctx); !errors.Is(err, ErrEngineStopped) {
		t.Fatalf("before start: expected ErrEngineStopped, got %v", err)
	}

	if err := eng.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := eng.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}

	if _, err := eng.ReadOutbound(ctx); !errors.Is(err, ErrEngineStopped) {
		t.Fatalf("after stop: expected ErrEngineStopped, got %v", err)
	}
}

func TestEngineReadOutboundContextCanceled(t *testing.T) {
	eng := NewEngine()
	if err := eng.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer func() {
		_ = eng.Stop()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := eng.ReadOutbound(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestEngineRestartDoesNotLeakQueuedPackets(t *testing.T) {
	eng := NewEngine()
	if err := eng.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}

	packet := make([]byte, minInboundPacketSize)
	packet[0] = 0x45
	if err := eng.InjectInbound(packet); err != nil {
		t.Fatalf("inject: %v", err)
	}
	if err := eng.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}
	if err := eng.Start(); err != nil {
		t.Fatalf("restart: %v", err)
	}
	defer func() {
		_ = eng.Stop()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	if _, err := eng.ReadOutbound(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline (no leaked packet), got %v", err)
	}
}

func TestEngineInjectInboundAndReadOutboundSuccess(t *testing.T) {
	eng := NewEngine()
	if err := eng.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	defer func() {
		_ = eng.Stop()
	}()

	want := make([]byte, minInboundPacketSize)
	want[0] = 0x45
	if err := eng.InjectInbound(want); err != nil {
		t.Fatalf("inject: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	got, err := eng.ReadOutbound(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("unexpected packet: got=%v want=%v", got, want)
	}
}

func TestEngineInjectInboundAfterStopReturnsEngineStopped(t *testing.T) {
	eng := NewEngine()
	if err := eng.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := eng.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}

	packet := make([]byte, minInboundPacketSize)
	packet[0] = 0x45
	if err := eng.InjectInbound(packet); !errors.Is(err, ErrEngineStopped) {
		t.Fatalf("expected ErrEngineStopped, got %v", err)
	}
}

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
