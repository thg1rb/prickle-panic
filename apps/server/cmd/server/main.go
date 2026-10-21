package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thg1rb/prickle-panic/apps/server/internal/appconfig"
	"github.com/thg1rb/prickle-panic/apps/server/internal/clientip"
	"github.com/thg1rb/prickle-panic/apps/server/internal/clock"
	"github.com/thg1rb/prickle-panic/apps/server/internal/httpapi"
	"github.com/thg1rb/prickle-panic/apps/server/internal/session"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := appconfig.Load(os.Getenv)
	if err != nil {
		return err
	}
	resolver, err := clientip.NewResolver(cfg.TrustedProxyCIDRs, cfg.TrustedProxyMaxHops)
	if err != nil {
		return fmt.Errorf("configure trusted client source resolver: %w", err)
	}
	manager, err := session.NewManager(session.Config{
		Clock:          clock.New(),
		Sink:           discardExpirySink{},
		IdleTTL:        cfg.IdleTTL,
		AbsoluteTTL:    cfg.AbsoluteTTL,
		CreationLimit:  cfg.CreationLimit,
		CreationWindow: cfg.CreationWindow,
		MaxSourceKeys:  cfg.MaxSourceKeys,
		OutboxCapacity: cfg.OutboxCapacity,
		RetryMin:       cfg.RetryMin,
		RetryMax:       cfg.RetryMax,
	})
	if err != nil {
		return fmt.Errorf("start session manager: %w", err)
	}
	defer manager.Close()

	api, err := httpapi.New(httpapi.Config{Mode: httpapi.Mode(cfg.Mode), Origin: cfg.Origin, Sessions: manager, ClientIP: resolver})
	if err != nil {
		return fmt.Errorf("configure HTTP session API: %w", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/api", api)
	mux.Handle("/api/", api)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serveErr := make(chan error, 1)
	go func() { serveErr <- server.ListenAndServe() }()
	slog.Info("HTTP server listening", "address", server.Addr, "environment", cfg.Mode)

	select {
	case err := <-serveErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	if err := <-serveErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// The W010 and W013 consumers do not exist yet. T03 has no live sockets or
// Matches, so the application acknowledges expiry events until those modules
// are integrated without adding their responsibilities here.
type discardExpirySink struct{}

func (discardExpirySink) Deliver(context.Context, session.ExpiryEvent) error { return nil }
