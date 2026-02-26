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
	if err := h.link.WriteFrame(ctx, request); err != nil {
		return nil, err
	}

	written := h.transport.takeWrittenFrame()
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
	mu        sync.Mutex
	readBuf   bytes.Buffer
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
	m.writtenCh <- frame
	return len(p), nil
}

func (m *mockUSBTransport) takeWrittenFrame() []byte {
	return <-m.writtenCh
}

func (m *mockUSBTransport) enqueueReadableFrame(frame []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, _ = m.readBuf.Write(frame)
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
