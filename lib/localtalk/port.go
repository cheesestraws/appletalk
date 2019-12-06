package localtalk

import (
	"log"
)

type Port struct {
	io Listener
	sendC chan<- []byte
	errorC <-chan error
	
	// Address
	iHaveAnAddress bool
	address uint8
}

func NewPort(l Listener) *Port {
	return &Port{
		io: l,
	}
}

func (p *Port) Start() error {
	err := p.io.Start()
	if err != nil {
		return err
	}
	
	recvC, sendC, errorC := p.io.Channels()
	
	// receive channel
	go func() {
		for packet := range recvC {
			// DEBUG
			log.Printf("Got packet: %d bytes", len(packet))
			frame, err := DecodeLLAPPacket(packet)
			if err != nil {
				log.Printf("    err: %v", err)
			}
			log.Printf("    %s", frame.PrettyHeaders())
			
			/*
			// For debugging: Do we have a DDP packet?
			if frame.LLAPType == 1 || frame.LLAPType == 2 {
				ddp, err := frame.DDP()
				if err != nil {
					log.Printf("    %v", err)
				}
				log.Printf("    %s", ddp.PrettyHeaders())
			}
			*/
			
			if frame.LLAPType >= lapLowestControlPacketType {
				p.handleLLAPControlPacket(frame)
			}
		}
	}()
	
	p.sendC = sendC
	p.errorC = errorC
	
	// for the moment force to have an address of 129
	p.iHaveAnAddress = true
	p.address = 250
	
	return nil
}

func (p *Port) SendRaw(packet []byte) {
	p.sendC <- packet
}

// handleLLAPControlPacket 
func (p *Port) handleLLAPControlPacket(l *LLAPPacket) {
	if l.LLAPType == lapENQ {
		// This is a host enquiring whether its address is actually unique or not.
		if !p.iHaveAnAddress {
			// If I have no address, then ignore this packet
			return
		}
		if l.Src == p.address {
			// Whoops!  This is my address!
			// Send an acknowledgement that I already own this address, so the other
			// will have to change its tune
			log.Printf("Detected address collision; looking sternly at other node")
			p.SendRaw([]byte{p.address, p.address, lapACK})
		}
	}
}