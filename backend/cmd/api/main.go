package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/saral-gupta7/recode/backend/internal/config"
	"github.com/saral-gupta7/recode/backend/internal/observability"
)


func main() {
	logger := observability.NewLogger("api")


	if err := run(logger); err != nil {
		logger.Error("api stopped", slog.Any("error", err),
	)
	

	os.Exit(1)
	}
}



func run (logger *slog.Logger) error {
	cfg, err := config.Load()

	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	logger.Info(
		"api configuration loaded",
		slog.String("http_address", cfg.HTTPAddress),
		slog.Duration("shutdown_timeout", cfg.ShutdownTimeout),
	)
	
	return nil
}