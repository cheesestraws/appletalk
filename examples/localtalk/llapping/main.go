package main

/* llapping is a utility which "pings" a localtalk address or a range thereof.
   To do so, it uses low-level LLAP ENQs and ACKs, which are used in address
   acquisition.  We send a bunch of ENQs for the address, and if it replies
   with an ACK, we know it's up.
   
   This is the localtalk equivalent of an ARP ping on an ethernet network. */
   
import (
	"log"
	"os"
	"strconv"

	lt "github.com/cheesestraws/appletalk/lib/localtalk"
)

const numberOfEnqs = 20

func printUsage() {
	pn := os.Args[0]
	log.Printf("usage: %s -a", pn)
	log.Printf("       Scan the LToUDP network and report on all up nodes")
	log.Printf("")
	log.Printf("       %s [node]", pn)
	log.Printf("       Pings the given node with LLAP ENQs")
}

func main() {
	log.SetFlags(0)

	p := lt.NewPort(&lt.LToUDPListener{})
	p.Start()
	
	// We do not acquire an address, because we're only mucking about
	// with control traffic.
	
	if len(os.Args) == 2 && os.Args[1] == "-a" {
		s := NewScanner(p)
		s.Scan()
		s.PrintResults()
		
		return
	}
	
	if len(os.Args) == 2 {
		node, err := strconv.Atoi(os.Args[1])
		if err == nil && node > 0 && node < 255 {
			for {
				pp := NewPinger(p, uint8(node))
				resp, when := pp.Ping()
				if resp {
					log.Printf("lapACK from %d in %v", node, when)
				} else {
					log.Printf("Request timed out")
				}
			}
			return
		}
	}
	
	printUsage()
}