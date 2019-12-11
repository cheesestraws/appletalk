package main

/* llapping is a utility which "pings" a localtalk address or a range thereof.
   To do so, it uses low-level LLAP ENQs and ACKs, which are used in address
   acquisition.  We send a bunch of ENQs for the address, and if it replies
   with an ACK, we know it's up.
   
   This is the localtalk equivalent of an ARP ping on an ethernet network. */
   
import (
	"sync"
	"time"
	"log"
	"os"
	"strconv"

	lt "github.com/cheesestraws/appletalk/lib/localtalk"
)

const numberOfEnqs = 20

type Pinger struct {
	l sync.RWMutex
	
	p *lt.Port
	
	address uint8
	iGotAResponse bool
	when time.Time
}

func NewPinger(p *lt.Port, target uint8) *Pinger {
	return &Pinger{
		p: p,
		address: target,
	}
}

func (pp *Pinger) Response() (bool, time.Time) {
	pp.l.RLock()
	defer pp.l.RUnlock()
	return pp.iGotAResponse, pp.when
}

func (pp *Pinger) MarkAsResponded() {
	pp.l.Lock()
	defer pp.l.Unlock()
	
	if pp.iGotAResponse {
		return // only pay attention to the first ACK
	}
	
	pp.iGotAResponse = true
	pp.when = time.Now()
}

func (pp *Pinger) Ping() (bool, time.Duration) {
	// get the target address
	pp.l.RLock()
	target := pp.address
	pp.l.RUnlock()

	// Add a callback to note when/whether we got a response, and defer
	// removing it again
	callback := func (p *lt.LLAPPacket) {
		if p.LLAPType == 0x82 && p.Src == target {
			pp.MarkAsResponded()
		}
	}
	pp.p.LLAPControlCallbacks.Add(&callback)
	defer pp.p.LLAPControlCallbacks.Remove(&callback)
	
	// Now fire off some ENQs
	startTime := time.Now()
	for i := 0; i < numberOfEnqs; i++ {
		pp.p.SendRaw([]byte{target, target, 0x81})
		time.Sleep(200 * time.Microsecond)
	}
	
	time.Sleep(1 * time.Second)
	
	responded, when := pp.Response()
	return responded, when.Sub(startTime)
}

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