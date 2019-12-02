package localtalk

import (
	"fmt"
	"strconv"
)

func llapAddressString(addr uint8) string {
	if addr == 255 {
		return "broadcast"
	}
	
	return strconv.Itoa(int(addr))
}

func llapTypeString(t uint8) string {
	if t == 1 {
		return "DDP (short header)"
	}
	if t == 2 {
		return "DDP (long header)"
	}
	
	if t >= 128 {
		return fmt.Sprintf("LLAP control packet: %x", t)
	}
	
	return strconv.Itoa(int(t))
}