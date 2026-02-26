package integration

import (
	"bytes"
	"context"
	"io"
	"sync"
	"testing"

	"xshare/core/pkg/protocol/mux"
	"xshare/core/pkg/transport/usb"
)

type loopbackHarness struct {
	link      usb.Link
	transport *mockUSBTransport
	responder *fakeOutboundResponder
}

func newLoopbackHarness(t *testing.T) *loopbackHarness {
	t.Helper()

	transport := newMockUSBTransport()
	return &loopbackHarness{
		link:      usb.NewLink(transport),
		transport: transport,
		responder: &fakeOutboundResponder{},
	}
}

func (h *loopbackHarness) SendAndReceive(ctx context.Context, packet []byte) ([]byte, error) {
	request := &mux.Frame{
		Version:  1,
		Channel:  1,
		Flags:    0,
		StreamID: 1,
		Seq:      1,
		Payload:  append([]byte(nil), packet...),
	}
	h.transport.setWriteContext(ctx)
	defer h.transport.setWriteContext(nil)

	if err := h.link.WriteFrame(ctx, request); err != nil {
		return nil, err
	}

	written, err := h.transport.takeWrittenFrame(ctx)
	if err != nil {
		return nil, err
	}

	responseWire, err := h.responder.respond(written)
	if err != nil {
		return nil, err
	}
	h.transport.enqueueReadableFrame(responseWire)

	response, err := h.link.ReadFrame(ctx)
	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}

type mockUSBTransport struct {
	mu       sync.Mutex
	readBuf  bytes.Buffer
	writeCtx context.Context

	writtenCh chan []byte
}

func newMockUSBTransport() *mockUSBTransport {
	return &mockUSBTransport{writtenCh: make(chan []byte, 1)}
}

func (m *mockUSBTransport) Read(p []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.readBuf.Len() == 0 {
		return 0, io.EOF
	}
	return m.readBuf.Read(p)
}

func (m *mockUSBTransport) Write(p []byte) (int, error) {
	frame := append([]byte(nil), p...)
	ctx := m.currentWriteContext()
	if ctx == nil {
		m.writtenCh <- frame
		return len(p), nil
	}

	select {
	case m.writtenCh <- frame:
		return len(p), nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func (m *mockUSBTransport) takeWrittenFrame(ctx context.Context) ([]byte, error) {
	select {
	case frame := <-m.writtenCh:
		return frame, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (m *mockUSBTransport) enqueueReadableFrame(frame []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, _ = m.readBuf.Write(frame)
}

func (m *mockUSBTransport) setWriteContext(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.writeCtx = ctx
}

func (m *mockUSBTransport) currentWriteContext() context.Context {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.writeCtx
}

type fakeOutboundResponder struct{}

func (f *fakeOutboundResponder) respond(requestWire []byte) ([]byte, error) {
	request, err := mux.Decode(requestWire)
	if err != nil {
		return nil, err
	}

	response := &mux.Frame{
		Version:  request.Version,
		Channel:  request.Channel,
		Flags:    request.Flags,
		StreamID: request.StreamID,
		Seq:      request.Seq + 1,
		Payload:  append([]byte(nil), request.Payload...),
	}

	return mux.Encode(response)
}
