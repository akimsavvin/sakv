package wal

import (
	"bytes"
	"context"
	"github.com/akimsavvin/sakv/internal/database/config"
	"github.com/akimsavvin/sakv/pkg/concur/promise"
	"github.com/akimsavvin/sakv/pkg/memory"
	"github.com/akimsavvin/sakv/pkg/sl"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"
)

type Segment interface {
	io.ReadWriteCloser
	Name() string
	Size() (uint, error)
}

type DirectorySegments interface {
	DirectorySegments() ([]Segment, error)
}

type SegmentFactory interface {
	Create(dir string) (Segment, error)
}

type SegmentsStreamer interface {
	Stream(ctx context.Context) (<-chan Segment, <-chan error)
}

type queryPromise struct {
	query      []byte
	errPromise *promise.Error
}

type WAL struct {
	log *slog.Logger
	cfg config.WAL

	sf SegmentFactory
	ss SegmentsStreamer

	mx             sync.Mutex
	currentSegment Segment
	maxSegmentSize uint

	flushingTimeout time.Duration
	queries         chan queryPromise
}

func New(ctx context.Context, log *slog.Logger, cfg config.WAL, sf SegmentFactory, ss SegmentsStreamer) (*WAL, error) {
	flushingTimeout, err := time.ParseDuration(cfg.FlushingBatchTimeout)
	if err != nil {
		return nil, newCreationError(err)
	}

	maxSegmentSize, err := memory.ParseSize(cfg.MaxSegmentSize)
	if err != nil {
		return nil, newCreationError(err)
	}

	wal := &WAL{
		log: log.With(sl.Comp("wal.WAL")),
		cfg: cfg,

		sf: sf,
		ss: ss,

		maxSegmentSize: maxSegmentSize,

		flushingTimeout: flushingTimeout,
		queries:         make(chan queryPromise),
	}

	go wal.run(ctx)

	return wal, nil
}

func (wal *WAL) newSegment(closeLast bool) (Segment, error) {
	wal.log.Debug("creating a new segment")
	s, err := wal.sf.Create(wal.cfg.DataDirectory)
	if err != nil {
		wal.log.Error("could not create a new segment", sl.Err(err))
		return nil, err
	}

	if closeLast {
		_ = wal.currentSegment.Close()
	}

	wal.currentSegment = s
	wal.log.Debug("created a new segment")

	return s, nil
}

func (wal *WAL) lastSegment() (Segment, error) {
	wal.log.Debug("getting last WAL segment")

	if wal.currentSegment == nil {
		return wal.newSegment(false)
	}

	wal.log.Debug("received last segment",
		slog.String("segment_name", wal.currentSegment.Name()))

	return wal.currentSegment, nil
}

func (wal *WAL) processBatchWithLock(batch []queryPromise) {
	wal.mx.Lock()
	defer wal.mx.Unlock()

	wal.processBatch(batch)
}

func (wal *WAL) processBatch(batch []queryPromise) {
	wal.log.Debug("starting batch processing")

	writeErr := func(batch []queryPromise, err error) {
		for _, qp := range batch {
			qp.errPromise.Set(err)
		}
	}

	s, err := wal.lastSegment()
	if err != nil {
		wal.log.Error("batch processing failed due to the last segment error", sl.Err(err))
		writeErr(batch, err)
		return
	}
	wal.log.Debug("received last segment", slog.String("segment_name", s.Name()))

	size, err := s.Size()
	if err != nil {
		wal.log.Error("batch processing failed due to the size of the segment", sl.Err(err))
		writeErr(batch, err)
		return
	}

	results := make([]bytes.Buffer, 1)
	for _, qp := range batch {
		n := uint(len(qp.query))

		if size+n <= wal.maxSegmentSize {
			results[len(results)-1].Write(qp.query)
			size += n
		} else {
			size = 0
			results = append(results, bytes.Buffer{})
		}
	}

	for i, result := range results {
		if i != 0 {
			s, err = wal.newSegment(true)
			if err != nil {
				writeErr(batch, err)
				return
			}
		}

		if _, err := s.Write(result.Bytes()); err != nil {
			writeErr(batch, err)
			return
		}
	}

	writeErr(batch, nil)
}

func (wal *WAL) run(ctx context.Context) {
	wal.log.DebugContext(ctx, "running WAL")

	ticker := time.NewTicker(wal.flushingTimeout)
	defer ticker.Stop()

	var queriesBatch []queryPromise
	reset := func() {
		queriesBatch = make([]queryPromise, 0, wal.cfg.FlushingBatchSize)
	}

	for {
		select {
		case <-ctx.Done():
			wal.log.DebugContext(ctx, "WAL stopped")
			return
		case <-ticker.C:
			if len(queriesBatch) > 0 {
				go wal.processBatchWithLock(queriesBatch)
				reset()
			}
		case query := <-wal.queries:
			queriesBatch = append(queriesBatch, query)
			if len(queriesBatch) == wal.cfg.FlushingBatchSize {
				go wal.processBatchWithLock(queriesBatch)
				reset()
			}
		}
	}
}

func (wal *WAL) Write(ctx context.Context, query string) error {
	p := promise.NewError()
	wal.queries <- queryPromise{
		query:      []byte(query + "\n"),
		errPromise: p,
	}

	select {
	case <-p.Awaiter():
		return p.MustGet()
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (wal *WAL) Recover(ctx context.Context) (<-chan string, <-chan error) {
	wal.log.DebugContext(ctx, "starting WAL recovering")

	ctx, cancel := context.WithCancel(ctx)

	errs := make(chan error, 1)
	segments, streamErrs := wal.ss.Stream(ctx)

	queries := make(chan string)
	go func() (err error) {
		wal.log.DebugContext(ctx, " segments processing started")

		defer close(queries)
		defer close(errs)

		defer func() {
			if err != nil {
				errs <- err
				wal.log.ErrorContext(ctx, "wal recovering failed", sl.Err(err))
			} else {
				wal.log.InfoContext(ctx, "wal recovering finished")
			}
		}()

		defer cancel()

		for {
			select {
			case err, ok := <-streamErrs:
				if ok {
					wal.log.ErrorContext(ctx, "segments processing failed", sl.Err(err))
					err = newSegmentsProcessingError(err)
					return err
				}
			default:
			}

			select {
			case <-ctx.Done():
				wal.log.WarnContext(ctx, "segments processing canceled", sl.Err(ctx.Err()))
				return ctx.Err()
			case err, ok := <-streamErrs:
				if ok {
					wal.log.ErrorContext(ctx, "segments processing failed", sl.Err(err))
					err = newSegmentsProcessingError(err)
					return err
				}
			case s, ok := <-segments:
				if !ok {
					wal.log.InfoContext(ctx, "segments processing finished")
					return nil
				}

				wal.mx.Lock()
				olgSegment := wal.currentSegment
				wal.currentSegment = s
				wal.mx.Unlock()

				if olgSegment != nil {
					_ = olgSegment.Close()
				}

				name := s.Name()
				log := wal.log.With(slog.String("segment_name", name))
				log.DebugContext(ctx, "segment processing started")

				size, err := s.Size()
				if err != nil {
					log.ErrorContext(ctx, "could not get segment size", sl.Err(err))
					err = newSegmentsProcessingError(err, name)
					return err
				}

				log.DebugContext(ctx, "segment size retrieved", slog.Uint64("segment_size", uint64(size)))

				data := make([]byte, size)
				if _, err = s.Read(data); err != nil {
					log.ErrorContext(ctx, "could not read segment", sl.Err(err))
					err = newSegmentsProcessingError(err, name)
					return err
				}
				log.DebugContext(ctx, "read segment")

				qs := strings.Split(string(data), "\n")
				log.DebugContext(ctx, "read segment queries", slog.Int("queries_count", len(qs)))

				for _, q := range qs {
					if q != "" {
						queries <- q
					}
				}

				log.DebugContext(ctx, "segment processing finished")
			}
		}
	}()

	wal.log.InfoContext(ctx, "WAL recovering started")

	return queries, errs
}
