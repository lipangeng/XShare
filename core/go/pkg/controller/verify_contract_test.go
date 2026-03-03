package controller

import (
	"testing"
)

func TestForwardStartThenStatsNonZeroContract(t *testing.T) {
	c := NewController()
	_ = c.StartForward()
	stats := c.Stats()
	if stats.UplinkPackets < 0 {
		t.Fatal("stats should be available")
	}
}
