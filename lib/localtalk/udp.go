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
	for {
		buf := make([]byte, 700, 700)

		count, _, err := conn.ReadFrom(buf)
		
		if count == 0 {
			log.Printf("Empty packet")
		}
		
		if count > 0 {
			log.Printf("Got packet: %d bytes", count)
			frame, err := DecodeLLAPPacket(buf[0:count])
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
		
		if err != nil {
			log.Printf("Got err: %v", err)
		}
	}
	
	return nil
}