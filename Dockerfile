FROM golang:1.26.1-alpine3.23 AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download && go mod verify
COPY . .
ARG VERSION
RUN VERSION=$(git describe --tags --always --dirty=-dev 2>/dev/null || echo "dev") && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X 'main.Version=${VERSION}'" \
    -o server.bin \
    ./cmd/server

FROM alpine:3.23.3
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/server.bin ./server
EXPOSE 50051
EXPOSE 50052
ENTRYPOINT ["/app/server"]
