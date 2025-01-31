package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

const (
	LogValueSignal = "signal"
)

func New() *slog.Logger {
	return slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
}
