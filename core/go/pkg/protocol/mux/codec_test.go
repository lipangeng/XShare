package mux

import "testing"

func TestEncodeDecodeRoundTrip(t *testing.T) {
	in := &Frame{
		Version:  1,
		Channel:  2,
		StreamID: 7,
		Seq:      1,
		Payload:  []byte{1, 2, 3},
	}
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

func TestEncodeDecodeWithFlags(t *testing.T) {
	in := &Frame{
		Version:  1,
		Channel:  ChannelControl,
		Flags:    FlagSyn,
		StreamID: 42,
		Seq:      100,
		Payload:  []byte("hello"),
	}
	b, err := Encode(in)
	if err != nil {
		t.Fatal(err)
	}
	out, err := Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	if out.Version != 1 || out.Flags != FlagSyn || out.StreamID != 42 {
		t.Fatalf("mismatch: got %+v", out)
	}
}

func TestDecodeInvalidMagic(t *testing.T) {
	_, err := Decode([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err == nil {
		t.Fatal("expected invalid magic error")
	}
}

func TestDecodeShortFrame(t *testing.T) {
	_, err := Decode([]byte{0x58, 0x53})
	if err == nil {
		t.Fatal("expected short frame error")
	}
}
