package logger

import (
	"log/slog"
	"os"
)

const (
	LogValueStatus     = "status"
	LogValuePath       = "path"
	LogValueFile       = "file"
	LogValueStartedAt  = "started_at"
	LogValueFinishedAt = "finished_at"
	LogValueDuration   = "duration"
)

var (
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		//AddSource: true,
		Level: slog.LevelDebug,
	}))
	Debug = Logger.Debug
	Info  = Logger.Info
	Warn  = Logger.Warn
	Error = Logger.Error
)
