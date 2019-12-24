package ethertalk

import (
	"github.com/google/gopacket/pcap"
)

type PCAPListener struct {
	Interface  string
	
	recvC chan []byte
	sendC chan []byte
	errC  chan error
}

func (l *PCAPListener) init() {
	if l.recvC == nil {
		l.recvC = make(chan []byte)
	}
	if l.sendC == nil {
		l.sendC = make(chan []byte)
	}
	if l.errC == nil {
		l.errC = make(chan error)
	}
}

func (l *PCAPListener) Channels() (<-chan []byte, chan<- []byte, <-chan error) {
	return l.recvC, l.sendC, l.errC
}

func (l *PCAPListener) Start() error {
	l.init()

	// open the pcap handles.
	readHandle, err := pcap.OpenLive(l.Interface, 65535, false, pcap.BlockForever)
	if err != nil {
		return err
	}
	err = readHandle.SetBPFFilter("atalk or aarp")
	if err != nil {
		readHandle.Close()
		return err
	}
	
	writeHandle, err := pcap.OpenLive(l.Interface, 0, false, pcap.BlockForever)
	if err != nil {
		readHandle.Close()
		return err
	}
	
	go l.runRecv(readHandle)
	go l.runSend(writeHandle)
	
	return nil
}

func (l *PCAPListener) runRecv(handle *pcap.Handle) {
	for {
		data, _, err := handle.ReadPacketData()
		if err != nil {
			l.errC <- err
			continue
		}
		l.recvC <- data
	}
}

func (l *PCAPListener) runSend(handle *pcap.Handle) {
	for bb := range l.sendC {
		err := handle.WritePacketData(bb)
		if err != nil {
			l.errC <- err
			continue
		}
	}
}