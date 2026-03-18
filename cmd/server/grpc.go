package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"os"

	pb "event-streamer/provider"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func RunGrpcServer(
	ctx context.Context,
	addr string,
	errChan chan<- error,
) {
	go func() {
		var opts []grpc.ServerOption

		if os.Getenv("EVENT_STREAMER_USE_CERTS") != "" {
			cert, err := tls.X509KeyPair([]byte(CertPEM), []byte(KeyPEM))
			if err != nil {
				slog.ErrorContext(ctx, "failed to load key pair", "err", err)
				errChan <- err
				return
			}
			slog.InfoContext(ctx, "grpc tls mode enabled")

			creds := credentials.NewServerTLSFromCert(&cert)
			opts = append(opts, grpc.Creds(creds))
		} else {
			slog.InfoContext(ctx, "grpc INSECURE mode enabled")
		}

		lis, err := net.Listen("tcp", addr)
		if err != nil {
			slog.ErrorContext(ctx, "grpc failed to listen", "err", err)
			errChan <- err
			return
		}

		s := grpc.NewServer(opts...)
		pb.RegisterStreamerServer(s, &GrpcServer{})
		reflection.Register(s)

		slog.InfoContext(ctx, "grpc server running", "address", addr)

		if err := s.Serve(lis); err != nil {
			slog.ErrorContext(ctx, "grpc failed to start server", "err", err)
			errChan <- err
		}
	}()
}

type GrpcServer struct {
	pb.UnimplementedStreamerServer
}

func (s *GrpcServer) GetStatus(ctx context.Context, in *pb.DataChunk) (*pb.DataChunk, error) {
	hostname, ok := os.LookupEnv("EVENT_STREAMER_HOSTNAME")
	if !ok {
		hostname = "unknown"
	}

	log.Printf("Unary request received: %s", in.Content)

	return &pb.DataChunk{
		Content: "Handled by: " + hostname + " | Echo: " + in.Content,
	}, nil
}

func (s *GrpcServer) EchoStream(stream pb.Streamer_EchoStreamServer) error {
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
