package main

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

func NewScanner(p *lt.Port) *Scanner {
	return &Scanner{
		p: p,
		recvTimes: make(map[uint8]time.Time),
	}
}

func (s *Scanner) Scan() {
	// Add a callback to note when/whether we got a response, and defer
	// removing it again
	callback := func (p *lt.LLAPPacket) {
		if p.LLAPType == 0x82 {
			log.Printf("pong from %d", p.Src)
			// Have we already got a response?
			_, ok := s.recvTimes[p.Src]
			if !ok {
				s.recvTimes[p.Src] = time.Now()
			}
		}
	}
	s.p.LLAPControlCallbacks.Add(&callback)
	defer s.p.LLAPControlCallbacks.Remove(&callback)
	
	// now fire off some ENQs
	var addr uint8
	for addr = 1; addr < 255; addr++ {
		for i := 0; i < numberOfEnqs; i++ {
			s.p.SendRaw([]byte{addr, addr, 0x81})
			time.Sleep(200 * time.Microsecond)
		}
	}
	
	time.Sleep(1 * time.Second)
}
