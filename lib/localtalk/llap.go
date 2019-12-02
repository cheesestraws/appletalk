package localtalk

import (
	"errors"
	"fmt"
)

type LLAPPacket struct {
	Src uint8
	Dst uint8
	LLAPType uint8
	Length uint16
	Data []byte
}

var IncompleteHeader error = errors.New("incomplete LLAP header")

func DecodeLLAPFrame(frame []byte) (*LLAPPacket, error) {
	if len(frame) < 5 {
		return nil, IncompleteHeader
	}
	
	// parse out the easy stuff
	packet := LLAPPacket{
		Dst: frame[0],
		Src: frame[1],
		LLAPType: frame[2],
	}
	
	size1 := uint16(frame[3] & 3) << 8
	size2 := uint16(frame[4])
	size := size1 + size2
	
	packet.Length = size
	
	return &packet, nil
}

func (l *LLAPPacket) PrettyHeaders() string {
	return fmt.Sprintf("from %s to %s, type %s, %d data bytes", 
		llapAddressString(l.Dst), llapAddressString(l.Src), llapTypeString(l.LLAPType), l.Length)
}
