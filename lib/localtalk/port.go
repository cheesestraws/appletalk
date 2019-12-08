package localtalk

import (
	"log"
	"sync"
)

type Port struct {
	io Listener
	sendC chan<- []byte
	errorC <-chan error
	
	// Address
	al sync.RWMutex
	iHaveAnAddress bool
	address uint8
	
	// Address discovery state
	addressAcqState addressAcqState
}

func NewPort(l Listener) *Port {
	return &Port{
		io: l,
	}
}

func (p *Port) SetAddress(addr uint8) {
	p.al.Lock()
	defer p.al.Unlock()
	
	p.iHaveAnAddress = true
	p.address = addr
}

func (p *Port) Address() (uint8, bool) {
	p.al.RLock()
	defer p.al.RUnlock()
	
	return p.address, p.iHaveAnAddress
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
	
	return nil
}

func (p *Port) SendRaw(packet []byte) {
	p.sendC <- packet
}

// handleLLAPControlPacket 
func (p *Port) handleLLAPControlPacket(l *LLAPPacket) {
	if l.LLAPType == lapACK {
		p.handleACK(l)
		return
	}
	
	if l.LLAPType == lapENQ {
		p.handleENQ(l)
		return
	}
}