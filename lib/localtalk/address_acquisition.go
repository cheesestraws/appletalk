package localtalk

/* This file implements address acquisition on a LocalTalk link.  The basic 
   approach here is that we pick a random candidate address, send out a load
   of packets saying "I'm claiming this address" (ENQ packets), and if nobody
   objects (via an ACK) then we use that address.  If anyone does object, then
   we sigh deeply, pick another random address, reset the retry counter, and
   try again.  This continues until we find an address nobody on the network
   objects to or the user gets bored.

   The code consists of two basic elements: handlers for LLAP ACK and ENQ 
   packets (which other computers use to notify us of their own address choices)
   and an AcquireAddress() function that sends out ENQ packets until we find
   an address that nobody else on the network objects to and have retried
   sending our ENQ enough times that we can be reasonably confident we're
   safe.
   
   These two elements are linked together by an addressAcqState, which keeps 
   track of what our current candidate address is and how many times we've
   tried so far.  An ACK for our candidate address generates a new address and
   resets the retry counter.
   
   See also: Inside Appletalk, Second Edition, p. 1-4.
*/

import (
	"sync"
	"math/rand"
	"time"
	"log"
)

var randomiser = rand.New(rand.NewSource(time.Now().UnixNano()))

type addressAcqState struct {
	l sync.Mutex

	discovering bool
	server bool
	addressCandidate uint8
	retriesRemaining int
}

// reset resets the address acquisition state to what it ought to be at the
// start of the address acquisition process.  We need a random address to
// try out and a maximum number of retries.
func (a *addressAcqState) reset() {
	a.l.Lock()
	defer a.l.Unlock()
	
	// pick a random address
	a.addressCandidate = uint8(randomiser.Intn(128))
	a.retriesRemaining = 30
	
	if a.server {
		a.addressCandidate += 128
		a.retriesRemaining *= 5
	}
}

// tellMeWhatToDoNext unwraps the addressAcqState for the use of the main
// address acquisition loop and decrements the retry count.
func (a *addressAcqState) tellMeWhatToDoNext() (uint8, int) {
	a.l.Lock()
	defer a.l.Unlock()
	
	candidate := a.addressCandidate
	retriesRemaining := a.retriesRemaining
	a.retriesRemaining--
	
	return candidate, retriesRemaining
}

// An ACK is someone on the network going "Oi! That's my address!".  If it's
// not our address we don't mind, but if it's ours and we're doing address
// acquisition then we need to try again.  If we're not doing address 
// acquisition and we get an ACK then we ignore it and hope the problem
// goes away, because I don't know what to do in that case.
func (p *Port) handleACK(packet *LLAPPacket) {
	if p.iHaveAnAddress {
		// if I already have an address, none of this matters.
		return
	}
	
	p.addressAcqState.l.Lock()
	candidate := p.addressAcqState.addressCandidate
	p.addressAcqState.l.Unlock()
	
	// If the ack isn't for our candidate address then ignore it
	if packet.Src != candidate {
		return
	}
	
	// it is.  oh well, let's try again
	p.addressAcqState.reset()
}

// An ENQ is someone on the network broadcasting that it intends to claim an
// address if nobody objects.  If we object, we should respond with an ACK
// message
func (p *Port) handleENQ(packet *LLAPPacket) {
	// This is a host enquiring whether its address is actually unique or not.
	if !p.iHaveAnAddress {
		// If I have no address, then ignore this packet
		return
	}
	if packet.Src == p.address {
		// Whoops!  This is my address!
		// Send an acknowledgement that I already own this address, so the other
		// will have to change its tune
		log.Printf("Detected address collision; looking sternly at other node")
		p.SendRaw([]byte{p.address, p.address, lapACK})
	}
}

// AcquireAddress runs the LocalTalk node discovery process to find a unique
// address for this localtalk port on its own network.  When it finds one,
// it sets the address on the port.  Blocks until address acquisition is
// complete.
func (p *Port) AcquireAddress() {
	log.Printf("acquiring address...")
	p.addressAcqState.reset()
	
	var candidate uint8
	var retriesRemaining int
	for {
		candidate, retriesRemaining = p.addressAcqState.tellMeWhatToDoNext()
		
		// Have we run out of retries?
		if retriesRemaining == 0 {
			break
		}
		
		// Transmit an ENQ
		p.SendRaw([]byte{candidate, candidate, lapENQ})
		time.Sleep(200 * time.Microsecond)
	}
	
	p.SetAddress(candidate)
	
	log.Printf("Got address: %v", p.address)
}