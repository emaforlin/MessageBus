package client

import (
	"context"
	"io"
	"log"

	pb "github.com/emaforlin/messagebus/proto/messagebus/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client pb.MessageBusServiceClient
}

func NewGRPCClient(addr string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GRPCClient{
		conn:   conn,
		client: pb.NewMessageBusServiceClient(conn),
	}, nil
}

// Publish send a message to a topic
func (c *GRPCClient) Publish(topic, message string) error {
	_, err := c.client.Publish(context.Background(), &pb.PublishRequest{
		Topic: topic,
		Msg:   message,
	})
	return err
}

// Subscribe to a topic and receive a stream of messages
func (c *GRPCClient) Subscribe(topic string) {
	stream, err := c.client.Subscribe(context.Background(), &pb.SubscribeRequest{Topic: topic})
	if err != nil {
		log.Fatal("Subscribe error:", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("Stream closed by server")
			break
		}

		if err != nil {
			log.Fatal("Stream error:", err)
		}

		log.Printf("[gRPC] [%s] %s\n", msg.Topic, msg.Msg)
	}
}

func (c *GRPCClient) Close() {
	_ = c.conn.Close()
}
