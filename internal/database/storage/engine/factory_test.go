package engine

import (
	"github.com/akimsavvin/sakv/internal/database/config"
	inmemory "github.com/akimsavvin/sakv/internal/database/storage/engine/in-memory"
	slogmock "github.com/samber/slog-mock"
	"github.com/stretchr/testify/suite"
	"log/slog"
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
	cfg := config.Engine{Type: InMemory}
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

	// Act
	action := func() {
		f.CreateEngine(cfg)
	}

	// Assert
	suite.Panics(action)
}

func TestFactory_CreateEngine(t *testing.T) {
	suite.Run(t, new(CreateEngineSuite))
}
