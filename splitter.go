// Package splitter implements functionality to download files by chunks
// asynchronously.
//
// The splitter package can handle only URL (RFC 3986) as source and save
// destination and file or directory. It won't create any new directory but
// file name only in case it was not provided.
//
// The number of chunks into which the file will be split is determined when
// the splitter instance is initialized.
package splitter

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
)

// Splitter allows to download source file by chunks asynchronously.
type Splitter struct {
	Ctx      context.Context
	PI       *PathInfo
	ChunkCnt int
	client   HTTPClient
}

type splitterError struct {
	context string
	err     error
}

func (se *splitterError) Error() string {
	return fmt.Sprintf("splitter: %s: %v", se.context, se.err)
}

// NewSplitter creates new Splitter instance.
func NewSplitter(ctx context.Context, pi *PathInfo, chunkCnt int, c HTTPClient) *Splitter {
	return &Splitter{Ctx: ctx, PI: pi, ChunkCnt: chunkCnt, client: c}
}

// Download initialize download process. It checks for content length and
// creates DownloadRange iterator. Each file's chunk will be downloaded
// asynchronously.
func (s *Splitter) Download() error {
	var g errgroup.Group
	rb := NewRangeBuilder(s.PI.Source.Size, s.ChunkCnt)

	for {
		nRange, err := rb.NextRange()
		if err == ErrOutOfRange {
			break
		}

		g.Go(func() error {
			return s.downloadChunk(nRange)
		})
	}

	err := g.Wait()

	return err
}

// downloadChunk creates and performs a new request for file chunk. The new
// request will fetch file's bytes range based on DownloadRange. After a
// successful response result will be written to dest path with an offset from
// DownloadRange.
func (s *Splitter) downloadChunk(dr DownloadRange) error {
	r, err := s.newChunkRequest(dr)
	if err != nil {
		return err
	}

	response, err := s.client.Do(r)
	if err != nil {
		return &splitterError{
			context: "chunk download error",
			err:     err,
		}
	}

	defer response.Body.Close()
	_, err = s.writeChunk(response.Body, int64(dr.Start))
	if err != nil {
		return err
	}

	return nil
}

// writeChunk writes result bytes range to destination file with specified offset.
func (s *Splitter) writeChunk(r io.Reader, offset int64) (int, error) {
	buf := make([]byte, 400)
	written := 0

	for {
		m, err := r.Read(buf[0:cap(buf)])
		if err == io.EOF {
			break
		}

		buf = buf[:m]

		_, err = s.PI.Dest.WriteAt(buf, offset)
		if err != nil {
			return 0, &splitterError{
				context: "error on writing data",
				err:     err,
			}
		}
		written += m
		offset += int64(m)
	}

	return written, nil
}

// newChunkRequest make new request to target source with provided DownloadRange
// info. Request will use "Range" header to download specific chunk of source.
func (s *Splitter) newChunkRequest(dr DownloadRange) (*http.Request, error) {
	request, err := http.NewRequestWithContext(
		s.Ctx,
		"GET",
		s.PI.Source.Path.String(),
		nil,
	)
	if err != nil {
		return nil, &splitterError{context: "cannot prepare request", err: err}
	}

	request.Header.Add("Range", dr.BuildRangeHeader())

	return request, nil
}
