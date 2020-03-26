package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"logstream/pkg/client"
	pb "logstream/pkg/proto"
)

func main() {
	raddr := flag.String("raddr", ":8500", "remote address of upstream server")
	id := flag.String("id", "1", "unique ID of client")
	readFx := flag.Duration("freq", 1*time.Second, "frequency of reads")
	flag.Parse()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	client := client.NewClient(*raddr, *readFx, nullWriter{})
	client.Run(ctx, *id)
}

type nullWriter struct{}

func (nullWriter) Write(*pb.Response) {}
