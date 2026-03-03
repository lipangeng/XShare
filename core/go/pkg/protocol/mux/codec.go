package mux

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

const (
	HeaderSize = 20
)

var (
	ErrInvalidMagic = errors.New("invalid frame magic")
	ErrShortFrame   = errors.New("frame too short")
	ErrCRCMismatch  = errors.New("CRC mismatch")
)

var crc32Table = crc32.MakeTable(crc32.IEEE)

func Encode(f *Frame) ([]byte, error) {
	payloadLen := len(f.Payload)
	frameLen := HeaderSize + payloadLen + 4

	buf := make([]byte, frameLen)

	binary.BigEndian.PutUint16(buf[0:], FrameMagic)
	buf[2] = f.Version
	buf[3] = f.Channel
	buf[4] = f.Flags
	binary.BigEndian.PutUint32(buf[5:], f.StreamID)
	binary.BigEndian.PutUint32(buf[9:], f.Seq)
	binary.BigEndian.PutUint32(buf[13:], uint32(payloadLen))

	copy(buf[HeaderSize:], f.Payload)

	crc := crc32.Checksum(buf[:HeaderSize+payloadLen], crc32Table)
	binary.BigEndian.PutUint32(buf[HeaderSize+payloadLen:], crc)

	return buf, nil
}

func Decode(buf []byte) (*Frame, error) {
	if len(buf) < HeaderSize+4 {
		return nil, ErrShortFrame
	}

	magic := binary.BigEndian.Uint16(buf[0:])
	if magic != FrameMagic {
		return nil, ErrInvalidMagic
	}

	payloadLen := binary.BigEndian.Uint32(buf[13:])
	expectedLen := HeaderSize + int(payloadLen) + 4
	if len(buf) < expectedLen {
		return nil, ErrShortFrame
	}

	storedCRC := binary.BigEndian.Uint32(buf[HeaderSize+payloadLen:])
	computedCRC := crc32.Checksum(buf[:HeaderSize+int(payloadLen)], crc32Table)
	if storedCRC != computedCRC {
		return nil, ErrCRCMismatch
	}

	f := &Frame{
		Version:  buf[2],
		Channel:  buf[3],
		Flags:    buf[4],
		StreamID: binary.BigEndian.Uint32(buf[5:]),
		Seq:      binary.BigEndian.Uint32(buf[9:]),
		Payload:  make([]byte, payloadLen),
	}
	copy(f.Payload, buf[HeaderSize:HeaderSize+int(payloadLen)])

	return f, nil
}
