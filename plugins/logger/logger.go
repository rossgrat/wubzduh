package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/rossgrat/wubzduh/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(cfg config.LogConfig) *slog.Logger {
	if cfg.Path == "" {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	fileWriter := &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
	}

	multi := io.MultiWriter(os.Stdout, fileWriter)
	return slog.New(slog.NewJSONHandler(multi, nil))
}
