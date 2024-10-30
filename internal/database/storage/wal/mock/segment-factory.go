package walmock

import (
	"github.com/akimsavvin/sakv/internal/database/storage/wal"
	"github.com/stretchr/testify/mock"
)

type SegmentFactory struct {
	mock.Mock
}

func NewSegmentFactory() *SegmentFactory {
	return new(SegmentFactory)
}

func (s *SegmentFactory) Create(dir string) (wal.Segment, error) {
	args := s.Called(dir)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	return args.Get(0).(wal.Segment), nil
}
