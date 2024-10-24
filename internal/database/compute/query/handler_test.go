package query

import (
	"context"
	"errors"
	enginemock "github.com/akimsavvin/sakv/internal/database/storage/engine/mock"
	slogmock "github.com/samber/slog-mock"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"strings"
	"testing"
)

type HandleQuerySuite struct {
	suite.Suite
	logMock *slog.Logger
	eMock   *enginemock.Engine
}

func (suite *HandleQuerySuite) SetupSuite() {
	suite.logMock = slog.New(new(slogmock.MockHandler))
}

func (suite *HandleQuerySuite) SetupTest() {
	suite.eMock = enginemock.NewEngine()
}

func (suite *HandleQuerySuite) TestGETSuccess() {
	// Arrange
	const query = "GET key"
	ctx := context.Background()
	suite.eMock.On("GET", ctx, "key").Return("value", nil)

	h := NewHandler(suite.logMock, nil, suite.eMock)

	// Act
	res := h.HandleQuery(ctx, query)

	// Assert
	suite.True(strings.HasPrefix(res, "[OK]"))
	suite.True(strings.HasSuffix(res, "value"))
	suite.eMock.AssertCalled(suite.T(), "GET", ctx, "key")
}

func (suite *HandleQuerySuite) TestGETError() {
	// Arrange
	const query = "GET key"
	ctx := context.Background()
	suite.eMock.On("GET", ctx, "key").Return("", errors.ErrUnsupported)

	h := NewHandler(suite.logMock, nil, suite.eMock)

	// Act
	res := h.HandleQuery(ctx, query)

	// Assert
	suite.True(strings.HasPrefix(res, "[ERROR]"))
	suite.eMock.AssertCalled(suite.T(), "GET", ctx, "key")
}

func (suite *HandleQuerySuite) TestSETSuccess() {
	// Arrange
	const query = "SET key value"
	ctx := context.Background()
	suite.eMock.On("SET", ctx, "key", "value").Return(nil)

	h := NewHandler(suite.logMock, nil, suite.eMock)

	// Act
	res := h.HandleQuery(ctx, query)

	// Assert
	suite.Equal("[OK]", res)
	suite.eMock.AssertCalled(suite.T(), "SET", ctx, "key", "value")
}

func (suite *HandleQuerySuite) TestSETError() {
	// Arrange
	const query = "SET key value"
	ctx := context.Background()
	suite.eMock.On("SET", ctx, "key", "value").Return(errors.ErrUnsupported)

	h := NewHandler(suite.logMock, nil, suite.eMock)

	// Act
	res := h.HandleQuery(ctx, query)

	// Assert
	suite.True(strings.HasPrefix(res, "[ERROR]"))
	suite.eMock.AssertCalled(suite.T(), "SET", ctx, "key", "value")
}

func (suite *HandleQuerySuite) TestDELSuccess() {
	// Arrange
	const query = "DEL key"
	ctx := context.Background()
	suite.eMock.On("DEL", ctx, "key").Return(nil)

	h := NewHandler(suite.logMock, nil, suite.eMock)

	// Act
	res := h.HandleQuery(ctx, query)

	// Assert
	suite.Equal("[OK]", res)
	suite.eMock.AssertCalled(suite.T(), "DEL", ctx, "key")
}

func (suite *HandleQuerySuite) TestDELError() {
	// Arrange
	const query = "DEL key"
	ctx := context.Background()
	suite.eMock.On("DEL", ctx, "key").Return(errors.ErrUnsupported)

	h := NewHandler(suite.logMock, nil, suite.eMock)

	// Act
	res := h.HandleQuery(ctx, query)

	// Assert
	suite.True(strings.HasPrefix(res, "[ERROR]"))
	suite.eMock.AssertCalled(suite.T(), "DEL", ctx, "key")
}

func TestHandler_HandleQuery(t *testing.T) {
	suite.Run(t, new(HandleQuerySuite))
}
