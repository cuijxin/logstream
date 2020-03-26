package proxy

import (
	pb "logstream/pkg/proto"
	"logstream/pkg/utils"
	"net"

	"log"

	"google.golang.org/grpc"
)

func NewProxy(laddr, raddr string) *proxy {
	return &proxy{
		laddr:          laddr,
		grpcSrv:        grpc.NewServer(),
		upstreamReader: utils.NewReceiverSender(raddr),
	}
}

type proxy struct {
	laddr          string
	grpcSrv        *grpc.Server
	upstreamReader *utils.ReceiverSender
}

func (p *proxy) Start() {
	log.Println("starting proxy on host port", p.laddr)
	p.startUpstreamConnection()

	p.startDownstreamServer()
}

func (p *proxy) startUpstreamConnection() {
	if err := p.upstreamReader.Connect(); err != nil {
		log.Fatal("failed to connect to upstream server ", err)
	}
}

func (p *proxy) startDownstreamServer() {
	pb.RegisterReaderServer(p.grpcSrv, p.upstreamReader)

	lis, err := net.Listen("tcp", p.laddr)
	if err != nil {
		log.Fatal("failed to listen on hostport ", p.laddr)
	}

	err = p.grpcSrv.Serve(lis)
	if err != nil {
		log.Fatal("failed to start server on ", p.laddr)
	}
}
