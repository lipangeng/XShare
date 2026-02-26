package usb

import (
	"bytes"
	"io"
)

// mockIO is a deterministic in-memory io.ReadWriter used by tests.
type mockIO struct {
	readBuf  bytes.Buffer
	writeBuf bytes.Buffer
}

func newMockIO(readData []byte) *mockIO {
	m := &mockIO{}
	if len(readData) > 0 {
		_, _ = m.readBuf.Write(readData)
	}
	return m
}

func (m *mockIO) Read(p []byte) (int, error) {
	if m.readBuf.Len() == 0 {
		return 0, io.EOF
	}
	return m.readBuf.Read(p)
}

func (m *mockIO) Write(p []byte) (int, error) {
	return m.writeBuf.Write(p)
}

func (m *mockIO) Written() []byte {
	out := make([]byte, m.writeBuf.Len())
	copy(out, m.writeBuf.Bytes())
	return out
}
