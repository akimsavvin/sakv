package filesegment

import (
	slogmock "github.com/samber/slog-mock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestFactory_Create(t *testing.T) {
	// Arrange
	log := slog.New(new(slogmock.MockHandler))
	fsf := NewFactory(log)

	// Act
	s, err := fsf.Create("./")
	var fs *FileSegment

	// Assert
	assert.NoError(t, err)
	if assert.NotNil(t, s) {
		assert.IsType(t, fs, s)
		assert.NoError(t, s.Close())

		// Cleanup
		_ = os.Remove(s.Name())
	}
}
