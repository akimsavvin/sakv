package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/akimsavvin/sakv/internal/database/config"
	"github.com/akimsavvin/sakv/internal/database/storage/wal"
	walmock "github.com/akimsavvin/sakv/internal/database/storage/wal/mock"
	"github.com/brianvoe/gofakeit"
	slogmock "github.com/samber/slog-mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"testing"
)

func randomSegmentName() string {
	return fmt.Sprintf("segment_%d", gofakeit.Int64())
}

type NewSuite struct {
	suite.Suite
	logMock *slog.Logger
}

func (suite *NewSuite) SetupSuite() {
	suite.logMock = slog.New(new(slogmock.MockHandler))
}

func (suite *NewSuite) TestInvalidFlushingTimeout() {
	// Arrange
	ctx := context.Background()
	sf := walmock.NewSegmentFactory()
	ss := walmock.NewSegmentsStreamer()

	cfg := config.WAL{
		FlushingBatchTimeout: "sdfjkls3",
	}

	// Act
	inst, err := wal.New(ctx, suite.logMock, cfg, sf, ss)

	// Assert
	suite.Nil(inst)
	suite.Error(err)
}

func (suite *NewSuite) TestInvalidSegmentSize() {
	// Arrange
	ctx := context.Background()
	sf := walmock.NewSegmentFactory()
	ss := walmock.NewSegmentsStreamer()

	cfg := config.WAL{
		FlushingBatchTimeout: "10ms",
		MaxSegmentSize:       "gdsf43g",
	}

	// Act
	inst, err := wal.New(ctx, suite.logMock, cfg, sf, ss)

	// Assert
	suite.Nil(inst)
	suite.Error(err)
}

func (suite *NewSuite) TestSuccess() {
	// Arrange
	ctx := context.Background()
	sf := walmock.NewSegmentFactory()
	ss := walmock.NewSegmentsStreamer()

	cfg := config.WAL{
		FlushingBatchSize:    10,
		FlushingBatchTimeout: "32ms",
		MaxSegmentSize:       "1mb",
		DataDirectory:        "./test-data/wal",
	}

	// Act
	inst, err := wal.New(ctx, suite.logMock, cfg, sf, ss)

	// Assert
	suite.NotNil(inst)
	suite.NoError(err)
}

func TestNew(t *testing.T) {
	suite.Run(t, new(NewSuite))
}

type RecoverSuite struct {
	suite.Suite
	logMock *slog.Logger
	cfg     config.WAL
	sfMock  *walmock.SegmentFactory
	ssMock  *walmock.SegmentsStreamer
}

func (suite *RecoverSuite) SetupSuite() {
	suite.logMock = slog.New(new(slogmock.MockHandler))
	suite.cfg = config.WAL{
		FlushingBatchSize:    10,
		FlushingBatchTimeout: "32ms",
		MaxSegmentSize:       "1mb",
		DataDirectory:        "./test-data/wal",
	}
}

func (suite *RecoverSuite) SetupTest() {
	suite.sfMock = walmock.NewSegmentFactory()
	suite.ssMock = walmock.NewSegmentsStreamer()
}

func (suite *RecoverSuite) TestSuccess() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inst, _ := wal.New(ctx, suite.logMock, suite.cfg, suite.sfMock, suite.ssMock)

	segmentsCh := make(chan wal.Segment)
	errsCh := make(chan error)
	suite.ssMock.
		On("Stream", mock.AnythingOfType("*context.cancelCtx")).
		Return(segmentsCh, errsCh)

	queries := make([]string, 0)
	segments := make([]wal.Segment, 0)

	for range gofakeit.Number(2, 16) {
		s := walmock.NewSegment(randomSegmentName())
		var contentBuilder bytes.Buffer
		for range gofakeit.Number(0, 16) {
			query := gofakeit.Word()
			queries = append(queries, query)
			contentBuilder.Write([]byte(query + "\n"))
		}

		s.SetContent(contentBuilder.Bytes())
		segments = append(segments, s)
	}

	// Act
	go func() {
		for _, s := range segments {
			segmentsCh <- s
		}
		close(segmentsCh)
		close(errsCh)
	}()

	resQueriesCh, resErrsCh := inst.Recover(ctx)

	// Assert
	resQueries := make([]string, 0)
	for query := range resQueriesCh {
		resQueries = append(resQueries, query)
	}

	suite.NoError(<-resErrsCh)
	suite.Equal(queries, resQueries)
}

func (suite *RecoverSuite) TestError() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inst, _ := wal.New(ctx, suite.logMock, suite.cfg, suite.sfMock, suite.ssMock)

	segmentsCh := make(chan wal.Segment)
	errsCh := make(chan error, 1)
	suite.ssMock.
		On("Stream", mock.AnythingOfType("*context.cancelCtx")).
		Return(segmentsCh, errsCh)

	// Act
	go func() {
		errsCh <- errors.ErrUnsupported
		close(segmentsCh)
		close(errsCh)
	}()

	resQueriesCh, resErrsCh := inst.Recover(ctx)

	// Assert
	resQueries := make([]string, 0)
	for query := range resQueriesCh {
		resQueries = append(resQueries, query)
	}

	suite.Error(<-resErrsCh)
	suite.Empty(resQueries)
}

func (suite *RecoverSuite) TestCancel() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())

	inst, _ := wal.New(ctx, suite.logMock, suite.cfg, suite.sfMock, suite.ssMock)

	segmentsCh := make(chan wal.Segment)
	errsCh := make(chan error)
	suite.ssMock.
		On("Stream", mock.AnythingOfType("*context.cancelCtx")).
		Return(segmentsCh, errsCh)

	// Act
	resQueriesCh, resErrsCh := inst.Recover(ctx)
	cancel()

	// Assert
	resQueries := make([]string, 0)
	for query := range resQueriesCh {
		resQueries = append(resQueries, query)
	}

	suite.Error(<-resErrsCh)
	suite.Empty(resQueries)
}

func TestWAL_Recover(t *testing.T) {
	suite.Run(t, new(RecoverSuite))
}

type WriteSuite struct {
	suite.Suite
	logMock *slog.Logger
	cfg     config.WAL
	sfMock  *walmock.SegmentFactory
	ssMock  *walmock.SegmentsStreamer
}

func (suite *WriteSuite) SetupSuite() {
	suite.logMock = slog.New(new(slogmock.MockHandler))
	suite.cfg = config.WAL{
		FlushingBatchSize:    10,
		FlushingBatchTimeout: "32ms",
		MaxSegmentSize:       "1mb",
		DataDirectory:        "./test-data/wal",
	}
}

func (suite *WriteSuite) SetupTest() {
	suite.sfMock = walmock.NewSegmentFactory()
	suite.ssMock = walmock.NewSegmentsStreamer()
}

func (suite *WriteSuite) TestSuccess() {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inst, _ := wal.New(ctx, suite.logMock, suite.cfg, suite.sfMock, suite.ssMock)

	s := walmock.NewSegment(randomSegmentName())
	s.SetContent([]byte(gofakeit.Sentence(10)))
	s.On("Write", mock.AnythingOfType("[]uint8")).Return(0, nil)

	suite.sfMock.On("Create", mock.AnythingOfType("string")).Return(s, nil)

	// Act
	err := inst.Write(ctx, gofakeit.BuzzWord())

	// Assert
	suite.NoError(err)
}

func TestWAL_Write(t *testing.T) {
	suite.Run(t, new(WriteSuite))
}
