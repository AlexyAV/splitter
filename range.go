package splitter

import (
	"errors"
	"fmt"
	"strconv"
)

// ErrOutOfRange is the error returned by NextRange when no more range is available.
var ErrOutOfRange = errors.New("ErrOutOfRange")

// DownloadRange is a basic data structure for storing bytes range data.
// Min Start value is 0 and max End value is file size.
type DownloadRange struct {
	Start, End int
}

// BuildRangeHeader builds bytes range for http Range header
func (dr *DownloadRange) BuildRangeHeader() string {
	return fmt.Sprintf(
		"bytes=%s-%s",
		strconv.Itoa(dr.Start),
		strconv.Itoa((dr.End)-1),
	)
}

// A RangeBuilder allows to iterate over convent length and split it on separate
// DownloadRange on each iteration.
type RangeBuilder struct {
	contentLen, rangeSize, remainder, start, end int
}

// NewRangeBuilder creates an instance of RangeBuilder based on total length
// and chunks count into which total length will be split.
func NewRangeBuilder(length, chunkCount, offset int) *RangeBuilder {
	adjLength := length - offset
	remainder := adjLength % chunkCount
	rangeSize := (adjLength - remainder) / chunkCount

	return &RangeBuilder{
		contentLen: length,
		remainder:  remainder,
		rangeSize:  rangeSize,
		start:      offset,
		end:        offset,
	}
}

// NextRange iterates over content length and creates new DownloadRange instance
// on each iteration. If end of range was reached ErrOutOfRange error will be returned.
func (rb *RangeBuilder) NextRange() (DownloadRange, error) {
	if rb.end == rb.contentLen {
		return DownloadRange{}, ErrOutOfRange
	}

	chunkSize := rb.rangeSize

	if rb.end == rb.start {
		chunkSize = rb.rangeSize + rb.remainder
	}

	rb.start = rb.end
	rb.end = rb.end + chunkSize

	return DownloadRange{rb.start, rb.end}, nil
}
