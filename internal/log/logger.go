package log

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func getLogger() *slog.Logger {
	if logger == nil {
		var handler slog.Handler

		if true {
			handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			})
		} else {
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			})
		}

		logger = slog.New(handler)
	}

	return logger
}

func Debug(msg string, args ...any) {
	logger := getLogger()

	logger.Debug(msg, args...)
}
