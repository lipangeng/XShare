package integration

import (
	"context"
	"os"
	"testing"

	"github.com/xshare/xshare/pkg/controller"
	"github.com/xshare/xshare/pkg/dataplane/netstack"
	"github.com/xshare/xshare/pkg/forwarder"
	"github.com/xshare/xshare/pkg/protocol/mux"
	"github.com/xshare/xshare/pkg/transport/usb"
)

type Harness struct {
	t          *testing.T
	ctrl       *controller.Controller
	engine     netstack.Engine
	fwd        *forwarder.Forwarder
	link       usb.Link
	mockIO     *mockUSB
}

type mockUSB struct {
	writeBuffer []byte
	readBuffer  []byte
}

func (m *mockUSB) Read(ctx context.Context, buf []byte) (int, error) {
	if len(m.readBuffer) == 0 {
		<-ctx.Done()
		return 0, ctx.Err()
	}
	n := copy(buf, m.readBuffer)
	m.readBuffer = m.readBuffer[n:]
	return n, nil
}

func (m *mockUSB) Write(ctx context.Context, data []byte) (int, error) {
	m.writeBuffer = append(m.writeBuffer, data...)
	m.readBuffer = append(m.readBuffer, data...)
	return len(data), nil
}

func (m *mockUSB) Close() error {
	return nil
}

func NewHarness(t *testing.T) *Harness {
	t.Helper()

	mock := &mockUSB{}

	h := &Harness{
		t:      t,
		ctrl:   controller.NewController(),
		engine: netstack.NewEngine(),
		fwd:    forwarder.New(),
		mockIO: mock,
	}

	link := usb.NewLink(mock)
	h.link = link

	return h
}

func (h *Harness) Close() {
	if h.ctrl != nil {
		h.ctrl.StopForward()
	}
	if h.engine != nil {
		h.engine.Stop()
	}
}

func (h *Harness) InjectFromEsp(pkt []byte) error {
	frame := &mux.Frame{
		Channel:  mux.ChannelData,
		Payload:  pkt,
		StreamID: 1,
		Seq:      1,
		Version:  1,
	}
	return h.link.WriteFrame(context.Background(), frame)
}

func (h *Harness) ReadToEsp(ctx context.Context) ([]byte, error) {
	frame, err := h.link.ReadFrame(ctx)
	if err != nil {
		return nil, err
	}
	return frame.Payload, nil
}

func (h *Harness) StartForwarding() error {
	return h.ctrl.StartForward()
}

func (h *Harness) StopForwarding() error {
	return h.ctrl.StopForward()
}

func loadPacket(t *testing.T, path string) []byte {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	if len(data) == 0 {
		data = []byte{0x45, 0x00, 0x00, 0x1c, 0x00, 0x01, 0x00, 0x00, 0x40, 0x11, 0x00, 0x00, 0xc0, 0xa8, 0x01, 0x64, 0x08, 0x08, 0x08, 0x08, 0x00, 0x35, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	}

	return data
}
