package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

const (
	LogValueStatus     = "status"
	LogValuePath       = "path"
	LogValueFolder     = "folder"
	LogValueFile       = "file"
	LogValueStartedAt  = "started_at"
	LogValueFinishedAt = "finished_at"
	LogValueDuration   = "duration"
)

var (
	Logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
	Debug = Logger.Debug
	Info  = Logger.Info
	Warn  = Logger.Warn
	Error = Logger.Error
)
