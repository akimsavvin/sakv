package app

import (
	"context"
	"errors"
	"github.com/akimsavvin/sakv/internal/database/compute/query"
	"github.com/akimsavvin/sakv/internal/database/config"
	"github.com/akimsavvin/sakv/internal/database/network/server"
	"github.com/akimsavvin/sakv/internal/database/storage/engine"
	"github.com/akimsavvin/sakv/internal/database/storage/restoration"
	"github.com/akimsavvin/sakv/internal/database/storage/wal"
	filesegment "github.com/akimsavvin/sakv/internal/database/storage/wal/file-segment"
	"github.com/akimsavvin/sakv/pkg/sl"
)

func Start(ctx context.Context, configFilePath string) error {
	cfg := config.MustNew(configFilePath)

	log, close, err := sl.NewFileLogger(cfg.Logging)
	if err != nil {
		return err
	}
	defer close()

	ef := engine.NewFactory(log)
	e := ef.CreateEngine(cfg.Engine)

	sf := filesegment.NewFactory()
	ss := filesegment.NewStreamer(log, cfg.WAL.DataDirectory)

	var walInst *wal.WAL

	if cfg.WAL.Enabled {
		walInst, err = wal.New(ctx, log, cfg.WAL, sf, ss)
		if err != nil {
			return err
		}

		queries, errs := walInst.Recover(ctx)
		r := restoration.NewRestorer(log)
		if err = r.Restore(ctx, e, queries, errs); err != nil {
			return err
		}
	}

	h := query.NewHandler(log, walInst, e)
	l, err := server.New(log, cfg.Network, h)
	if err != nil {
		return err
	}

	if err := l.Start(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			log.Info("application stopped")
			return nil
		}

		log.Error("some error occurred while listening", sl.Err(err))
		return err
	}

	return nil
}
