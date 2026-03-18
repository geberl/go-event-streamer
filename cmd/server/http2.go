package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func RunHTTP2Server(
	ctx context.Context,
	addr string,
	errChan chan<- error,
) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello HTTP/2\n"))
	})

	go func() {
		if os.Getenv("EVENT_STREAMER_USE_CERTS") != "" {
			cert, err := tls.X509KeyPair([]byte(CertPEM), []byte(KeyPEM))
			if err != nil {
				slog.ErrorContext(ctx, "http2 failed to load key pair", "err", err)
				errChan <- err
				return
			}
			slog.InfoContext(ctx, "http2 tls mode enabled")

			srv := &http.Server{
				Addr:    addr,
				Handler: mux,
				TLSConfig: &tls.Config{
					MinVersion:   tls.VersionTLS12,
					Certificates: []tls.Certificate{cert},
					NextProtos:   []string{"h2"},
				},
			}

			if err := http2.ConfigureServer(srv, &http2.Server{}); err != nil {
				errChan <- err
				return
			}

			ln, err := tls.Listen("tcp", addr, srv.TLSConfig)
			if err != nil {
				errChan <- err
				return
			}

			slog.InfoContext(ctx, "http2 server running", "address", addr)

			if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.ErrorContext(ctx, "http2 failed to start server", "err", err)
				errChan <- err
			}
		} else {
			slog.InfoContext(ctx, "http2 INSECURE mode enabled")

			h2s := &http2.Server{}

			srv := &http.Server{
				Addr:    addr,
				Handler: h2c.NewHandler(mux, h2s),
			}

			slog.InfoContext(ctx, "http2 h2c server running", "address", addr)

			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.ErrorContext(ctx, "http2 failed to start h2c server", "err", err)
				errChan <- err
			}
		}
	}()
}
