package main

import (
	et "github.com/cheesestraws/appletalk/lib/ethertalk"
)

func main() {
	p := et.NewPort(&et.PCAPListener{
		Interface: "tap0",
	})
	
	p.Start()
	
	select {}
}
