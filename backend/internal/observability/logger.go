package observability

import (
	"log/slog"
	"os"
)

func NewLogger(service string) *slog.Logger {

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	)

	logger := slog.New(handler)

	return logger.With(
		slog.String("service", service),
	)
}
