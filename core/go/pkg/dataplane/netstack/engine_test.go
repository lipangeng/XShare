package netstack

import (
	"testing"
)

func TestEngineInjectAndRead(t *testing.T) {
	e := NewEngine()
	if err := e.Start(); err != nil {
		t.Fatal(err)
	}
	defer e.Stop()
	if err := e.InjectInbound([]byte{0x45, 0, 0, 20}); err == nil {
		t.Fatal("expected short packet error")
	}
}

func TestEngineStartStop(t *testing.T) {
	e := NewEngine()
	if err := e.Start(); err != nil {
		t.Fatal(err)
	}
	if err := e.Stop(); err != nil {
		t.Fatal(err)
	}
}
