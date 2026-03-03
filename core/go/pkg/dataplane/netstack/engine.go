package netstack

import (
	"context"
	"errors"
	"sync"
)

var ErrNotRunning = errors.New("engine not running")
var ErrPacketTooShort = errors.New("packet too short")

type Engine interface {
	InjectInbound(pkt []byte) error
	ReadOutbound(ctx context.Context) ([]byte, error)
	Start() error
	Stop() error
}

type engine struct {
	mu       sync.RWMutex
	running  bool
	outbound chan []byte
}

func NewEngine() Engine {
	return &engine{
		outbound: make(chan []byte, 100),
	}
}

func (e *engine) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.running {
		return nil
	}
	e.running = true
	return nil
}

func (e *engine) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.running {
		return nil
	}
	e.running = false
	close(e.outbound)
	return nil
}

func (e *engine) InjectInbound(pkt []byte) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if !e.running {
		return ErrNotRunning
	}
	if len(pkt) < 20 {
		return ErrPacketTooShort
	}
	return nil
}

func (e *engine) ReadOutbound(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case pkt, ok := <-e.outbound:
		if !ok {
			return nil, ErrNotRunning
		}
		return pkt, nil
	}
}
