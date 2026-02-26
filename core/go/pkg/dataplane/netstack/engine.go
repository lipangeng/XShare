package netstack

import (
	"context"
	"errors"
	"sync"
)

const minInboundPacketSize = 20

var (
	ErrShortPacket       = errors.New("netstack: short packet")
	ErrEngineStopped     = errors.New("netstack: engine stopped")
	ErrOutboundQueueFull = errors.New("netstack: outbound queue full")
)

// Engine is a deterministic scaffold for future netstack integration.
type Engine struct {
	mu       sync.RWMutex
	started  bool
	stopCh   chan struct{}
	outbound chan []byte
}

func NewEngine() *Engine {
	return &Engine{}
}

func (e *Engine) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.started {
		return nil
	}

	e.stopCh = make(chan struct{})
	e.outbound = make(chan []byte, 64)
	e.started = true
	return nil
}

func (e *Engine) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.started {
		return nil
	}

	e.started = false
	close(e.stopCh)
	e.stopCh = nil
	e.outbound = nil
	return nil
}

func (e *Engine) InjectInbound(packet []byte) error {
	if len(packet) < minInboundPacketSize {
		return ErrShortPacket
	}

	copyPacket := append([]byte(nil), packet...)

	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.started || e.outbound == nil {
		return ErrEngineStopped
	}

	select {
	case e.outbound <- copyPacket:
		return nil
	default:
		if !e.started {
			return ErrEngineStopped
		}
		return ErrOutboundQueueFull
	}
}

func (e *Engine) ReadOutbound(ctx context.Context) ([]byte, error) {
	e.mu.RLock()
	started := e.started
	stopCh := e.stopCh
	e.mu.RUnlock()

	if !started {
		return nil, ErrEngineStopped
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-stopCh:
		return nil, ErrEngineStopped
	case packet := <-e.outbound:
		return packet, nil
	}
}
