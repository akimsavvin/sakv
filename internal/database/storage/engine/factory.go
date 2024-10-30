package engine

import (
	"github.com/akimsavvin/sakv/internal/database/config"
	"github.com/akimsavvin/sakv/internal/database/storage/engine/in-memory"
	"github.com/akimsavvin/sakv/pkg/sl"
	"log/slog"
)

const (
	InMemory = "in_memory"
)

type Factory interface {
	CreateEngine(cfg config.Engine) Engine
}

type factory struct {
	log *slog.Logger
}

func NewFactory(log *slog.Logger) Factory {
	return &factory{
		log: log,
	}
}

func (f *factory) CreateEngine(cfg config.Engine) Engine {
	log := f.log.With(sl.Comp("engine.factory"))
	log.Debug("creating an engine", slog.Any("config", cfg))

	switch cfg.Type {
	case InMemory:
		log.Info("creating the in memory engine")
		return inmemory.NewEngine(f.log)
	default:
		log.Error("unknown engine type")
		panic("unknown engine type")
	}
}
