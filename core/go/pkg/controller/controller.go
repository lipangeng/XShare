package controller

import (
	"errors"
	"sync"
)

var (
	ErrAlreadyRunning = errors.New("already running")
	ErrNotRunning     = errors.New("not running")
)

type State int

const (
	StateIdle State = iota
	StateRunning
)

type Controller struct {
	mu    sync.Mutex
	state State
}

func NewController() *Controller {
	return &Controller{
		state: StateIdle,
	}
}

func (c *Controller) StartForward() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == StateRunning {
		return ErrAlreadyRunning
	}
	c.state = StateRunning
	return nil
}

func (c *Controller) StopForward() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == StateIdle {
		return ErrNotRunning
	}
	c.state = StateIdle
	return nil
}

func (c *Controller) State() State {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

func (c *Controller) IsRunning() bool {
	return c.State() == StateRunning
}
