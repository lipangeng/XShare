package controller

import (
	"errors"
	"sync"

	"xshare/core/pkg/diag"
)

var (
	ErrForwardAlreadyRunning = errors.New("cannot start forwarding: already running")
	ErrForwardAlreadyIdle    = errors.New("cannot stop forwarding: controller is idle")
)

type state uint8

const (
	stateIdle state = iota
	stateRunning
)

type Controller struct {
	mu    sync.Mutex
	state state
	stats *diag.Stats
}

func NewController() *Controller {
	return &Controller{
		state: stateIdle,
		stats: diag.NewStats(),
	}
}

func (c *Controller) StartForward() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == stateRunning {
		return ErrForwardAlreadyRunning
	}

	c.state = stateRunning
	return nil
}

func (c *Controller) StopForward() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == stateIdle {
		return ErrForwardAlreadyIdle
	}

	c.state = stateIdle
	return nil
}

func (c *Controller) Stats() *diag.Stats {
	return c.stats
}
