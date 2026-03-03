package controller

import (
	"testing"
)

func TestStartStopStateTransition(t *testing.T) {
	c := NewController()
	if err := c.StartForward(); err != nil {
		t.Fatal(err)
	}
	if err := c.StopForward(); err != nil {
		t.Fatal(err)
	}
}

func TestDoubleStartError(t *testing.T) {
	c := NewController()
	if err := c.StartForward(); err != nil {
		t.Fatal(err)
	}
	if err := c.StartForward(); err == nil {
		t.Fatal("expected error on double start")
	}
}

func TestStopWithoutStartError(t *testing.T) {
	c := NewController()
	if err := c.StopForward(); err == nil {
		t.Fatal("expected error on stop without start")
	}
}
