package splitter

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestSplitter_Download(t *testing.T) {
	dir, f := initTmpStorage()
	defer os.RemoveAll(dir)

	mockResponse := &http.Response{
		StatusCode:    200,
		ProtoMajor:    1,
		ProtoMinor:    0,
		Header:        http.Header{"Content-Type": []string{"text/plain"}},
		Body:          ioutil.NopCloser(strings.NewReader("abcdef")),
		ContentLength: 6,
	}

	GetGetFunc = func(url string) (resp *http.Response, err error) {
		return mockResponse, nil
	}

	GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return mockResponse, nil
	}

	pr := PathResolver{
		Source: "http://test-url.com/test/text",
		Dest:   f.Name(),
		client: &mockClient{},
	}

	pi, err := pr.PathInfo()
	assert.NoError(t, err)

	s := NewSplitter(context.Background(), pi, 1, &mockClient{})
	err = s.Download()
	assert.NoError(t, err)
}

func TestDownloadChunkError(t *testing.T) {
	s := splitterStub(nil)

	err := s.downloadChunk(DownloadRange{0, 6})
	assert.EqualError(
		t,
		err,
		"splitter: cannot prepare request: net/http: nil Context",
	)
}

func TestDownloadChunkReqError(t *testing.T) {
	GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("bad request")
	}

	s := splitterStub(context.Background())

	err := s.downloadChunk(DownloadRange{0, 6})
	assert.EqualError(
		t,
		err,
		"splitter: chunk download error: bad request",
	)
}

func TestWriteChunkError(t *testing.T) {
	GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode:    200,
			ProtoMajor:    1,
			ProtoMinor:    0,
			Header:        http.Header{"Content-Type": []string{"text/plain"}},
			Body:          ioutil.NopCloser(strings.NewReader("abcdef")),
			ContentLength: 6,
		}, nil
	}

	s := splitterStub(context.Background())
	s.PI.Dest = nil

	err := s.downloadChunk(DownloadRange{0, 6})
	assert.EqualError(
		t,
		err,
		"splitter: error on writing data: invalid argument",
	)
}

func splitterStub(ctx context.Context) Splitter {
	testURL, _ := url.ParseRequestURI("http://source.com/file.txt")

	return Splitter{
		Ctx: ctx,
		PI: &PathInfo{
			Source: &Source{
				Path:   testURL,
				Size:   0,
				Ext:    "",
				client: nil,
			},
			Dest: nil,
		},
		client: &mockClient{},
	}
}
