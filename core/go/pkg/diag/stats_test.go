package diag

import (
	"testing"
)

func TestStatsCounterIncrement(t *testing.T) {
	s := NewStats()
	s.IncUplinkPackets(3)
	if s.Snapshot().UplinkPackets != 3 {
		t.Fatal("bad uplink counter")
	}
}

func TestStatsBytesIncrement(t *testing.T) {
	s := NewStats()
	s.IncUplinkBytes(100)
	s.IncDownlinkBytes(200)
	stats := s.Snapshot()
	if stats.UplinkBytes != 100 || stats.DownlinkBytes != 200 {
		t.Fatal("bad byte counters")
	}
}

func TestStatsMultipleIncrements(t *testing.T) {
	s := NewStats()
	s.IncUplinkPackets(1)
	s.IncUplinkPackets(2)
	s.IncUplinkPackets(3)
	if s.Snapshot().UplinkPackets != 6 {
		t.Fatal("bad cumulative counter")
	}
}
