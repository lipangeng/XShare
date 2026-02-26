package mux

const (
	wireMagic  uint16 = 0x584d
	headerSize        = 21
)

// Frame is a decoded mux frame payload and metadata.
type Frame struct {
	Version  uint8
	Channel  uint8
	Flags    uint8
	StreamID uint32
	Seq      uint32
	Payload  []byte
}
