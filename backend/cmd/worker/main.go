package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saral-gupta7/recode/backend/internal/config"
	"github.com/saral-gupta7/recode/backend/internal/database"
	"github.com/saral-gupta7/recode/backend/internal/media"
	"github.com/saral-gupta7/recode/backend/internal/observability"
	"github.com/saral-gupta7/recode/backend/internal/queue"
	"github.com/saral-gupta7/recode/backend/internal/repository/postgres"
	"github.com/saral-gupta7/recode/backend/internal/storage"
	"github.com/saral-gupta7/recode/backend/internal/worker"
)

func main() {
	logger := observability.NewLogger("worker")
	if err := run(logger); err != nil {
		logger.Error("worker stopped", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	signalContext, stopSignals := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stopSignals()

	startupContext, cancelStartup := context.WithTimeout(signalContext, 5*time.Second)
	pool, err := database.Open(startupContext, cfg.DatabaseURL)
	if err != nil {
		cancelStartup()
		return fmt.Errorf("open database: %w", err)
	}
	defer pool.Close()

	jobQueue, err := queue.OpenRedis(startupContext, cfg.RedisAddress, cfg.QueueName)
	cancelStartup()
	if err != nil {
		return fmt.Errorf("open queue: %w", err)
	}
	defer jobQueue.Close()

	store, err := storage.NewLocal(cfg.StorageRoot)
	if err != nil {
		return fmt.Errorf("open storage: %w", err)
	}

	runner := worker.New(
		logger,
		postgres.NewTaskRepository(pool),
		jobQueue,
		store,
		media.NewFFprobe(cfg.FFprobePath),
		media.NewFFmpeg(cfg.FFmpegPath),
		cfg.ResultRetention,
		cfg.MaxUploadBytes*2,
	)
	if err := runner.Run(signalContext); err != nil {
		return fmt.Errorf("run worker: %w", err)
	}

	logger.Info("worker stopped gracefully")
	return nil
}
