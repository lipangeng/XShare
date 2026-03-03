package forwarder

import (
	"fmt"
	"sync"
	"time"
)

type FiveTuple struct {
	SrcIP   string
	SrcPort uint16
	DstIP   string
	DstPort uint16
	Proto   uint8
}

type Session struct {
	ID        FiveTuple
	CreatedAt time.Time
	LastSeen  time.Time
}

type Forwarder struct {
	mu       sync.RWMutex
	sessions map[FiveTuple]*Session
}

func New() *Forwarder {
	return &Forwarder{
		sessions: make(map[FiveTuple]*Session),
	}
}

func (f *Forwarder) GetOrCreate(id FiveTuple) (*Session, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	if sess, exists := f.sessions[id]; exists {
		sess.LastSeen = now
		return sess, false
	}

	sess := &Session{
		ID:        id,
		CreatedAt: now,
		LastSeen:  now,
	}
	f.sessions[id] = sess
	return sess, true
}

func (f *Forwarder) Get(id FiveTuple) (*Session, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	sess, exists := f.sessions[id]
	return sess, exists
}

func (f *Forwarder) Delete(id FiveTuple) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.sessions, id)
}

func (f *Forwarder) Count() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.sessions)
}

func (id FiveTuple) String() string {
	return fmt.Sprintf("%s:%d->%s:%d proto=%d", id.SrcIP, id.SrcPort, id.DstIP, id.DstPort, id.Proto)
}
