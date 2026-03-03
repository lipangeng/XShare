package forwarder

import (
	"testing"
)

func TestForwarderCreateUDPMapping(t *testing.T) {
	f := New()
	id := FiveTuple{SrcIP: "1.1.1.1", SrcPort: 1234, DstIP: "8.8.8.8", DstPort: 53, Proto: 17}
	_, created := f.GetOrCreate(id)
	if !created {
		t.Fatal("expected new mapping")
	}
}

func TestForwarderReuseMapping(t *testing.T) {
	f := New()
	id := FiveTuple{SrcIP: "1.1.1.1", SrcPort: 1234, DstIP: "8.8.8.8", DstPort: 53, Proto: 17}
	_, created1 := f.GetOrCreate(id)
	if !created1 {
		t.Fatal("expected new mapping on first call")
	}
	_, created2 := f.GetOrCreate(id)
	if created2 {
		t.Fatal("expected reuse on second call")
	}
}

func TestForwarderTCPMapping(t *testing.T) {
	f := New()
	id := FiveTuple{SrcIP: "192.168.1.100", SrcPort: 8080, DstIP: "93.184.216.34", DstPort: 80, Proto: 6}
	_, created := f.GetOrCreate(id)
	if !created {
		t.Fatal("expected new TCP mapping")
	}
}
