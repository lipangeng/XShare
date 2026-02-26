package mux

import (
	"encoding/binary"
	"errors"
	"testing"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	in := &Frame{Version: 1, Channel: 2, StreamID: 7, Seq: 1, Payload: []byte{1, 2, 3}}
	b, err := Encode(in)
	if err != nil {
		t.Fatal(err)
	}

	out, err := Decode(b)
	if err != nil {
		t.Fatal(err)
	}

	if out.Channel != in.Channel || len(out.Payload) != 3 {
		t.Fatal("roundtrip mismatch")
	}
}

func TestDecodeShortFrameReturnsFrameTooShort(t *testing.T) {
	_, err := Decode(make([]byte, headerSize-1))
	if !errors.Is(err, ErrFrameTooShort) {
		t.Fatalf("expected %v, got %v", ErrFrameTooShort, err)
	}
}

func TestDecodeInvalidMagicReturnsInvalidMagic(t *testing.T) {
	frame := &Frame{Version: 1, Channel: 2, StreamID: 3, Seq: 4, Payload: []byte{9, 8, 7}}
	encoded, err := Encode(frame)
	if err != nil {
		t.Fatal(err)
	}

	binary.BigEndian.PutUint16(encoded[0:2], 0x0000)

	_, err = Decode(encoded)
	if !errors.Is(err, ErrInvalidMagic) {
		t.Fatalf("expected %v, got %v", ErrInvalidMagic, err)
	}
}

func TestDecodeInvalidLengthReturnsInvalidLength(t *testing.T) {
	frame := &Frame{Version: 1, Channel: 2, StreamID: 3, Seq: 4, Payload: []byte{9, 8, 7}}
	encoded, err := Encode(frame)
	if err != nil {
		t.Fatal(err)
	}

	binary.BigEndian.PutUint32(encoded[13:17], uint32(len(frame.Payload)+1))

	_, err = Decode(encoded)
	if !errors.Is(err, ErrInvalidLength) {
		t.Fatalf("expected %v, got %v", ErrInvalidLength, err)
	}
}

func TestDecodeChecksumMismatchReturnsChecksumMismatch(t *testing.T) {
	frame := &Frame{Version: 1, Channel: 2, StreamID: 3, Seq: 4, Payload: []byte{9, 8, 7}}
	encoded, err := Encode(frame)
	if err != nil {
		t.Fatal(err)
	}

	encoded[len(encoded)-1] ^= 0xff

	_, err = Decode(encoded)
	if !errors.Is(err, ErrChecksumMismatch) {
		t.Fatalf("expected %v, got %v", ErrChecksumMismatch, err)
	}
}

func TestEncodeNilFrameReturnsNilFrameError(t *testing.T) {
	_, err := Encode(nil)
	if !errors.Is(err, ErrNilFrame) {
		t.Fatalf("expected %v, got %v", ErrNilFrame, err)
	}
}
