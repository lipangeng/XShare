package diag

import "testing"

func TestStatsCounterIncrement(t *testing.T) {
	s := NewStats()
	s.IncUplinkPackets(3)
	if s.Snapshot().UplinkPackets != 3 {
		t.Fatal("bad uplink counter")
	}
}
