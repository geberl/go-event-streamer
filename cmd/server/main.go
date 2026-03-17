package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	pb "event-streamer/provider"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStreamerServer
}

func (s *server) GetStatus(ctx context.Context, in *pb.DataChunk) (*pb.DataChunk, error) {
	hostname, ok := os.LookupEnv("EVENT_STREAMER_HOSTNAME")
	if !ok {
		hostname = "unknown"
	}

	log.Printf("Unary request received: %s", in.Content)

	return &pb.DataChunk{
		Content: "Handled by: " + hostname + " | Echo: " + in.Content,
	}, nil
}

func (s *server) EchoStream(stream pb.Streamer_EchoStreamServer) error {
	log.Println("New stream connection established.")

	hostname, ok := os.LookupEnv("EVENT_STREAMER_HOSTNAME")
	if !ok {
		hostname = "unknown"
	}

	for {
		// Read incoming message from client
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Received: %s", in.Content)

		// Send response back
		reply := fmt.Sprintf("Server Echo from %s: %s", hostname, in.Content)
		if err := stream.Send(&pb.DataChunk{Content: reply}); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterStreamerServer(s, &server{})

	log.Println("Server listening on :50051 (Unary + Streaming) ...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
