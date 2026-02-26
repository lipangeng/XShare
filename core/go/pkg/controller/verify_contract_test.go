package controller

import "testing"

func TestForwardStartThenStatsAvailableContract(t *testing.T) {
	ctrl := NewController()

	if err := ctrl.StartForward(); err != nil {
		t.Fatalf("expected StartForward to succeed, got: %v", err)
	}

	stats := ctrl.Stats()
	if stats.UplinkPackets != 0 {
		t.Fatalf("expected zero uplink packets at startup, got %d", stats.UplinkPackets)
	}
	if stats.DownlinkPackets != 0 {
		t.Fatalf("expected zero downlink packets at startup, got %d", stats.DownlinkPackets)
	}
}
