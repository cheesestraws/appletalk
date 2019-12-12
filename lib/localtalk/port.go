package localtalk

import (
	"log"
	"sync"
)

type Port struct {
	io     Listener
	sendC  chan<- []byte
	errorC <-chan error

	// Address
	al             sync.RWMutex
	iHaveAnAddress bool
	address        uint8

	// Address discovery state
	addressAcqState addressAcqState

	LLAPControlCallbacks CallbackChain
}

// NewPort creates a LocalTalk port using the given listener.  It does not
// start the port.
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

// Start enables the port and its underlying listener, starts listening for 
// traffic on it, and enables sending to it.
func (p *Port) Start() error {
	err := p.io.Start()
	if err != nil {
		return err
	}

	recvC, sendC, errorC := p.io.Channels()

	// receive channel
	go func() {
		for packet := range recvC {
			frame, err := DecodeLLAPPacket(packet)
			if err != nil {
				log.Printf("    err: %v", err)
			}

			if frame.LLAPType >= LAPLowestControlPacketType {
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

func (p *Port) SendLLAP(packet LLAPPacket) {
	p.sendC <- packet.EncodeBytes()
}

// handle an LLAP control packet
func (p *Port) handleLLAPControlPacket(l *LLAPPacket) {
	p.LLAPControlCallbacks.Run(l)

	if l.LLAPType == LAPACK {
		p.handleACK(l)
		return
	}

	if l.LLAPType == LAPENQ {
		p.handleENQ(l)
		return
	}
}
