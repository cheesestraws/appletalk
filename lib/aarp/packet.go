package aarp

import (
	"fmt"
	"errors"
	"encoding/binary"
)

// Packet encodes the information in an AARP packet.
type Packet struct {
	HardwareType HardwareType
	ProtocolType ProtocolType
	HardwareAddressLength byte
	ProtocolAddressLength byte
	
	Function Function
	HWAddr1 string // Addresses as strings rather than []byte 
	               // so we can use them as hash keys
	ProtoAddr1 string
	HWAddr2 string
	ProtoAddr2 string
}

var ErrShortPacket = errors.New("malformed packet: too short")
var ErrShortAddressPacket = errors.New("malformed packet: too short for addresses (but header ok)")

// Decode turns an AARP packet from its on-the-wire representation into a
// struct.  For the sources of the offsets, please consult /Inside Appletalk/
// 2nd Ed., p. 2-12.
func Decode(packet []byte) (*Packet, error) {
	// if we have less than 8 bytes, something's gone terribly wrong
	if len(packet) < 8 {
		return nil, ErrShortPacket
	}
	
	p := &Packet{
	}
	
	p.HardwareType = HardwareType(binary.BigEndian.Uint16(packet))
	p.ProtocolType = ProtocolType(binary.BigEndian.Uint16(packet[2:]))
	p.HardwareAddressLength = packet[4]
	p.ProtocolAddressLength = packet[5]
	p.Function = Function(binary.BigEndian.Uint16(packet[6:]))
	
	// Now, we know how big the rest of the packet should be, because the
	// rest of the packet consists of two hardware addresses and two protocol
	// addresses.
	
	if len(packet) < 8 + (2 * int(p.HardwareAddressLength)) + (2 * int(p.ProtocolAddressLength)) {
		return nil, ErrShortAddressPacket
	}
	
	nextAddress := packet[8:]
	p.HWAddr1 = string(nextAddress[0:p.HardwareAddressLength])
	nextAddress = nextAddress[p.HardwareAddressLength:]
	
	p.ProtoAddr1 = string(nextAddress[0:p.ProtocolAddressLength])
	nextAddress = nextAddress[p.ProtocolAddressLength:]
	
	p.HWAddr2 = string(nextAddress[0:p.HardwareAddressLength])
	nextAddress = nextAddress[p.HardwareAddressLength:]
	
	p.ProtoAddr2 = string(nextAddress[0:p.ProtocolAddressLength])
	nextAddress = nextAddress[p.ProtocolAddressLength:]
	
	return p, nil
}

func (p *Packet) String() string {
	prefix := fmt.Sprintf("aarp: %v %v/%v: ", p.Function, p.HardwareType, p.ProtocolType)
	
	return prefix
}