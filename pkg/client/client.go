package client

import (
	pb "logstream/pkg/proto"
	"time"

	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

// ResponseWriter is used to write all responses from the upstream server for
// subsequent processing.
type ResponseWriter interface {
	Write(*pb.Response)
}

// NewClient is the initializer for Client.
func NewClient(raddr string, readFx time.Duration, w ResponseWriter) *Client {
	return &Client{
		raddr: raddr,
		fx:    readFx,
		w:     w,
	}
}

// Client is a consumer of a ReaderClient stream
type Client struct {
	raddr string
	fx    time.Duration
	w     ResponseWriter
}

// Run consumes from a stream until the context expires.
func (c *Client) Run(ctx context.Context, reqID string) {
	conn, err := grpc.Dial(
		c.raddr,
		grpc.WithInsecure(),
		grpc.WithStatsHandler(logHandler{}),
	)
	if err != nil {
		// NOTE grpc.Dial returns an error only when the provided options are
		// together invalid. It will not return an error if the server is down.
		log.Fatal("failed to create gRPC connection ", c.raddr, err)
	}

	client := pb.NewReaderClient(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.read(ctx, client, reqID)
		}
	}
}

func (c *Client) read(ctx context.Context, client pb.ReaderClient, id string) {
	stream, err := client.CreateStream(
		ctx,
		&pb.Request{AppId: id},
	)
	if err != nil {
		log.Printf("failed to create stream: %s", err)
		// Sleep to avoid hammering the stream with multiple connection
		// attempts.
		time.Sleep(c.fx)
		return
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Printf("stream has closed: %s", err)
			return
		}
		log.Printf("client received response: %s", resp)
		c.w.Write(resp)

		time.Sleep(c.fx)
	}
}

type logHandler struct{}

func (logHandler) TagRPC(ctx context.Context, i *stats.RPCTagInfo) context.Context {
	log.Printf("TagRPC", i)
	return ctx
}

func (logHandler) HandleRPC(ctx context.Context, r stats.RPCStats) {
	log.Printf("HandleRPC", r)
}

func (logHandler) TagConn(ctx context.Context, c *stats.ConnTagInfo) context.Context {
	log.Printf("TagConn", c)
	return ctx
}

func (logHandler) HandleConn(ctx context.Context, c stats.ConnStats) {
	log.Printf("HandleConn", c)
}
