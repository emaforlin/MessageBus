package server

import (
	"context"
	"fmt"

	"github.com/emaforlin/messagebus/internal/core"
	pb "github.com/emaforlin/messagebus/proto/messagebus/v1"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedMessageBusServiceServer
	bus *core.InMemoryBus
}

func NewGRPCServer(bus *core.InMemoryBus) *GRPCServer {
	return &GRPCServer{bus: bus}
}

func (s *GRPCServer) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	var success bool = false
	topic := req.GetTopic()
	msg := req.GetMsg()

	err := s.bus.Publish(ctx, topic, msg)
	if err == nil {
		success = true
	}

	return &pb.PublishResponse{
		Success: success,
	}, err
}

func (s *GRPCServer) Subscribe(req *pb.SubscribeRequest, stream grpc.ServerStreamingServer[pb.SubscribeResponse]) error {
	topic := req.GetTopic()

	done := make(chan struct{})

	id, err := s.bus.Subscribe(topic, func(msg string) error {
		select {
		case <-done:
			return fmt.Errorf("client disconnected")
		default:
			return stream.Send(&pb.SubscribeResponse{
				Topic: topic,
				Msg:   msg,
			})
		}
	})
	if err != nil {
		return err
	}

	<-stream.Context().Done()
	close(done)

	if err := s.bus.Unsubscribe(topic, id); err != nil {
		fmt.Printf("Error unsubscribing: %v\n", err)
	}

	return nil
}
