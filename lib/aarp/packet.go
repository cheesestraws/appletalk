package aarp

import (
	"errors"
	"encoding/binary"
)

// Packet encodes the information in an AARP packet.
type Packet struct {
	HardwareType uint16
	ProtocolType uint16
	HardwareAddressLength byte
	ProtocolAddressLength byte
	
	Function uint16
	HWAddr1 string // Addresses as strings rather than []byte 
	               // so we can use them as hash keys
	ProtoAddr1 string
	HWAddr2 string
	ProtoAddr2 string
}

var ErrShortPacket = errors.New("malformed packet: too short")

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
	
	p.HardwareType = binary.BigEndian.Uint16(packet)
	p.ProtocolType = binary.BigEndian.Uint16(packet[2:])
	p.HardwareAddressLength = packet[4]
	p.ProtocolAddressLength = packet[5]
	p.Function = binary.BigEndian.Uint16(packet)
	
	return p, nil
}