package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	lv := slog.LevelInfo
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.Level(lv)}))
	slog.SetDefault(logger)

	slog.Info(
		"starting",
		"name", "event-streamer",
	)

	ctx := context.Background()
	if err := run(ctx); err != nil {
		slog.Error("failed to start servers", "err", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	errChan := make(chan error, 2)

	RunHTTP2Server(ctx, ":50051", errChan)
	RunGrpcServer(ctx, ":50052", errChan)

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		slog.InfoContext(ctx, "shutting down servers")
	}

	return nil
}
