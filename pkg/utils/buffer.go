package utils

import (
	pb "logstream/pkg/proto"

	"log"
)

func NewBuffer(capacity int) *Buffer {
	return &Buffer{
		data: make(chan *pb.Response, capacity),
	}
}

type Buffer struct {
	data chan *pb.Response
}

func (b *Buffer) Write(r *pb.Response) {
	select {
	case b.data <- r:
	default:
		<-b.data
		b.data <- r
	}

	log.Printf("wrote to buffer: %v", r)
}

func (b *Buffer) Read() *pb.Response {
	return <-b.data
}
