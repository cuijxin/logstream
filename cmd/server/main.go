package main

import (
	"flag"
	"logstream/pkg/server"
	"time"
)

func main() {
	addr := flag.String("addr", ":8000", "address on which to listen")
	writeFx := flag.Duration("freq", 1*time.Second, "frequency of writes")
	flag.Parse()

	s := server.NewServer(*addr, *writeFx)
	s.Start()
}
