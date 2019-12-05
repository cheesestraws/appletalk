package localtalk

import (
	"errors"
	"fmt"

	"github.com/cheesestraws/appletalk/lib/ddp"
)

type LLAPPacket struct {
	Src      uint8
	Dst      uint8
	LLAPType uint8
	Length   uint16
	Data     []byte
}

var IncompleteHeader error = errors.New("incomplete LLAP header")
var Short error = errors.New("packet data too short for its declared LLAP length")
var NotDDP error = errors.New("not a DDP packet")

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

	if len(p) > 5 {
		size1 := uint16(p[3]&3) << 8
		size2 := uint16(p[4])
		size := size1 + size2

		if len(p) < int(size)+3 {
			return nil, Short
		}

		packet.Length = size
		packet.Data = p[3 : size+3]
	}

	return &packet, nil
}

func (l *LLAPPacket) PrettyHeaders() string {
	return fmt.Sprintf("from %s to %s, type %s, %d data bytes",
		llapAddressString(l.Dst), llapAddressString(l.Src), llapTypeString(l.LLAPType), l.Length)
}

func (l *LLAPPacket) DDP() (*ddp.Packet, error) {
	if l.LLAPType == 1 {
		return ddp.DecodePacket(l.Data, true, false)
	}

	return nil, NotDDP
}
