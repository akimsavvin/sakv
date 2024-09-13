package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"sakv/internal/compute/listener"
	"sakv/internal/compute/query"
	inmemory "sakv/internal/storage/engine/in-memory"
	"sakv/pkg/sl"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	inMemoryEngine := inmemory.NewEngine()
	h := query.NewHandler(inMemoryEngine)
	l := listener.New(h)

	if err := l.StartListening(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Info("application stopped")
			return
		}

		log.Error("some error occurred while listening", sl.Err(err))
	}
}
