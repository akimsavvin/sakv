package filesegment

import (
	"fmt"
	"github.com/akimsavvin/sakv/internal/database/storage/wal"
	"github.com/akimsavvin/sakv/pkg/sl"
	"log/slog"
	"time"
)

type Factory struct {
	log *slog.Logger
}

func NewFactory(log *slog.Logger) *Factory {
	return &Factory{
		log: log.With(sl.Comp("filesegment.Factory")),
	}
}

func (f *Factory) Create(dir string) (s wal.Segment, err error) {
	f.log.Debug("creating a new file segment")
	now := time.Now()
	name := fmt.Sprintf("%s/segment_%d.txt", dir, now.Unix())
	s, err = New(name)
	if err != nil {
		f.log.Error("could not create the file segment", sl.Err(err))
	} else {
		f.log.Info("created the file segment")
	}

	return
}
