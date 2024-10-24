package restoration

import (
	"context"
	"github.com/akimsavvin/sakv/internal/database/compute/query"
	"github.com/akimsavvin/sakv/internal/database/storage/engine"
	"github.com/akimsavvin/sakv/pkg/sl"
	"log/slog"
)

type Restorer struct {
	log *slog.Logger
}

func NewRestorer(log *slog.Logger) *Restorer {
	return &Restorer{
		log: log.With(sl.Comp("restoration.Restorer")),
	}
}

func (r *Restorer) Restore(ctx context.Context, e engine.Engine, queries <-chan string, errs <-chan error) error {
	r.log.DebugContext(ctx, "engine queries restoring started")

	for {
		select {
		case queryStr, ok := <-queries:
			if !ok {
				r.log.InfoContext(ctx, "engine queries restoring finished")
				return nil
			}

			log := r.log.With(slog.String("query", queryStr))
			log.InfoContext(ctx, "handling query")

			q, _ := query.ParseQueryStr(queryStr)
			switch q.Command() {
			case query.CommandSET:
				_ = e.SET(ctx, q.Arg(0), q.Arg(1))
			case query.CommandDEL:
				_ = e.DEL(ctx, q.Arg(0))
			}

			log.InfoContext(ctx, "query handled", slog.String("query", queryStr))
		case err := <-errs:
			if err != nil {
				r.log.ErrorContext(ctx, "engine queries restoring failed", sl.Err(err))
				return err
			}

			return nil
		}
	}
}
