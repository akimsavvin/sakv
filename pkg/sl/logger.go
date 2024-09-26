package sl

import (
	slogmulti "github.com/samber/slog-multi"
	"log/slog"
	"os"
	"sakv/internal/database/config"
)

func NewLogger(cfg config.Logging) (log *slog.Logger, close func() error, err error) {
	f, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	close = f.Close

	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}

	log = slog.New(slogmulti.Fanout(
		slog.NewJSONHandler(f, opts),
		slog.NewTextHandler(os.Stdout, opts),
	))

	return
}
