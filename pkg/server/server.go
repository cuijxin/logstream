package server

import (
	"log"
	"net"
	"time"

	// "github.com/lexkong/log"

	"google.golang.org/grpc"

	pb "logstream/pkg/proto"
	"logstream/pkg/utils"
)

func NewServer(hostPort string, writeFx time.Duration) *Server {
	return &Server{
		grpcSrv:   grpc.NewServer(),
		localAddr: hostPort,
		writeFx:   writeFx,
	}
}

type Server struct {
	grpcSrv   *grpc.Server
	localAddr string
	writeFx   time.Duration
}

func (s *Server) Start() {
	log.Println("starting server on host port", s.localAddr)
	pb.RegisterReaderServer(s.grpcSrv, utils.NewReaderService(s.writeFx))

	lis, err := net.Listen("tcp", s.localAddr)
	if err != nil {
		log.Fatal("failed to listen on hostport", s.localAddr)
	}
	s.grpcSrv.Serve(lis)
}
