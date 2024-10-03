package inmemory

import (
	"context"
	slogmock "github.com/samber/slog-mock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestEngine_GET(t *testing.T) {
	// Arrange
	logmock := slog.New(new(slogmock.MockHandler))
	e := NewEngine(logmock)
	e.data = map[string]string{"key": "value"}
	ctx := context.Background()

	// Act
	val, err := e.GET(ctx, "key")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "value", val)
}

func TestEngine_SET(t *testing.T) {
	// Arrange
	logmock := slog.New(new(slogmock.MockHandler))
	e := NewEngine(logmock)
	ctx := context.Background()

	// Act
	err := e.SET(ctx, "key", "value")
	val, ok := e.data["key"]

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "value", val)
	assert.True(t, ok)
}

func TestEngine_DEL(t *testing.T) {
	// Arrange
	logmock := slog.New(new(slogmock.MockHandler))
	e := NewEngine(logmock)
	e.data = map[string]string{"key": "value"}
	ctx := context.Background()

	// Act
	err := e.DEL(ctx, "key")
	val, ok := e.data["key"]

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, val)
	assert.False(t, ok)
}
