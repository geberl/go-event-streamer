brew install protoc-gen-go-grpc
brew install protoc-gen-go

protoc \
  --go_out=. \
  --go-grpc_out=. \
  *.proto

go run ./cmd/server/...

go run ./cmd/client/...