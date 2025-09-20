/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/emaforlin/messagebus/internal/core"
	"github.com/emaforlin/messagebus/internal/server"
	pb "github.com/emaforlin/messagebus/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessageBusServiceServer(grpcServer, server.NewGRPCServer(core.NewInMemoryBus()))

	go func() {
		log.Println("gRPC server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting down server...")
	grpcServer.GracefulStop()
}
