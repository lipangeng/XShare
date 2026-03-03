package integration

import (
	"context"
	"testing"
)

func TestE2E_UplinkToDownlinkLoopback(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	s := NewHarness(t)
	defer s.Close()

	pkt := []byte{0x45, 0x00, 0x00, 0x1c, 0x00, 0x01, 0x00, 0x00, 0x40, 0x11, 0x00, 0x00, 0xc0, 0xa8, 0x01, 0x64, 0x08, 0x08, 0x08, 0x08, 0x00, 0x35, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	if err := s.InjectFromEsp(pkt); err != nil {
		t.Fatal(err)
	}

	resp, err := s.ReadToEsp(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp) == 0 {
		t.Fatal("expected response packet")
	}
}

func TestE2E_StartAndStopForwarding(t *testing.T) {
	s := NewHarness(t)
	defer s.Close()

	if err := s.StartForwarding(); err != nil {
		t.Fatal(err)
	}

	if err := s.StopForwarding(); err != nil {
		t.Fatal(err)
	}
}
