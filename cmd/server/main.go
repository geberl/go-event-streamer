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
MIIC4TCCAcmgAwIBAgIUaqVFfQS+EMoMPy4rjv9K9EfkpzswDQYJKoZIhvcNAQEL
BQAwADAeFw0yNjAzMTcxNTE4NDlaFw0yNzAzMTcxNTE4NDlaMAAwggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQDAzQnzVskGD++WQd4huNdzxsJc6LD0nsLq
TgpmPOdVSK1gyvS6olTOwOO9+n1U3JwGDOgMRhHc24dnXurYadNMm+xmHDjkgEgp
Vo0Hi/mL6N6FC+GygFMHSeuW9EePo9X/eZL1WjzLmcHDwyzrfCWJ3dDOj50TuV4h
3xWKJU2O4xO2hQvXVxyEqGwrSL9Zdq7GgmYBUijt1MMBhKe9SYpb8KM/Qeh5nFJ7
gMk4ZIrtyy6AWQIhfQh6LbBm6R7x8A7la4HtUJHAj/FFR74NG/vuOYaKYSkWMVIK
0+j1TQEjGcm9H5ViVH8RyqWoOCRRKPyUyQuo4Pv8Jcjhrbq4ejAZAgMBAAGjUzBR
MB0GA1UdDgQWBBQykgO2NYAwKqyWaXfbQ4g0y4g+QzAfBgNVHSMEGDAWgBQykgO2
NYAwKqyWaXfbQ4g0y4g+QzAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUA
A4IBAQAt+I/E3DondUw3RZ89/ycCo167tu0NSUIdFQ818g7GMVonMJOYcP7q38j5
8ydOZ8OziH7Y9xMGVXo1y9Xc5V1wKciThgi39CJJkepo6+dGlWriz/qghExeJc0Z
SQWeWOkFP67+QKRNCWrfPY59RtJDi10ewj82qYXQw3qsEj2m+eupgHzi3c03i/xe
XKH3LDMhqnbDSswCqlamciBe7WSNSdJjjvPq2ftvDFpZ5gyPqlBxRnpEsyDi+qcF
O9Cg4jg7tBwhmwWE/+V3SA23qntpEAwFXP5/9l1yOp9VUzajrRo28C/yalR4cpok
XpD9Fy9BRz1EiavGr4ihq78ESGS9
-----END CERTIFICATE-----`

const keyPEM = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDAzQnzVskGD++W
Qd4huNdzxsJc6LD0nsLqTgpmPOdVSK1gyvS6olTOwOO9+n1U3JwGDOgMRhHc24dn
XurYadNMm+xmHDjkgEgpVo0Hi/mL6N6FC+GygFMHSeuW9EePo9X/eZL1WjzLmcHD
wyzrfCWJ3dDOj50TuV4h3xWKJU2O4xO2hQvXVxyEqGwrSL9Zdq7GgmYBUijt1MMB
hKe9SYpb8KM/Qeh5nFJ7gMk4ZIrtyy6AWQIhfQh6LbBm6R7x8A7la4HtUJHAj/FF
R74NG/vuOYaKYSkWMVIK0+j1TQEjGcm9H5ViVH8RyqWoOCRRKPyUyQuo4Pv8Jcjh
rbq4ejAZAgMBAAECggEAQ3GWWQLTWGULtR6+g0JjT/NH+3NEr5W37nm1RpVogRtm
1xS4Lm9pxleQc10kKaLwi2dJZz29suoygBUihujiCwsCU6fsuPYtCBToSasL9QbV
jGofHi+om8SefpReUh+IVRGkuGJEIR7cusvUM14ezY8EI7X2RzeRd7zPjp9E9cXb
GbiVOEpS9RiTAVF/FcRhjvsIPJn3+TrTxcrBuDYu7oNMUl5r0wKik/w7KYtFPg/e
eyRAlhr5aSwVvRm8dlKZyjwgDX7rUlOCmcnidsuta+jTsE7h9UMRvnDUNG2ruo88
TXl2PgCWRinSxkGgkoi3HQ3ewoXM3g0YKr85eU7gvwKBgQDkbvXsGJhaOIXfjDBE
TXBUPCNKAPlzDV59WKMVS24yx9WwUpZtiF1mPY6hSJLb3hXE3PyGr+8UbBYlEwp1
cspX7mFYpaVct8HV8c7g0HSqUnmwcTBAb2SvW1377m20UDCA7AfICeeGmvVjnEMp
JHZ+UY5ckgEoiX2rgZyIBQ2n/wKBgQDYEUaytGHOLffQRonyEw9mi9ELn8DLJZ/l
Rhhu9pWLK/DL/YGhmKlcpjXsqr5DDHU5/XeywNvrG+TZTrpTRHou6EdEbLQAvniM
jVZvS0qbCf0oGU64fQZYLRWKKoAuYv/82RdE8K272VxJB9Cz6FOUvOxOwCuo3C/n
YBXL1XBn5wKBgEnj9pqDLizo4az4/NfrMK2esk+K1yW3KlxjYoVN2/yDFYUugcg2
dvfOa6eSAScrxGDklq6+lBhICjW93gE1u2wMCOMS2dWO/x1EVYX1B/fcK86+HjyJ
i8kJRfJrIoNT+QyKzM2RHpo037Fz52mUiNu9Z85b0BIbv1HN4CNDdzJjAoGAWLc3
PRb9dafAMb9U0pVq5GMSMWCly4OmVIBkdeM/YcZn94oeWNiS6ZzBVWyB9Iu/8lCV
fkrbwXxRibxemuPp+yqaYIj1m7yZSLSbwdS7TE9cp8NEZFHJchkI2BM9UE6L5yjH
+iGMZC4KS14vHj+NWev8ZxVWl93YuXrlWC1KGw8CgYEAqysx0edR7nHA9aGseyCM
/sEi4ZMZ3JdH7Zf084sWNKD5jPdpr5Qi2jsa4paX1CmLdAAr6Q+DOgJshalzguhQ
TDCCb93OHZtRe/7jfysA04hesYPOuyxhxOmgAh2dNHd4L2GBfSRkFcEbtEY8H8Iz
K5xiIBwYIG9t3953e4Ob26E=
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
