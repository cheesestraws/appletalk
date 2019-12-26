package aarp

/* The types.go file contains auxiliary types and constants to make life easier
   while working with AARP packets. */

import (
	"fmt"
)

// HardwareType encodes what kind of hardware address is being asked for.
// See /Inside Appletalk/, 2nd ed., p. 2-12.s
type HardwareType uint16

func (h HardwareType) String() string {
	if h == EthernetHardware {
		return "EthernetHardware"
	}
	if h == TokenRingHardware {
		return "TokenRingHardware"
	}
	return fmt.Sprintf("UnknownHardware %d", uint16(h))
}

const EthernetHardware HardwareType = 1
const TokenRingHardware HardwareType = 2

// ProtocolType encodes what kind of layer 3 ("protocol") address is being 
// asked for.  See /Inside Appletalk/, 2nd ed., p. 2-12.s
type ProtocolType uint16

const AppleTalkProtocol ProtocolType = 0x809B

func (p ProtocolType) String() string {
	if p == AppleTalkProtocol {
		return "AppleTalkProtocol"
	}
	return fmt.Sprintf("UnknownProtocol %d", uint16(p))
}