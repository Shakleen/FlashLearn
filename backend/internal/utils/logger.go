package utils

import (
	"log/slog"
	"os"
)

func GetLogger() *slog.Logger {
	handlerOptions := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))
	return logger
}
