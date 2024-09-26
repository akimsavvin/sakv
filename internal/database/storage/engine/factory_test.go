package engine

import (
	slogmock "github.com/samber/slog-mock"
	"github.com/stretchr/testify/suite"
	"log/slog"
	"sakv/internal/database/config"
	inmemory "sakv/internal/database/storage/engine/in-memory"
	"testing"
)

type CreateEngineSuite struct {
	suite.Suite

	logmock *slog.Logger
}

func (suite *CreateEngineSuite) SetupSuite() {
	suite.logmock = slog.New(new(slogmock.MockHandler))
}

func (suite *CreateEngineSuite) SetupTest() {}

func (suite *CreateEngineSuite) TestInMemory() {
	// Arrange
	cfg := config.Engine{Type: "in_memory"}
	f := NewFactory(suite.logmock)

	// Act
	e := f.CreateEngine(cfg)
	var inMemoryE *inmemory.Engine

	// Assert
	suite.IsType(inMemoryE, e)
}

func (suite *CreateEngineSuite) TestUnknown() {
	// Arrange
	cfg := config.Engine{Type: "unknown"}
	f := NewFactory(suite.logmock)

	// Act & Assert
	suite.Panics(func() {
		f.CreateEngine(cfg)
	})
}

func TestFactory_CreateEngine(t *testing.T) {
	suite.Run(t, new(CreateEngineSuite))
}
