# go-event-streamer

**

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

## Run Client

```sh
go run ./cmd/client/... -mode stream
go run ./cmd/client/... -mode unary
```
