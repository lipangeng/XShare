package usb

import (
	"context"
	"sync"
)

type mockIO struct {
	mu      sync.Mutex
	buf     []byte
	hasData bool
}

func newMockIO() *mockIO {
	return &mockIO{}
}

func (m *mockIO) Read(ctx context.Context, buf []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for !m.hasData {
		m.mu.Unlock()
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			m.mu.Lock()
		}
	}

	n := copy(buf, m.buf)
	if n < len(m.buf) {
		m.buf = m.buf[n:]
	} else {
		m.buf = nil
		m.hasData = false
	}
	return n, nil
}

func (m *mockIO) Write(ctx context.Context, data []byte) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.buf = append([]byte{}, data...)
	m.hasData = true
	return len(data), nil
}

func (m *mockIO) Close() error {
	return nil
}
