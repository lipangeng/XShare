package mux

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
)

var (
	ErrNilFrame         = errors.New("mux: nil frame")
	ErrFrameTooShort    = errors.New("mux: frame too short")
	ErrInvalidMagic     = errors.New("mux: invalid magic")
	ErrInvalidLength    = errors.New("mux: invalid payload length")
	ErrChecksumMismatch = errors.New("mux: checksum mismatch")
)

// Encode serializes a Frame into the mux wire format.
func Encode(frame *Frame) ([]byte, error) {
	if frame == nil {
		return nil, ErrNilFrame
	}
	if len(frame.Payload) > int(^uint32(0)) {
		return nil, fmt.Errorf("%w: %d", ErrInvalidLength, len(frame.Payload))
	}

	buf := make([]byte, headerSize+len(frame.Payload))
	binary.BigEndian.PutUint16(buf[0:2], wireMagic)
	buf[2] = frame.Version
	buf[3] = frame.Channel
	buf[4] = frame.Flags
	binary.BigEndian.PutUint32(buf[5:9], frame.StreamID)
	binary.BigEndian.PutUint32(buf[9:13], frame.Seq)
	binary.BigEndian.PutUint32(buf[13:17], uint32(len(frame.Payload)))
	copy(buf[headerSize:], frame.Payload)

	sum := crc32.NewIEEE()
	_, _ = sum.Write(buf[0:17])
	_, _ = sum.Write(buf[headerSize:])
	binary.BigEndian.PutUint32(buf[17:21], sum.Sum32())

	return buf, nil
}

// Decode parses a mux wire frame and validates header and CRC32.
func Decode(data []byte) (*Frame, error) {
	if len(data) < headerSize {
		return nil, ErrFrameTooShort
	}

	if binary.BigEndian.Uint16(data[0:2]) != wireMagic {
		return nil, ErrInvalidMagic
	}

	payloadLen := binary.BigEndian.Uint32(data[13:17])
	if len(data) != headerSize+int(payloadLen) {
		return nil, fmt.Errorf("%w: got %d bytes, want %d", ErrInvalidLength, len(data), headerSize+int(payloadLen))
	}

	expectedCRC := binary.BigEndian.Uint32(data[17:21])
	sum := crc32.NewIEEE()
	_, _ = sum.Write(data[0:17])
	_, _ = sum.Write(data[headerSize:])
	if expectedCRC != sum.Sum32() {
		return nil, ErrChecksumMismatch
	}

	payload := make([]byte, payloadLen)
	copy(payload, data[headerSize:])

	return &Frame{
		Version:  data[2],
		Channel:  data[3],
		Flags:    data[4],
		StreamID: binary.BigEndian.Uint32(data[5:9]),
		Seq:      binary.BigEndian.Uint32(data[9:13]),
		Payload:  payload,
	}, nil
}
