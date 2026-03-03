package mux

const (
	FrameMagic = 0x5853

	ChannelControl = 1
	ChannelData    = 2
	ChannelOta     = 3

	FlagSyn  = 0x01
	FlagAck  = 0x02
	FlagFin  = 0x04
	FlagRst  = 0x08
	FlagFrag = 0x10
)

type Frame struct {
	Version  uint8
	Channel  uint8
	Flags    uint8
	StreamID uint32
	Seq      uint32
	Payload  []byte
}
