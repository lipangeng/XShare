package usb

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"xshare/core/pkg/protocol/mux"
)

const (
	muxHeaderSize     = 21
	muxPayloadLenFrom = 13
	muxPayloadLenTo   = 17
)

// Link defines a minimal USB transport capable of sending and receiving mux frames.
type Link interface {
	ReadFrame(ctx context.Context) (*mux.Frame, error)
	WriteFrame(ctx context.Context, frame *mux.Frame) error
}

type streamLink struct {
	rw io.ReadWriter
}

// NewLink creates a stream-backed USB link.
func NewLink(rw io.ReadWriter) Link {
	return &streamLink{rw: rw}
}

func (l *streamLink) ReadFrame(ctx context.Context) (*mux.Frame, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	header := make([]byte, muxHeaderSize)
	if err := readExactWithContext(ctx, l.rw, header); err != nil {
		return nil, err
	}

	payloadLen := binary.BigEndian.Uint32(header[muxPayloadLenFrom:muxPayloadLenTo])
	if int(payloadLen) < 0 || int(payloadLen) > maxPayloadSize() {
		return nil, fmt.Errorf("%w: %d", mux.ErrInvalidLength, payloadLen)
	}

	frameBytes := make([]byte, muxHeaderSize+int(payloadLen))
	copy(frameBytes, header)
	if payloadLen > 0 {
		if err := readExactWithContext(ctx, l.rw, frameBytes[muxHeaderSize:]); err != nil {
			return nil, err
		}
	}

	return mux.Decode(frameBytes)
}

func (l *streamLink) WriteFrame(ctx context.Context, frame *mux.Frame) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	encoded, err := mux.Encode(frame)
	if err != nil {
		return err
	}

	return writeAllWithContext(ctx, l.rw, encoded)
}

func readExactWithContext(ctx context.Context, r io.Reader, buf []byte) error {
	for off := 0; off < len(buf); {
		if err := ctx.Err(); err != nil {
			return err
		}

		n, err := r.Read(buf[off:])
		off += n
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrNoProgress
		}
	}
	return nil
}

func writeAllWithContext(ctx context.Context, w io.Writer, buf []byte) error {
	for off := 0; off < len(buf); {
		if err := ctx.Err(); err != nil {
			return err
		}

		n, err := w.Write(buf[off:])
		off += n
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrNoProgress
		}
	}
	return nil
}

func maxPayloadSize() int {
	maxInt := int(^uint(0) >> 1)
	return maxInt - muxHeaderSize
}
