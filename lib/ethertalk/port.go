package ethertalk

import (
	"log"
	"bytes"
)

var appletalkSNAPProtocol = []byte{0x08, 0x00, 0x07, 0x80, 0x9B}
var aarpSNAPProtocol = []byte{0x00, 0x00, 0x00, 0x80, 0xF3}

type Port struct {
	io     Listener
	sendC  chan<- []byte
	errorC <-chan error
}

func NewPort(io Listener) *Port {
	return &Port{
		io: io,
	}
}

func (p *Port) Start() error {
	err := p.io.Start()
	if err != nil {
		return err
	}
	
	recvC, sendC, errorC := p.io.Channels()
	p.sendC = sendC
	p.errorC = errorC
	
	// start error logging channel
	go func() {
		for err := range errorC {
			log.Printf("err: %v", err)
		}
	}()
	
	// and processing received frames
	go func() {
		for f := range recvC {
			log.Printf("port: recvd %v bytes", len(f))
			if isAppleTalkPhase2(f) {
				// do something
			}
		}
	}()
	
	return nil
}

func isAppleTalkPhase2(frame []byte) bool {
	// An EtherTalk Phase 2 frame uses 802.2 LLC/SNAP.
	
	// Do we actually have room for an LLC/SNAP header?
	if len(frame) < 22 {
		log.Printf("frame too short")
		return false
	}
	
	// Do we have a length not an ethertype?
	var length int
	length = (int(frame[12]) << 8) | int(frame[13])
	
	// if the length is <= 1500 it's a length and we're OK (we're in 802.2 land)
	// but if it's > 1500 it's an ethertype and we have an Ethernet II frame.
	// An AppleTalk ethertype means Phase 1 (don't ask.) and we're not doing
	// phase 1, if only because I don't have any docs for it.
	if length > 1500 {
		log.Printf("got ethernet ii frame")
		return false
	}
	
	// Check that the DSAP, SSAP and control field are correct
	if !bytes.Equal(frame[14:17], []byte{0xAA, 0xAA, 0x3}) {
		log.Printf("no SNAP LLC header")
		return false
	}
	
	// Do we have a valid SNAP protocol discriminator?
	if bytes.Equal(frame[17:22], appletalkSNAPProtocol) {
		log.Printf("got appletalk")
		return true
	}
	if bytes.Equal(frame[17:22], aarpSNAPProtocol) {
		log.Printf("got aarp")
		return true
	}
	
	return false
}