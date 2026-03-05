package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/DaisukeMatsuoh/simple-gallery/internal/config"
	"github.com/DaisukeMatsuoh/simple-gallery/internal/store"
)

type Server struct {
	config *config.Config
	db     *store.DB
	mux    *http.ServeMux
	log    *slog.Logger
	srv    *http.Server
}

// New creates a new Server instance
func New(cfg *config.Config, db *store.DB, logger *slog.Logger) *Server {
	mux := http.NewServeMux()

	srv := &Server{
		config: cfg,
		db:     db,
		mux:    mux,
		log:    logger,
		srv: &http.Server{
			Addr:         net.JoinHostPort(cfg.Server.Host, fmt.Sprintf("%d", cfg.Server.Port)),
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}

	// Register routes
	registerRoutes(srv)

	return srv
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.log.Info("starting server", slog.String("addr", s.srv.Addr))
	return s.srv.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	s.log.Info("shutting down server")
	return s.srv.Shutdown(ctx)
}

// GetDB returns the database connection
func (s *Server) GetDB() *store.DB {
	return s.db
}

// GetLog returns the logger
func (s *Server) GetLog() *slog.Logger {
	return s.log
}
