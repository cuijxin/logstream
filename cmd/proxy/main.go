package main

import (
	"flag"
	"logstream/pkg/proxy"
)

func main() {
	laddr := flag.String(
		"laddr",
		":8500",
		"local address on which to listen for incoming connections",
	)
	raddr := flag.String(
		"raddr",
		":8000",
		"remote address to which to connect",
	)
	flag.Parse()

	p := proxy.NewProxy(*laddr, *raddr)
	p.Start()
}
