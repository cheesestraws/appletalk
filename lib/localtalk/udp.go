package localtalk

import (
	"log"
	"net"
	"os"
)

// LToUDPListener implements LocalTalk over UDP multicast
type LToUDPListener struct {
	recvC chan []byte
	sendC chan []byte
	errC  chan error
}

var _ Listener = &LToUDPListener{}

var multicastGroup = &net.UDPAddr{
	IP:   net.ParseIP("239.192.76.84"),
	Port: 1954,
}
var broadcast = &net.UDPAddr{
	IP:   net.ParseIP("255.255.255.255"),
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

func packetIsOneOfMine(recvaddr *net.UDPAddr, buf []byte) bool {
	// Does the pid in the packet match mine?
	pid := os.Getpid()
	for i := 0; i < 4; i++ {
		if buf[i] != uint8((pid>>(3-i)*8)&0xff) {
			return false
		}
	}

	// It does!  Is the IP address one of mine?
	// todo: better error handling. or any error handling.
	intfs, _ := net.Interfaces()
	for _, intf := range intfs {
		addrs, _ := intf.Addrs()
		for _, addr := range addrs {
			ipaddr, ok := addr.(*net.IPAddr)
			if ok {
				if ipaddr.IP.Equal(recvaddr.IP) {
					return true
				}
			}
		}
	}

	return false
}

func pidEmbeddedPacket(packet []byte) []byte {
	buf := make([]byte, len(packet)+4, len(packet)+4)
	copy(buf[4:], packet)

	pid := os.Getpid()
	for i := 0; i < 4; i++ {
		buf[i] = uint8((pid >> (3 - i) * 8) & 0xff)
	}
	return buf
}

func (l *LToUDPListener) runRecv(conn *net.UDPConn) {
	// Now enter the packet run loop
	for {
		buf := make([]byte, 700, 700)

		count, addr, err := conn.ReadFromUDP(buf)

		if count > 4 && !packetIsOneOfMine(addr, buf) {
			l.recvC <- buf[4:count]
		}

		if err != nil {
			log.Printf("Got err: %v", err)
			l.errC <- err
		}
	}
}

func (l *LToUDPListener) runSend(conn *net.UDPConn) {
	for packet := range l.sendC {
		/*
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
		*/
		conn.WriteToUDP(pidEmbeddedPacket(packet), multicastGroup)
	}
}
