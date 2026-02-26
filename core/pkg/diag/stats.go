package diag

import "sync/atomic"

type Snapshot struct {
	UplinkPackets   uint64
	DownlinkPackets uint64
}

type Stats struct {
	uplinkPackets   atomic.Uint64
	downlinkPackets atomic.Uint64
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) IncUplinkPackets(delta uint64) {
	s.uplinkPackets.Add(delta)
}

func (s *Stats) IncDownlinkPackets(delta uint64) {
	s.downlinkPackets.Add(delta)
}

func (s *Stats) Snapshot() Snapshot {
	return Snapshot{
		UplinkPackets:   s.uplinkPackets.Load(),
		DownlinkPackets: s.downlinkPackets.Load(),
	}
}
