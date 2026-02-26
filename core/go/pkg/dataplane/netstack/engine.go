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
	return &Engine{
		outbound: make(chan []byte, 64),
	}
}

func (e *Engine) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.started {
		return nil
	}

	e.started = true
	e.stopCh = make(chan struct{})
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
	return nil
}

func (e *Engine) InjectInbound(packet []byte) error {
	if len(packet) < minInboundPacketSize {
		return ErrShortPacket
	}

	e.mu.RLock()
	started := e.started
	stopCh := e.stopCh
	e.mu.RUnlock()

	if !started {
		return ErrEngineStopped
	}

	copyPacket := append([]byte(nil), packet...)
	select {
	case e.outbound <- copyPacket:
		return nil
	case <-stopCh:
		return ErrEngineStopped
	default:
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
