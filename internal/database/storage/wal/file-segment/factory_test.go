package filesegment

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFactory_Create(t *testing.T) {
	// Arrange
	fsf := NewFactory()

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
