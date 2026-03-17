package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"time"

	pb "event-streamer/provider"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            nil,
	}
	creds := credentials.NewTLS(tlsConfig)

	conn, err := grpc.NewClient("eberl.se:5000", grpc.WithTransportCredentials(creds))

	// the following works with service tcp 5000 -> 30009 but errors once switched to https
	// conn, err := grpc.NewClient("eberl.se:5000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewStreamerClient(conn)
	stream, err := client.EchoStream(context.Background())
	if err != nil {
		log.Fatalf("error opening stream: %v", err)
	}

	// Goroutine to handle incoming messages from the server
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note: %v", err)
			}
			log.Printf("Got server response: %s", in.Content)
		}
	}()

	// Loop to send messages every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	log.Println("Starting stream...")
	for range 10 {
		<-ticker.C
		msg := "Ping"
		if err := stream.Send(&pb.DataChunk{Content: msg}); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}

	stream.CloseSend()
	<-waitc
}
