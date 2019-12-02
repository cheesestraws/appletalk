package localtalk

import (
	"net"
	"log"
)

type LToUDPListener struct {

}

func (l *LToUDPListener) Start() error {
	// Set up the multicast listener
	conn, err := net.ListenMulticastUDP("udp", nil, &net.UDPAddr{
		IP: net.ParseIP("224.0.0.3"),
		Port: 1954,
	})
	
	if err != nil {
		return err
	}
	
	// Now enter the packet run loop
	buf := make([]byte, 700, 700)
	
	for {
		count, _, err := conn.ReadFrom(buf)
		
		if count == 0 {
			log.Printf("Empty packet")
		}
		
		if count > 0 {
			log.Printf("Got packet: %d bytes", count)
			frame, err := DecodeLLAPFrame(buf[0:count])
			if err != nil {
				log.Printf("    err: %v", err)
			}
			log.Printf("    %s", frame.PrettyHeaders())
		}
		
		if err != nil {
			log.Printf("Got err: %v", err)
		}
	}
	
	return nil
}