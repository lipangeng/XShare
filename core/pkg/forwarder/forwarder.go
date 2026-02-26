package forwarder

import (
	"net/netip"
	"strings"
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

	tuple = canonicalizeTuple(tuple)
	now := time.Now()
	if sess, ok := f.sessions[tuple]; ok {
		sess.LastSeen = now
		return copySession(sess), false
	}

	sess := &Session{
		Tuple:     tuple,
		CreatedAt: now,
		LastSeen:  now,
	}
	f.sessions[tuple] = sess
	return copySession(sess), true
}

func (f *Forwarder) SessionTimeout() time.Duration {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.sessionTimeout
}

func canonicalizeTuple(tuple FiveTuple) FiveTuple {
	tuple.Proto = strings.ToUpper(tuple.Proto)
	return tuple
}

func copySession(sess *Session) *Session {
	if sess == nil {
		return nil
	}
	clone := *sess
	return &clone
}
