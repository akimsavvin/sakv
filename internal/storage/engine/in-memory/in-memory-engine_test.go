package inmemory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine_GET(t *testing.T) {
	// Arrange
	e := NewEngine()
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
	e := NewEngine()
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
	e := NewEngine()
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
