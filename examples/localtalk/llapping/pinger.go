package main

/* pinger.go contains code to ping a localtalk node */

import (
	"sync"
	"time"

	lt "github.com/cheesestraws/appletalk/lib/localtalk"
)

// A Pinger pings a LocalTalk node
type Pinger struct {
	l sync.RWMutex
	
	p *lt.Port
	
	address uint8
	iGotAResponse bool
	when time.Time
}

// NewPinger returns a new Pinger pointing out of the given LocalTalk port,
// aimed at the target given.
func NewPinger(p *lt.Port, target uint8) *Pinger {
	return &Pinger{
		p: p,
		address: target,
	}
}

// Response returns whether a response was forthcoming, and if so, how long
// it took to get here
func (pp *Pinger) Response() (bool, time.Time) {
	pp.l.RLock()
	defer pp.l.RUnlock()
	return pp.iGotAResponse, pp.when
}

// markAsResponded is to be called when we get an ACK from the node
func (pp *Pinger) markAsResponded() {
	pp.l.Lock()
	defer pp.l.Unlock()
	
	if pp.iGotAResponse {
		return // only pay attention to the first ACK
	}
	
	pp.iGotAResponse = true
	pp.when = time.Now()
}

// Ping pings the target and returns whether it responded and how long it took
// to do so.
func (pp *Pinger) Ping() (bool, time.Duration) {
	// get the target address
	pp.l.RLock()
	target := pp.address
	pp.l.RUnlock()
	
	// Reset the stats
	pp.l.Lock()
	pp.iGotAResponse = false
	pp.l.Unlock()
	
	// Add a callback to note when/whether we got a response, and defer
	// removing it again
	callback := func (p *lt.LLAPPacket) {
		if p.LLAPType == 0x82 && p.Src == target {
			pp.markAsResponded()
		}
	}
	pp.p.LLAPControlCallbacks.Add(&callback)
	defer pp.p.LLAPControlCallbacks.Remove(&callback)
	
	// Now fire off some ENQs
	startTime := time.Now()
	for i := 0; i < numberOfEnqs; i++ {
		pp.p.SendLLAP(lt.LLAPPacket{target, target, 0x81, nil})
		time.Sleep(200 * time.Microsecond)
	}
	
	// Give it a second to respond.
	time.Sleep(1 * time.Second)
	
	responded, when := pp.Response()
	return responded, when.Sub(startTime)
}
