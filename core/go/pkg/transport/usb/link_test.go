package usb

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"xshare/core/pkg/protocol/mux"
)

func TestLinkWriteFrameEncodesToStream(t *testing.T) {
	stream := newMockIO(nil)
	link := NewLink(stream)

	in := &mux.Frame{
		Version:  1,
		Channel:  2,
		Flags:    3,
		StreamID: 4,
		Seq:      5,
		Payload:  []byte("hello"),
	}

	if err := link.WriteFrame(context.Background(), in); err != nil {
		t.Fatalf("WriteFrame returned error: %v", err)
	}

	want, err := mux.Encode(in)
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}

	if !bytes.Equal(stream.Written(), want) {
		t.Fatalf("written bytes mismatch")
	}
}

func TestLinkReadFrameDecodesFromStream(t *testing.T) {
	want := &mux.Frame{
		Version:  1,
		Channel:  9,
		Flags:    7,
		StreamID: 99,
		Seq:      123,
		Payload:  []byte("payload"),
	}
	encoded, err := mux.Encode(want)
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}

	stream := newMockIO(encoded)
	link := NewLink(stream)

	got, err := link.ReadFrame(context.Background())
	if err != nil {
		t.Fatalf("ReadFrame returned error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ReadFrame mismatch: got %+v want %+v", got, want)
	}
}
