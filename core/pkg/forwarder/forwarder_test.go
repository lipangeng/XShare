package forwarder

import (
	"net/netip"
	"sync"
	"testing"
	"time"
)

func TestGetOrCreate_UDPTuple_FirstCallCreatesSession(t *testing.T) {
	fwd := NewForwarder()
	tuple := FiveTuple{
		SrcIP:   netip.MustParseAddr("10.0.0.2"),
		SrcPort: 5050,
		DstIP:   netip.MustParseAddr("10.0.0.9"),
		DstPort: 53,
		Proto:   "udp",
	}

	sess, created := fwd.GetOrCreate(tuple)
	if !created {
		t.Fatalf("expected first GetOrCreate call to create session")
	}
	if sess == nil {
		t.Fatalf("expected session to be non-nil")
	}
	expectedTuple := tuple
	expectedTuple.Proto = "UDP"
	if sess.Tuple != expectedTuple {
		t.Fatalf("expected returned session tuple to match canonicalized input tuple")
	}
	if sess.CreatedAt.IsZero() {
		t.Fatalf("expected CreatedAt to be set")
	}
	if sess.LastSeen.IsZero() {
		t.Fatalf("expected LastSeen to be set")
	}

	sess2, created2 := fwd.GetOrCreate(tuple)
	if created2 {
		t.Fatalf("expected second GetOrCreate call to return created=false")
	}
	if sess2 == sess {
		t.Fatalf("expected GetOrCreate to return a copy, not internal session pointer")
	}
	if sess2.Tuple != sess.Tuple {
		t.Fatalf("expected session tuple to remain consistent on repeated GetOrCreate")
	}
}

func TestGetOrCreate_NormalizesProtocolKey(t *testing.T) {
	fwd := NewForwarder()
	lower := FiveTuple{
		SrcIP:   netip.MustParseAddr("10.1.0.2"),
		SrcPort: 4000,
		DstIP:   netip.MustParseAddr("10.1.0.9"),
		DstPort: 53,
		Proto:   "udp",
	}
	upper := lower
	upper.Proto = "UDP"

	sess1, created1 := fwd.GetOrCreate(lower)
	if !created1 {
		t.Fatalf("expected first call to create a session")
	}

	sess2, created2 := fwd.GetOrCreate(upper)
	if created2 {
		t.Fatalf("expected protocol variants to map to the same session")
	}
	if sess1.Tuple.Proto != "UDP" || sess2.Tuple.Proto != "UDP" {
		t.Fatalf("expected canonical protocol to be UDP, got %q and %q", sess1.Tuple.Proto, sess2.Tuple.Proto)
	}
	if len(fwd.sessions) != 1 {
		t.Fatalf("expected exactly one session entry, got %d", len(fwd.sessions))
	}
}

func TestGetOrCreate_UpdatesLastSeenOnSecondCall(t *testing.T) {
	fwd := NewForwarder()
	tuple := FiveTuple{
		SrcIP:   netip.MustParseAddr("10.2.0.2"),
		SrcPort: 5051,
		DstIP:   netip.MustParseAddr("10.2.0.9"),
		DstPort: 80,
		Proto:   "tcp",
	}

	sess1, created1 := fwd.GetOrCreate(tuple)
	if !created1 {
		t.Fatalf("expected first call to create session")
	}

	time.Sleep(2 * time.Millisecond)

	sess2, created2 := fwd.GetOrCreate(tuple)
	if created2 {
		t.Fatalf("expected second call to return created=false")
	}
	if !sess2.LastSeen.After(sess1.LastSeen) {
		t.Fatalf("expected LastSeen to advance; first=%v second=%v", sess1.LastSeen, sess2.LastSeen)
	}
}

func TestGetOrCreate_ConcurrentSameTupleCreatesOnlyOnce(t *testing.T) {
	fwd := NewForwarder()
	tuple := FiveTuple{
		SrcIP:   netip.MustParseAddr("10.3.0.2"),
		SrcPort: 5052,
		DstIP:   netip.MustParseAddr("10.3.0.9"),
		DstPort: 443,
		Proto:   "udp",
	}

	const workers = 32
	var wg sync.WaitGroup
	results := make(chan bool, workers)

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_, created := fwd.GetOrCreate(tuple)
			results <- created
		}()
	}
	wg.Wait()
	close(results)

	createdCount := 0
	for created := range results {
		if created {
			createdCount++
		}
	}
	if createdCount != 1 {
		t.Fatalf("expected exactly one creator, got %d", createdCount)
	}
}
