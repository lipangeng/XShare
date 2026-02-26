package controller

import (
	"strings"
	"testing"
)

func TestController_StartStop_StateTransitions(t *testing.T) {
	ctrl := NewController()

	if err := ctrl.StartForward(); err != nil {
		t.Fatalf("expected first StartForward to succeed, got error: %v", err)
	}

	err := ctrl.StartForward()
	if err == nil {
		t.Fatalf("expected second StartForward to fail when already running")
	}
	if !strings.Contains(err.Error(), "already running") {
		t.Fatalf("expected running-state error message, got: %v", err)
	}

	if err := ctrl.StopForward(); err != nil {
		t.Fatalf("expected StopForward to succeed from running state, got error: %v", err)
	}
}

func TestController_StopForward_WhenIdle_ReturnsError(t *testing.T) {
	ctrl := NewController()

	err := ctrl.StopForward()
	if err == nil {
		t.Fatalf("expected StopForward to fail when controller is idle")
	}
	if !strings.Contains(err.Error(), "idle") {
		t.Fatalf("expected idle-state error message, got: %v", err)
	}
}
