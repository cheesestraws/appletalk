package localtalk

import (
	"log"
	"net"
)

// LToUDPListener implements LocalTalk over UDP multicast
type LToUDPListener struct {
	recvC chan []byte
	sendC chan []byte
	errC  chan error
}

var _ Listener = &LToUDPListener{}

var multicastGroup = &net.UDPAddr{
	IP:   net.ParseIP("224.0.0.3"),
	Port: 1954,
}

func (l *LToUDPListener) init() {
	if l.recvC == nil {
		l.recvC = make(chan []byte)
	}
	if l.sendC == nil {
		l.sendC = make(chan []byte)
	}
	if l.errC == nil {
		l.errC = make(chan error)
	}
}

func (l *LToUDPListener) debugRecv() {
	for packet := range l.recvC {
		log.Printf("Got packet: %d bytes", len(packet))
		frame, err := DecodeLLAPPacket(packet)
		if err != nil {
			log.Printf("    err: %v", err)
		}
		log.Printf("    %s", frame.PrettyHeaders())

		// For debugging: Do we have a DDP packet?
		if frame.LLAPType == 1 || frame.LLAPType == 2 {
			ddp, err := frame.DDP()
			if err != nil {
				log.Printf("    %v", err)
			}
			log.Printf("    %s", ddp.PrettyHeaders())
		}
	}
}

func (l *LToUDPListener) Channels() (<-chan []byte, chan<- []byte, <-chan error) {
	return l.recvC, l.sendC, l.errC
}

func (l *LToUDPListener) Start() error {
	l.init()

//	go l.debugRecv()

	// Set up the multicast listener
	conn, err := net.ListenMulticastUDP("udp", nil, multicastGroup)
	if err != nil {
		return err
	}

	go l.runRecv(conn)
	go l.runSend(conn)

	return nil
}

func (l *LToUDPListener) runRecv(conn *net.UDPConn) {
	// Now enter the packet run loop
	for {
		buf := make([]byte, 700, 700)

		count, _, err := conn.ReadFrom(buf)
		if count > 0 {
			l.recvC <- buf[0:count]
		}

		if err != nil {
			log.Printf("Got err: %v", err)
			l.errC <- err
		}
	}
}

func (l *LToUDPListener) runSend(conn *net.UDPConn) {
	for packet := range l.sendC {
		log.Printf("Sending packet: %d bytes", len(packet))
		frame, err := DecodeLLAPPacket(packet)
		if err != nil {
			log.Printf("    err: %v", err)
		}
		log.Printf("    %s", frame.PrettyHeaders())

		// For debugging: Do we have a DDP packet?
		if frame.LLAPType == 1 || frame.LLAPType == 2 {
			ddp, err := frame.DDP()
			if err != nil {
				log.Printf("    %v", err)
			}
			log.Printf("    %s", ddp.PrettyHeaders())
		}
		
		_, err = conn.WriteToUDP(packet, multicastGroup)
	}
}
