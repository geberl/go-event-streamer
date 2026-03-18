package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	pb "event-streamer/provider"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
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

const certPEM = `-----BEGIN CERTIFICATE-----
MIIDJTCCAg2gAwIBAgIUVfa5+1D+HD8gIamjp9ZX8cJyjIQwDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTI2MDMxODEyNTQzOFoXDTM2MDMx
NTEyNTQzOFowFDESMBAGA1UEAwwJbG9jYWxob3N0MIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEAuq4HreXSgvPbn9e1t6dUZ+C2HM8omgWicjs5kB5ub7RP
A9E5HQH3naTiA4ytiRLHvUdvTlYUbGmc6wMOD+dCUK2eraJAAc3JngHfWxdgw3iT
gvB+QzA/KKVk7Sl7GIFV5aeSkI3Z7UpfNPEFDJU4q7+IfoZMe7efRf5A4DK4ouhd
2mYD1qImwhzx/gQ5O/htI8EeueRqL6/oTIwm2wlITmPQmd9XVtCy9wipW1md6DEM
OfahNoq0gJZOobqv82iwjwhT0K8T3rBB0pbRrzT0TMT+a0YYho/JOwLQuRLG4eNO
OD7wypyKybwUfuNlA7tTX9+j/6s2vGiHxiUIloRzJQIDAQABo28wbTAdBgNVHQ4E
FgQUZJCfQSzu1ztTvDMKz7fgq28m9ocwHwYDVR0jBBgwFoAUZJCfQSzu1ztTvDMK
z7fgq28m9ocwDwYDVR0TAQH/BAUwAwEB/zAaBgNVHREEEzARgglsb2NhbGhvc3SH
BH8AAAEwDQYJKoZIhvcNAQELBQADggEBAGJFz0NxTGK6JJgMMIXZ/oZmYzGJxnyV
nzQ5Qme6GAgDnFPOD1MWZQkh9de/oANWPe+rzc5owGLsgKAQZSw0H02mjAL/1kg7
s+G7VhyXP7YJoIrBJRZ7xCkbYQ/+Ke4poN3kFDwXeqVoSiWeDL0pZg9ByA8Z7iZH
VIZzoHNkshp6eMQuNXXh85+k9TQBxKin6Iwozz6yGLHAeGZUOmvH6+GU9EbNIA2O
bVNaeq3ILV+4DWtlOrFBER/t4k/hXHzecXqYz4rCEv54uiYZgJ4bPViWjXgQXy2T
SkSq0uYOvE9360TNNV6zXwE0imVtjW0pIBH/9XLsXiGs4Q7SgPpF74A=
-----END CERTIFICATE-----`

const keyPEM = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC6rget5dKC89uf
17W3p1Rn4LYczyiaBaJyOzmQHm5vtE8D0TkdAfedpOIDjK2JEse9R29OVhRsaZzr
Aw4P50JQrZ6tokABzcmeAd9bF2DDeJOC8H5DMD8opWTtKXsYgVXlp5KQjdntSl80
8QUMlTirv4h+hkx7t59F/kDgMrii6F3aZgPWoibCHPH+BDk7+G0jwR655Govr+hM
jCbbCUhOY9CZ31dW0LL3CKlbWZ3oMQw59qE2irSAlk6huq/zaLCPCFPQrxPesEHS
ltGvNPRMxP5rRhiGj8k7AtC5Esbh4044PvDKnIrJvBR+42UDu1Nf36P/qza8aIfG
JQiWhHMlAgMBAAECggEACqOWbnO1np71OlPZ2GCh79Wfq16nCrgdfPMhIbSKSLV2
91m6LowJJ6PY+ajPzwsR9RiYIFfJjDAssDwZVhCw99YdP/oKOdAXmHi02QUpD5rU
lVbNa1jZkKB0cwu1Jz1fvtnhAXoEHIDrkiHWTtRGSYt95PAUdcyOODf4TI63dRk9
31x74Ckzeb+s0e+IDcWMqqAKOt+Wdz0Np7voyE7WR7TIgBtDzKNVyOLP9UzQB8JG
82LjawyYzMIkdpHyMPDV2j7VXYHRT5Y6Be+PvIjcLamkY6fbRxP3NxiBKlsaEZrD
PtGTt8zdb9R3Ol1OqujLc6VnKiHvZ661U1l9eSFTFQKBgQD/gkztfV0A8Y5xgZnJ
pfa/GjPTtr5ZxwXDn80sI+rCi5bU33AiolzS6HLqAm1TYHNnuqQMzyaz7E2jyLVM
djJum9h+t/mL6cTNQVRJF30sVVa3M6roCuPEUSOoM6jL1Axth8bjUcKKD3R8Ffvf
5F/n5Adx8NjoMEKP0+U8I3QgrwKBgQC7Cd5U4d3hsEN6r8wbdzC8wvmY8vxMqwuM
DIcDv3BR6CawaNF/+ClpMSg62aiM1q3C0MKWvRDN1az+94/tDu2Myn6/5OE0PZrP
PF3ESP38ke8houye1pvkWjq1/721Xj/S8PphaKLPrF9gAycxNrKlyCPyNa1lH6Py
L1Tr3gBWawKBgQCZrMwJ9uGOJLrwp+tQLgK3M9JCHuJj6uEbpKxpRPz4n65LQEwY
eKDttSMQff81K4idtdLfZWQ4yQJ1ZM0uPNTeU9ulc4+iyCo27Xj9MSR3Gqi6LVg4
kfwl4ktY6iE23sXOxuAnbtBb6ym7TBmesqPAPBUCQcKj/Aq8qMxyHDzHPwKBgQCF
p9eo8H6N+FdAJL/GILZDLVEPaxO/9bqaqZkRpIuu/CYpib2rpLpy4R3OcBtyCTbC
MEvdS93mOPsWd/HxhOlb4pgQqI4FtsAZtxmKWl6lTeOENdjA6LsdwxyRUd9O67rQ
EkPZt9wgaxz8j0RCdsPSk+KcAp+V07ZkKk6U/l9fYQKBgEA73dfCxOtZpkO8J/+X
SgSLQdtmfZ/MGmafnJeQyWEIjTlAuZT/aXjPBvoSNMUHPf0g1c7dmhXJX7Z5AQdy
34nM2g7pCHcUt4deVgBrLI2esH95X2tKaI6HHF5jVZQ9bCIULsvxk95ZqHDW7fcM
Hz/xP8KkY3hu8Z4w9Kg2e9g9
-----END PRIVATE KEY-----`

func main() {
	var opts []grpc.ServerOption

	if os.Getenv("EVENT_STREAMER_USE_CERTS") != "" {
		cert, err := tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
		if err != nil {
			log.Fatalf("failed to load key pair: %s", err)
		}

		log.Println("Mode: TLS Enabled")
		creds := credentials.NewServerTLSFromCert(&cert)
		opts = append(opts, grpc.Creds(creds))
	} else {
		log.Println("Mode: INSECURE (Plaintext)")
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(opts...)
	pb.RegisterStreamerServer(s, &server{})
	reflection.Register(s)

	log.Println("Server listening on :50051 (Unary + Streaming) ...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
