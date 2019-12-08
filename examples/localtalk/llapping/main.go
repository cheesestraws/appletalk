package main

/* llapping is a utility which "pings" a localtalk address or a range thereof.
   To do so, it uses low-level LLAP ENQs and ACKs, which are used in address
   acquisition.  We send a bunch of ENQs for the address, and if it replies
   with an ACK, we know it's up.
   
   This is the localtalk equivalent of an ARP ping on an ethernet network. */
   
import (
	"sync"
	"time"

	lt "github.com/cheesestraws/appletalk/lib/localtalk"
)

const numberOfEnqs = 5

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

func (pp *Pinger) Ping() (bool, time.Duration) {
	// get the target address
	pp.l.RLock()
	target := pp.address
	pp.l.RUnlock()

	// add listener here
	// defer removing it
	
	startTime := time.Now()
	for i := 0; i < numberOfEnqs; i++ {
		pp.p.SendRaw([]byte{target, target, 0x81})
		time.Sleep(200 * time.Microsecond)
	}
	
	responded, when := pp.Response()
	return responded, when.Sub(startTime)
}

func main() {

}