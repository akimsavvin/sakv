package memory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSize(t *testing.T) {
	// Arrange
	type ParseSizeSuite struct {
		SizeStr  string
		Expected uint
		Error    bool
	}

	suites := []*ParseSizeSuite{
		{
			SizeStr:  "128",
			Expected: 128,
			Error:    false,
		},
		//{
		//	SizeStr:  "256b",
		//	Expected: 256,
		//	Error:    false,
		//},
		{
			SizeStr:  "8KB",
			Expected: 8192,
			Error:    false,
		},
		{
			SizeStr:  "2mb",
			Expected: 2097152,
			Error:    false,
		},
		{
			SizeStr:  "1Gb",
			Expected: 1073741824,
			Error:    false,
		},
		{
			SizeStr:  "3tB",
			Expected: 3298534883328,
			Error:    false,
		},
		{
			SizeStr:  "8lb",
			Expected: 0,
			Error:    true,
		},
		{
			SizeStr:  "2pb",
			Expected: 0,
			Error:    true,
		},
	}

	for _, suite := range suites {
		// Act
		res, err := ParseSize(suite.SizeStr)

		// Assert
		assert.Equal(t, suite.Expected, res)

		if suite.Error {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
