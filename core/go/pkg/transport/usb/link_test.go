package usb

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"strings"
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

func TestReadExactWithContextAllowsEOFWhenBufferCompleted(t *testing.T) {
	src := &readOnceReader{
		data: []byte("hello"),
		err:  io.EOF,
	}

	buf := make([]byte, 5)
	err := readExactWithContext(context.Background(), src, buf)
	if err != nil {
		t.Fatalf("readExactWithContext returned error: %v", err)
	}
	if string(buf) != "hello" {
		t.Fatalf("unexpected payload: %q", string(buf))
	}
}

func TestLinkReadFrameRejectsOversizedPayloadLength(t *testing.T) {
	const oversizedLength = uint32(1<<20 + 1)

	header := make([]byte, muxHeaderSize)
	binary.BigEndian.PutUint32(header[muxPayloadLenFrom:muxPayloadLenTo], oversizedLength)

	link := NewLink(newMockIO(header))
	_, err := link.ReadFrame(context.Background())
	if err == nil {
		t.Fatal("ReadFrame returned nil error")
	}
	if !errors.Is(err, mux.ErrInvalidLength) {
		t.Fatalf("expected ErrInvalidLength, got %v", err)
	}
	if !strings.Contains(err.Error(), "exceeds max") {
		t.Fatalf("expected max length message, got: %v", err)
	}
}

func TestLinkReadFrameCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	link := NewLink(newMockIO(nil))
	_, err := link.ReadFrame(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestLinkReadFramePropagatesShortPayloadReadError(t *testing.T) {
	frame := &mux.Frame{
		Version:  1,
		Channel:  1,
		Flags:    0,
		StreamID: 1,
		Seq:      1,
		Payload:  []byte("hello"),
	}
	encoded, err := mux.Encode(frame)
	if err != nil {
		t.Fatalf("Encode returned error: %v", err)
	}
	truncated := encoded[:len(encoded)-1]

	link := NewLink(newMockIO(truncated))
	_, err = link.ReadFrame(context.Background())
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected io.EOF, got %v", err)
	}
}

type readOnceReader struct {
	data []byte
	err  error
	read bool
}

func (r *readOnceReader) Read(p []byte) (int, error) {
	if r.read {
		return 0, io.EOF
	}
	r.read = true
	n := copy(p, r.data)
	return n, r.err
}
