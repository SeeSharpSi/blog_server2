package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"seesharpsi/htmx_quickstart/config"
	"seesharpsi/htmx_quickstart/handlers"
	"seesharpsi/htmx_quickstart/logger"
	"seesharpsi/htmx_quickstart/services"
	"seesharpsi/htmx_quickstart/session"
)

// custom404Handler wraps a mux to provide custom 404 pages
func custom404Handler(mux *http.ServeMux, h *handlers.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the route exists in the mux
		_, pattern := mux.Handler(r)
		if pattern == "" {
			// Route not found, serve custom 404
			h.NotFound(w, r)
			return
		}
		// Route exists, serve normally
		mux.ServeHTTP(w, r)
	})
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Setup structured logging with config
	logger.SetupLogger(cfg)

	slog.Info("configuration loaded",
		"server_addr", cfg.GetServerAddr(),
		"environment", os.Getenv("ENV"))

	sessionManager := session.NewManager(cfg)

	// Create service layer with dependencies
	service := services.NewService(sessionManager, slog.Default())

	// Create handler with injected service
	h := &handlers.Handler{
		Service: service,
	}

	// set up mux with middleware
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", h.Index)
	mux.HandleFunc("/test", h.Test)
	mux.HandleFunc("/health", h.Health)

	// Custom 404 handler for unmatched routes
	mux.HandleFunc("/404", h.NotFound)

	// Wrap with middleware and custom 404 handler
	handler := logger.PanicRecovery(logger.RequestLogger(custom404Handler(mux, h)))

	server := &http.Server{
		Addr:         cfg.GetServerAddr(),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Channel to listen for interrupt signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// start server
	slog.Info("starting server", "address", cfg.GetServerAddr())
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-done
	slog.Info("shutting down server")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited")
}
