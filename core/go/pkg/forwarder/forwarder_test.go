package forwarder

import (
	"net/netip"
	"testing"
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
	if sess.Tuple != tuple {
		t.Fatalf("expected returned session tuple to match input tuple")
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
	if sess2 != sess {
		t.Fatalf("expected GetOrCreate to return existing session for same tuple")
	}
}
