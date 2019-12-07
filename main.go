package main

import (
	lt "github.com/cheesestraws/appletalk/lib/localtalk"
)

func main() {
	p := lt.NewPort(&lt.LToUDPListener{})
	p.Start()
	p.AcquireAddress()
	
	select { }	
}