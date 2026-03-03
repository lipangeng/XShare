package usb

import (
	"context"

	"github.com/xshare/xshare/pkg/protocol/mux"
)

type Link interface {
	ReadFrame(ctx context.Context) (*mux.Frame, error)
	WriteFrame(ctx context.Context, f *mux.Frame) error
	Close() error
}

type IO interface {
	Read(ctx context.Context, buf []byte) (int, error)
	Write(ctx context.Context, buf []byte) (int, error)
	Close() error
}

type link struct {
	io IO
}

func NewLink(io IO) Link {
	return &link{io: io}
}

func (l *link) ReadFrame(ctx context.Context) (*mux.Frame, error) {
	header := make([]byte, mux.HeaderSize)
	n, err := l.io.Read(ctx, header)
	if err != nil {
		return nil, err
	}
	if n < mux.HeaderSize {
		return nil, mux.ErrShortFrame
	}

	payloadLen := int(header[13])<<24 | int(header[14])<<16 | int(header[15])<<8 | int(header[16])

	frameSize := mux.HeaderSize + payloadLen + 4
	fullBuf := make([]byte, frameSize)
	copy(fullBuf, header)

	_, err = l.io.Read(ctx, fullBuf[mux.HeaderSize:])
	if err != nil {
		return nil, err
	}

	return mux.Decode(fullBuf)
}

func (l *link) WriteFrame(ctx context.Context, f *mux.Frame) error {
	buf, err := mux.Encode(f)
	if err != nil {
		return err
	}
	_, err = l.io.Write(ctx, buf)
	return err
}

func (l *link) Close() error {
	return l.io.Close()
}
