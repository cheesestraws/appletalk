package localtalk

import (
	"sync"
)

// A CallbackChain value stores an unordered set of callback functions to be
// called when some kind of packet reception event occurs.
type CallbackChain struct {
	l sync.RWMutex

	m map[*func(l *LLAPPacket)]struct{} // well this is a hack
}

// Add adds a function to the callback chain.  Callbacks must NOT mutate the
// packet they are given.
func (c *CallbackChain) Add(p *func(l *LLAPPacket)) {
	c.l.Lock()
	defer c.l.Unlock()

	if c.m == nil {
		c.m = make(map[*func(l *LLAPPacket)]struct{})
	}

	c.m[p] = struct{}{}
}

// Remove removes a function from the callback chain.  Note that the function
// pointer passed to remove must be EXACTLY the pointer that was passed to
// Add.  It is not sufficient for it to be a different pointer to the same
// function!  This is because there is no function equality in Go.
func (c *CallbackChain) Remove(p *func(l *LLAPPacket)) {
	c.l.Lock()
	defer c.l.Unlock()

	if c.m == nil {
		return
	}

	delete(c.m, p)
}

// Run runs every callback in the chain with the given packet.
func (c *CallbackChain) Run(p *LLAPPacket) {
	c.l.RLock()
	defer c.l.RUnlock()

	for callback := range c.m {
		if callback != nil {
			(*callback)(p)
		}
	}
}
