package splitter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNextRange(t *testing.T) {
	rb := NewRangeBuilder(100, 10, 0)
	lastRangeHeader := ""

	for i := 0; i < 10; i++ {
		r, err := rb.NextRange()

		assert.NoError(t, err)
		assert.NotEqual(t, r.Start, r.End)

		lastRangeHeader = r.BuildRangeHeader()
	}

	assert.Equal(t, "bytes=90-99", lastRangeHeader)

	_, err := rb.NextRange()
	assert.EqualError(t, err, "ErrOutOfRange")
}
