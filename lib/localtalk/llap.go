package localtalk

import (
	"errors"
	"fmt"

	"github.com/cheesestraws/appletalk/lib/ddp"
)

// An LLAPPacket contains information from a decoded LLAP packet.
type LLAPPacket struct {
	Src      uint8
	Dst      uint8
	LLAPType uint8
	Data     []byte
}

// EncodeBytes returns an LLAP packet as a slice of bytes ready to send to the
// physical layer.
func (l LLAPPacket) EncodeBytes() []byte {
	buf := make([]byte, 3+len(l.Data), 3+len(l.Data))
	buf[0] = l.Dst
	buf[1] = l.Src
	buf[2] = l.LLAPType
	copy(buf[3:], l.Data)
	return buf
}

var IncompleteHeader error = errors.New("incomplete LLAP header")
var Short error = errors.New("packet data too short for its declared LLAP length")
var NotDDP error = errors.New("not a DDP packet")

// DecodeLLAPPacket takes an LLAP packet as a slice of bytes read from the 
// listener and parses it into an LLAPPacket struct.
func DecodeLLAPPacket(p []byte) (*LLAPPacket, error) {
	if len(p) < 3 {
		return nil, IncompleteHeader
	}

	// parse out the easy stuff
	packet := LLAPPacket{
		Dst:      p[0],
		Src:      p[1],
		LLAPType: p[2],
	}
	packet.Data = p[3:]

	return &packet, nil
}

// PrettyHeaders returns a summary of the LLAP header of the packet suitable
// for human consumption
func (l *LLAPPacket) PrettyHeaders() string {
	return fmt.Sprintf("from %s to %s, type %s, %d data bytes",
		llapAddressString(l.Dst), llapAddressString(l.Src), llapTypeString(l.LLAPType), len(l.Data))
}

// DDP isn't finished.
func (l *LLAPPacket) DDP() (*ddp.Packet, error) {
	if l.LLAPType == 1 {
		return ddp.DecodePacket(l.Data, true, false)
	}

	return nil, NotDDP
}
