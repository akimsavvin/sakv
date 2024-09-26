package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"sakv/internal/database/compute/query"
	"sakv/internal/database/config"
	"sakv/internal/database/network/listener"
	"sakv/internal/database/storage/engine"
	"sakv/pkg/sl"
	"syscall"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "./config/config.yaml", "path to sakv config file")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	cfg := config.MustNew(configFilePath)

	log, close, err := sl.NewLogger(cfg.Logging)
	if err != nil {
		panic(err)
	}
	defer close()

	ef := engine.NewFactory(log)
	e := ef.CreateEngine(cfg.Engine)
	h := query.NewHandler(log, e)
	l, err := listener.New(log, cfg.Network, h)
	if err != nil {
		panic(err)
	}

	if err := l.StartListening(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Info("application stopped")
			return
		}

		log.Error("some error occurred while listening", sl.Err(err))
	}
}
