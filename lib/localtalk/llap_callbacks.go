package localtalk

import (
	"sync"
)

type PacketCallback func(l *LLAPPacket)

type CallbackChain struct {
	l sync.RWMutex
	
	m map[*PacketCallback]struct{} // well this is a hack
}

func (c *CallbackChain) AddCallback(p *PacketCallback) {
	c.l.Lock()
	defer c.l.Unlock()
	
	if c.m == nil {
		c.m = make(map[*PacketCallback]struct{})
	}
	
	c.m[p] = struct{}{}
}

func (c *CallbackChain) RemoveCallback(p *PacketCallback) {
	c.l.Lock()
	defer c.l.Unlock()
	
	if c.m == nil {
		return
	}
	
	delete(c.m, p)
}

func (c *CallbackChain) Run(p *LLAPPacket) {
	c.l.RLock()
	defer c.l.RUnlock()
	
	for callback := range c.m {
		if callback != nil {
			(*callback)(p)
		}
	}
}