package main

/* A Scanner scans a LocalTalk network for nodes using LLAP ENQ packets, 
   printing out the nodes it finds. */

import (
	"time"
	"log"

	lt "github.com/cheesestraws/appletalk/lib/localtalk"
)

type Scanner struct {
	p *lt.Port
	
	startTimes [256]time.Time
	recvTimes map[uint8]time.Time
}

// NewScanner returns a scanner that will scan the LocalTalk network that the
// port given is attached to.
func NewScanner(p *lt.Port) *Scanner {
	return &Scanner{
		p: p,
		recvTimes: make(map[uint8]time.Time),
	}
}

// Scan scans a network and notes which nodes are up.  Only call Scan once
// on a scanner.
func (s *Scanner) Scan() {
	// Add a callback to note when/whether we got a response, and defer
	// removing it again
	callback := func (p *lt.LLAPPacket) {
		if p.LLAPType == 0x82 {
			// Have we already got a response from this node?
			_, ok := s.recvTimes[p.Src]
			if !ok {
				// no!  Note when we received the reply.
				s.recvTimes[p.Src] = time.Now()
			}
		}
	}
	s.p.LLAPControlCallbacks.Add(&callback)
	defer s.p.LLAPControlCallbacks.Remove(&callback)
	
	// now fire off some ENQs for the entire namespace.  This feels slightly
	// rude, somehow.
	var addr uint8
	for addr = 1; addr < 255; addr++ {
		// note when we started sending ENQs
		s.startTimes[addr] = time.Now()
		for i := 0; i < numberOfEnqs; i++ {
			s.p.SendLLAP(lt.LLAPPacket{addr, addr, 0x81, nil})
			time.Sleep(200 * time.Microsecond)
		}
	}
	
	// Sleep to wait for any remaining nodes to shout at us
	time.Sleep(1 * time.Second)
}

// PrintResults prints the results from a scan
func (s *Scanner) PrintResults() {
	var addr uint8
	for addr = 1; addr < 255; addr++ {
		endTime, ok := s.recvTimes[addr]
		
		if ok {
			latency := endTime.Sub(s.startTimes[addr])
			log.Printf("Node %d responded in %v", addr, latency)
		}
	}
}
