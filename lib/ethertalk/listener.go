package ethertalk

// A Listener represents the interface to the hardware or the emulated hardware.  It
// provides two channels for sending and receiving 802.2 LLC frames and a means of starting
// the packet listener.
type Listener interface {
	Start() error
	Channels() (<-chan []byte, chan<- []byte, <-chan error)
}
