package debug

import (
	"log/slog"
	"os"
	"strings"
)

// FYLINE_LOG controls verbosity: warn, info (default), debug
// debug traces everything (sends, receives, channel switches)
var logger *slog.Logger

func init() {
	level := parseLevel(os.Getenv("FYLINE_LOG"))
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: level == slog.LevelDebug,
	}
	logger = slog.New(slog.NewTextHandler(os.Stderr, opts))
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}
