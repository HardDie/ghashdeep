package crawler

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

const (
	LogValueStatus     = "status"
	LogValuePath       = "path"
	LogValueFolder     = "folder"
	LogValueFile       = "file"
	LogValueStartedAt  = "started_at"
	LogValueFinishedAt = "finished_at"
	LogValueDuration   = "duration"
)
