package controller

import (
	"errors"
	"testing"
)

func TestForwardControllerPublicContract(t *testing.T) {
	ctrl := NewController()

	before := ctrl.Stats()
	if before.UplinkPackets != 0 {
		t.Fatalf("expected zero uplink packets before start, got %d", before.UplinkPackets)
	}
	if before.DownlinkPackets != 0 {
		t.Fatalf("expected zero downlink packets before start, got %d", before.DownlinkPackets)
	}

	if err := ctrl.StartForward(); err != nil {
		t.Fatalf("expected StartForward to succeed, got: %v", err)
	}

	if err := ctrl.StartForward(); !errors.Is(err, ErrForwardAlreadyRunning) {
		t.Fatalf("expected ErrForwardAlreadyRunning on repeated start, got: %v", err)
	}

	during := ctrl.Stats()
	if during.UplinkPackets != 0 {
		t.Fatalf("expected zero uplink packets while running, got %d", during.UplinkPackets)
	}
	if during.DownlinkPackets != 0 {
		t.Fatalf("expected zero downlink packets while running, got %d", during.DownlinkPackets)
	}

	if err := ctrl.StopForward(); err != nil {
		t.Fatalf("expected StopForward to succeed, got: %v", err)
	}

	if err := ctrl.StopForward(); !errors.Is(err, ErrForwardAlreadyIdle) {
		t.Fatalf("expected ErrForwardAlreadyIdle on repeated stop, got: %v", err)
	}

	after := ctrl.Stats()
	if after.UplinkPackets != 0 {
		t.Fatalf("expected zero uplink packets after stop, got %d", after.UplinkPackets)
	}
	if after.DownlinkPackets != 0 {
		t.Fatalf("expected zero downlink packets after stop, got %d", after.DownlinkPackets)
	}
}
