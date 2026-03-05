package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DaisukeMatsuoh/simple-gallery/internal/config"
	"github.com/DaisukeMatsuoh/simple-gallery/internal/server"
	"github.com/DaisukeMatsuoh/simple-gallery/internal/store"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "config.toml", "path to configuration file")
	flag.Parse()

	// Set up structured logging
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(logHandler)

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	// Initialize database
	db, err := store.New(cfg.Storage.DBPath, logger)
	if err != nil {
		logger.Error("failed to initialize database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	// Apply migrations
	migrator := store.NewMigrator(db, logger)
	if err := migrator.ApplyMigrations(); err != nil {
		logger.Error("failed to apply migrations", slog.Any("error", err))
		os.Exit(1)
	}

	// Create server
	srv := server.New(cfg, db, logger)

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Start()
	}()

	// Wait for shutdown signal or server error
	select {
	case <-sigChan:
		logger.Info("received shutdown signal")
	case err := <-errChan:
		if err != nil {
			logger.Error("server error", slog.Any("error", err))
			os.Exit(1)
		}
	}

	// Graceful shutdown
	if err := srv.Shutdown(5 * time.Second); err != nil {
		logger.Error("failed to shutdown server", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("server stopped")
}
