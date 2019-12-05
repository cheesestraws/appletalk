package localtalk

// A Listener represents the interface to the hardware or the emulated hardware.  It
// provides two channels for sending and receiving LLAP packets and a means of starting
// the packet listener.
type Listener interface {
	Start() error
	Channels() (<-chan []byte, chan<- []byte, <-chan error)
}
