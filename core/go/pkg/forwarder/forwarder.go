package forwarder

import (
	"net/netip"
	"sync"
	"time"
)

const defaultSessionTimeout = 2 * time.Minute

// FiveTuple identifies a transport flow.
type FiveTuple struct {
	SrcIP   netip.Addr
	SrcPort uint16
	DstIP   netip.Addr
	DstPort uint16
	Proto   string
}

// Session stores flow state and timeout-related metadata.
type Session struct {
	Tuple     FiveTuple
	CreatedAt time.Time
	LastSeen  time.Time
}

// Forwarder maintains an in-memory session table.
type Forwarder struct {
	mu             sync.Mutex
	sessionTimeout time.Duration
	sessions       map[FiveTuple]*Session
}

func NewForwarder() *Forwarder {
	return &Forwarder{
		sessionTimeout: defaultSessionTimeout,
		sessions:       make(map[FiveTuple]*Session),
	}
}

func (f *Forwarder) GetOrCreate(tuple FiveTuple) (*Session, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	if sess, ok := f.sessions[tuple]; ok {
		sess.LastSeen = now
		return sess, false
	}

	sess := &Session{
		Tuple:     tuple,
		CreatedAt: now,
		LastSeen:  now,
	}
	f.sessions[tuple] = sess
	return sess, true
}

func (f *Forwarder) SessionTimeout() time.Duration {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.sessionTimeout
}
