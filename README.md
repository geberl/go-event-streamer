# go-event-streamer

*Minimal reproducer to debug a gRPC streaming API*

## Generate gRPC Client

```sh
brew install protoc-gen-go-grpc
brew install protoc-gen-go

protoc \
  --go_out=. \
  --go-grpc_out=. \
  *.proto
```

## Run Server

```sh
export EVENT_STREAMER_HOSTNAME="t1000"
export EVENT_STREAMER_USE_CERTS="true"

go run ./cmd/server/...
```

## Run gRPC Clients

```sh
go run ./cmd/client/... -mode stream
go run ./cmd/client/... -mode unary

grpcurl -insecure -plaintext=false example.com:5000 list
```
## Run HTTP/2 Client

```sh
# secure h2
curl -H2 -k https://example.com:5001/

# insecure h2c
curl --http2-prior-knowledge http://example.com:5001/
```
