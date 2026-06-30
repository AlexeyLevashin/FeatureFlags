package logger

import (
	"log/slog"
	"os"
)

func Setup() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	log := slog.New(handler)

	slog.SetDefault(log)

	return log
}
