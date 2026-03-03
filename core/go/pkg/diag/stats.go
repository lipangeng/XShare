package diag

import "sync/atomic"

type Stats struct {
	uplinkPackets   atomic.Int64
	uplinkBytes     atomic.Int64
	downlinkPackets atomic.Int64
	downlinkBytes   atomic.Int64
	activeTCPSess   atomic.Int32
	activeUDPSess   atomic.Int32
}

type StatsSnapshot struct {
	UplinkPackets     int64
	UplinkBytes       int64
	DownlinkPackets   int64
	DownlinkBytes     int64
	ActiveTCPSessions int32
	ActiveUDPSessions int32
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) IncUplinkPackets(n int64) {
	s.uplinkPackets.Add(n)
}

func (s *Stats) IncUplinkBytes(n int64) {
	s.uplinkBytes.Add(n)
}

func (s *Stats) IncDownlinkPackets(n int64) {
	s.downlinkPackets.Add(n)
}

func (s *Stats) IncDownlinkBytes(n int64) {
	s.downlinkBytes.Add(n)
}

func (s *Stats) IncActiveTCPSessions(n int32) {
	s.activeTCPSess.Add(n)
}

func (s *Stats) IncActiveUDPSessions(n int32) {
	s.activeUDPSess.Add(n)
}

func (s *Stats) Snapshot() StatsSnapshot {
	return StatsSnapshot{
		UplinkPackets:     s.uplinkPackets.Load(),
		UplinkBytes:       s.uplinkBytes.Load(),
		DownlinkPackets:   s.downlinkPackets.Load(),
		DownlinkBytes:     s.downlinkBytes.Load(),
		ActiveTCPSessions: s.activeTCPSess.Load(),
		ActiveUDPSessions: s.activeUDPSess.Load(),
	}
}
