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

	applicationjobs "github.com/saral-gupta7/recode/backend/internal/application/jobs"
	"github.com/saral-gupta7/recode/backend/internal/config"
	"github.com/saral-gupta7/recode/backend/internal/database"
	"github.com/saral-gupta7/recode/backend/internal/httpapi"
	"github.com/saral-gupta7/recode/backend/internal/observability"
	"github.com/saral-gupta7/recode/backend/internal/queue"
	"github.com/saral-gupta7/recode/backend/internal/repository/postgres"
	"github.com/saral-gupta7/recode/backend/internal/storage"
)

const dependencyStartupTimeout = 5 * time.Second

func main() {
	logger := observability.NewLogger("api")
	if err := run(logger); err != nil {
		logger.Error("api stopped", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	startupContext, cancelStartup := context.WithTimeout(
		context.Background(),
		dependencyStartupTimeout,
	)
	defer cancelStartup()

	pool, err := database.Open(startupContext, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer pool.Close()

	jobQueue, err := queue.OpenRedis(startupContext, cfg.RedisAddress, cfg.QueueName)
	if err != nil {
		return fmt.Errorf("open queue: %w", err)
	}
	defer jobQueue.Close()

	store, err := storage.NewLocal(cfg.StorageRoot)
	if err != nil {
		return fmt.Errorf("open storage: %w", err)
	}

	repository := postgres.NewTaskRepository(pool)
	jobService := applicationjobs.New(
		repository,
		jobQueue,
		store,
		cfg.MaxUploadBytes,
		cfg.ResultRetention,
	)

	server := &http.Server{
		Addr: cfg.HTTPAddress,
		Handler: httpapi.NewRouter(logger, httpapi.Dependencies{
			Jobs:           jobService,
			MaxUploadBytes: cfg.MaxUploadBytes,
			Readiness: func(ctx context.Context) error {
				if err := pool.Ping(ctx); err != nil {
					return err
				}
				return jobQueue.Ping(ctx)
			},
		}),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	logger.Info(
		"api listening",
		slog.String("http_address", cfg.HTTPAddress),
	)

	signalContext, stopSignals := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stopSignals()

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return nil
	case <-signalContext.Done():
		logger.Info("shutdown signal received")
	}

	shutdownContext, cancelShutdown := context.WithTimeout(
		context.Background(),
		cfg.ShutdownTimeout,
	)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownContext); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	if err := <-serverErrors; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP during shutdown: %w", err)
	}

	logger.Info("api stopped gracefully")
	return nil
}
