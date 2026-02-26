package integration

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMockUSBTransportTakeWrittenFrameHonorsContext(t *testing.T) {
	t.Parallel()

	transport := newMockUSBTransport()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := transport.takeWrittenFrame(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("takeWrittenFrame error = %v, want %v", err, context.DeadlineExceeded)
	}
}

func TestLoopbackHarnessSendAndReceiveWriteTimeout(t *testing.T) {
	t.Parallel()

	h := newLoopbackHarness(t)
	h.transport.writtenCh <- []byte("preload")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		_, err := h.SendAndReceive(ctx, []byte{0x45, 0x00, 0x00, 0x1c})
		done <- err
	}()

	select {
	case err := <-done:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("SendAndReceive error = %v, want %v", err, context.DeadlineExceeded)
		}
	case <-time.After(200 * time.Millisecond):
		<-h.transport.writtenCh
		select {
		case <-done:
		case <-time.After(200 * time.Millisecond):
		}
		t.Fatal("SendAndReceive did not return after context timeout")
	}
}
