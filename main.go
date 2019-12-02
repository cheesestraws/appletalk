package main

import (
	lt "github.com/cheesestraws/appletalk/lib/localtalk"
)

func main() {
	l := lt.LToUDPListener{}
	l.Start()	
}