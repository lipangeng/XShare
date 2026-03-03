package usb

import (
	"context"
	"testing"

	"github.com/xshare/xshare/pkg/protocol/mux"
)

func TestLinkReadWriteFrame(t *testing.T) {
	io := newMockIO()
	l := NewLink(io)
	f := &mux.Frame{
		Version:  1,
		Channel:  2,
		Payload:  []byte("abc"),
		StreamID: 1,
		Seq:      1,
	}
	if err := l.WriteFrame(context.Background(), f); err != nil {
		t.Fatal(err)
	}
	got, err := l.ReadFrame(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if string(got.Payload) != "abc" {
		t.Fatal("payload mismatch")
	}
}
