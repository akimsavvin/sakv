package walmock

import (
	"context"
	"github.com/akimsavvin/sakv/internal/database/storage/wal"
	"github.com/stretchr/testify/mock"
)

type SegmentsStreamer struct {
	mock.Mock
}

func NewSegmentsStreamer() *SegmentsStreamer {
	return new(SegmentsStreamer)
}

func (ss *SegmentsStreamer) Stream(ctx context.Context) (<-chan wal.Segment, <-chan error) {
	args := ss.Called(ctx)
	return args.Get(0).(chan wal.Segment), args.Get(1).(chan error)
}
