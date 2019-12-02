package ddp

import (
	"errors"
	"fmt"
)

type Packet struct {
	ShortHeaders bool

	HopCount uint8
	Length uint16
	Checksum uint16
	DstNet uint16
	SrcNet uint16
	DstNode uint8
	SrcNode uint8
	DstSocket uint8
	SrcSocket uint8
	DDPType uint8
	
	Data []byte
}

var IncompleteHeaders = errors.New("incomplete DDP headers")
var Short error = errors.New("packet data too short for its declared DDP length")

// DecodePacket decodes a DDP packet into a DDPPacket struct.  Note that if you only
// have short DDP headers, the resulting DDP struct is incomplete and requires the src
// and dst nodes filling in from the 
func DecodePacket(p []byte, withShortHeaders bool, validateChecksum bool) (*Packet, error) {
	if withShortHeaders {
		return decodeShortHeaders(p)
	}
	return nil, errors.New("unimplemented")
}

func decodeShortHeaders(p []byte) (*Packet, error) {
	packet := &Packet{
		ShortHeaders: true,
	}
	
	if len(p) < 5 {
		return nil, IncompleteHeaders
	}
	
	size1 := uint16(p[0] & 3) << 8
	size2 := uint16(p[1])
	size := size1 + size2
	
	if len(p) < int(size) {
		return nil, Short
	}
	
	packet.Length = size
	
	packet.DstSocket = p[2]
	packet.SrcSocket = p[3]
	packet.DDPType = p[4]
	packet.Data = p[5:len(p)]
	
	return packet, nil
}

func (p *Packet) PrettyHeaders() string {
	if p.ShortHeaders {
		return fmt.Sprintf("%v -> %v, DDP type %v, %v bytes", p.SrcNode, p.DstNode, 
			p.DDPType, p.Length)
	} else {
		return "unimplemented"
	}
}
