package filesegment

import (
	"context"
	"fmt"
	"github.com/akimsavvin/sakv/internal/database/storage/wal"
	"github.com/akimsavvin/sakv/pkg/sl"
	"log/slog"
	"os"
)

type Streamer struct {
	log           *slog.Logger
	dataDirectory string
}

func NewStreamer(log *slog.Logger, dataDirectory string) *Streamer {
	return &Streamer{
		log:           log.With(sl.Comp("filesegment.Streamer")),
		dataDirectory: dataDirectory,
	}
}

func (s *Streamer) Stream(ctx context.Context) (<-chan wal.Segment, <-chan error) {
	s.log.DebugContext(ctx, "starting file segments reading")

	errs := make(chan error, 1)

	s.log.DebugContext(ctx, "reading data entries directory")
	dataEntries, err := os.ReadDir(s.dataDirectory)
	if err != nil {
		s.log.ErrorContext(ctx, "could not read directory entries", sl.Err(err))

		errs <- err
		close(errs)

		return nil, errs
	}
	s.log.DebugContext(ctx, "read directory entries", slog.Int("data_entries_count", len(dataEntries)))

	segments := make(chan wal.Segment, len(dataEntries))
	go func() (err error) {
		s.log.DebugContext(ctx, "streaming segments")

		defer close(segments)
		defer close(errs)

		defer func() {
			if err != nil {
				errs <- err
				s.log.ErrorContext(ctx, "segments streaming failed", sl.Err(err))
			} else {
				s.log.InfoContext(ctx, "segments streaming finished")
			}
		}()

		for _, dataEntry := range dataEntries {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			log := s.log.With(slog.String("segment_name", dataEntry.Name()))

			log.DebugContext(ctx, "creating segment from data entry")
			segment, err := New(fmt.Sprintf("%s/%s", s.dataDirectory, dataEntry.Name()))
			if err != nil {
				log.ErrorContext(ctx, "could not create segment from data entry", sl.Err(err))
				return err
			}

			segments <- segment
			log.DebugContext(ctx, "created segment from data entry")
		}

		return nil
	}()

	s.log.InfoContext(ctx, "file segments reading started")

	return segments, errs
}
