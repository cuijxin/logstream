package utils

import (
	pb "logstream/pkg/proto"
	"time"

	"log"
)

func NewReaderService(writeFx time.Duration) *ReaderService {
	return &ReaderService{
		writeFx: writeFx,
	}
}

type ReaderService struct {
	writeFx time.Duration
}

func (s *ReaderService) CreateStream(req *pb.Request, srv pb.Reader_CreateStreamServer) error {
	var err error
	var msgID int64 = 1
	errChan := make(chan error, 1)

	buf := NewBuffer(10)
	go s.drainBuffer(buf, srv, errChan)

	for range time.Tick(s.writeFx) {
		err = s.checkError(errChan)
		if err != nil {
			break
		}

		resp := &pb.Response{
			AppId: req.AppId,
			Epoch: time.Now().Unix(),
			Id:    msgID,
		}

		buf.Write(resp)

		msgID++
	}

	return err
}

func (s *ReaderService) checkError(errCh <-chan error) error {
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

type Sender interface {
	Send(*pb.Response) error
}

func (s *ReaderService) drainBuffer(b *Buffer, sdr Sender, errChan chan<- error) {
	for {
		resp := b.Read()
		err := sdr.Send(resp)
		if err != nil {
			errChan <- err
			return
		}
		log.Printf("sent resp: %v", resp)
	}
}
