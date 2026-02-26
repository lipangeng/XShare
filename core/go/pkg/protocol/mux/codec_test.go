package mux

import "testing"

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
