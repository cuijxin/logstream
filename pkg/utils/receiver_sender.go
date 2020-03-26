package utils

import (
	pb "logstream/pkg/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"log"
)

func NewReceiverSender(raddr string) *ReceiverSender {
	return &ReceiverSender{
		raddr: raddr,
	}
}

type ReceiverSender struct {
	raddr     string
	connector pb.ReaderClient
}

func (u *ReceiverSender) CreateStream(req *pb.Request, srv pb.Reader_CreateStreamServer) error {
	log.Printf("creating stream to server with id %s", req.AppId)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	client, err := u.connector.CreateStream(ctx, req)
	if err != nil {
		return err
	}
	for {
		resp, err := client.Recv()
		if err != nil {
			return err
		}

		if err := srv.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

func (u *ReceiverSender) Connect() error {
	conn, err := grpc.Dial(u.raddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	u.connector = pb.NewReaderClient(conn)

	return nil
}
