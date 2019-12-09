package localtalk

import (
	"sync"
)

type CallbackChain struct {
	l sync.RWMutex
	
	m map[*func(l *LLAPPacket)]struct{} // well this is a hack
}

func (c *CallbackChain) Add(p *func(l *LLAPPacket)) {
	c.l.Lock()
	defer c.l.Unlock()
	
	if c.m == nil {
		c.m = make(map[*func(l *LLAPPacket)]struct{})
	}
	
	c.m[p] = struct{}{}
}

func (c *CallbackChain) Remove(p *func(l *LLAPPacket)) {
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