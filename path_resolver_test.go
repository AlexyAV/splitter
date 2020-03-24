package splitter

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"testing"
)

func TestPathResolver_PathInfo(t *testing.T) {
	dir, f := initTmpStorage()
	defer os.RemoveAll(dir)
	prepareHttpClientResp()

	pr := PathResolver{
		Source: "http://test-url.com/image/source.jpg",
		Dest:   f.Name(),
		client: &mockClient{},
	}

	_, err := pr.PathInfo()
	assert.NoError(t, err)
}

func TestPathResolver_PathInfoSourceError(t *testing.T) {
	pr := PathResolver{
		Source: "image_source.jpg",
		Dest:   os.TempDir(),
		client: &mockClient{},
	}

	_, err := pr.PathInfo()
	assert.Error(t, err)
}

func TestPathResolver_PathInfoSourceInfoError(t *testing.T) {
	GetHeadFunc = func(url string) (*http.Response, error) {
		return nil, errors.New("http error")
	}

	pr := PathResolver{
		Source: "http://test-url.com/image/source.jpg",
		Dest:   os.TempDir(),
		client: &mockClient{},
	}

	_, err := pr.PathInfo()
	assert.Error(t, err)
}

func TestPathResolver_PathInfoDestError(t *testing.T) {
	dir, f := initTmpStorage()
	defer os.RemoveAll(dir)

	_ = os.Chmod(f.Name(), 0000)
	prepareHttpClientResp()

	pr := PathResolver{
		Source: "http://test-url.com/image/source.jpg",
		Dest:   f.Name(),
		client: &mockClient{},
	}

	_, err := pr.PathInfo()
	assert.Error(t, err)
}

func TestResolveSource(t *testing.T) {
	testSource := "http://test-url.com/image/source.jpg"

	pr := PathResolver{
		Source: testSource,
		Dest:   os.TempDir(),
		client: nil,
	}

	s, err := pr.resolveSource()
	assert.NoError(t, err)
	assert.Equal(t, testSource, s.String())
}

func TestResolveSourceError(t *testing.T) {
	pr := PathResolver{
		Source: "image_source.jpg",
		Dest:   os.TempDir(),
		client: nil,
	}

	_, err := pr.resolveSource()
	assert.EqualError(
		t,
		err,
		"splitter: path resolver: invalid source path: parse image_source.jpg: invalid URI for request",
	)
}

func TestResolveDest(t *testing.T) {
	dir, f := initTmpStorage()
	defer os.RemoveAll(dir)

	testURL, _ := url.ParseRequestURI("http://source.com/file.txt")
	noExtURL, _ := url.ParseRequestURI("http://source.com/test")

	destTests := []struct {
		source           Source
		dest, resultDest string
		valid            bool
	}{
		{
			Source{testURL, 100, ".txt", nil},
			dir,
			path.Join(dir, "file.txt"),
			true,
		},
		{
			Source{noExtURL, 100, ".txt", nil},
			f.Name(),
			f.Name(),
			true,
		},
		{
			Source{noExtURL, 100, ".txt", nil},
			dir,
			path.Join(dir, "test.txt"),
			true,
		},
		{
			Source{noExtURL, 100, ".txt", nil},
			"fakeDest",
			"fakeDest",
			false,
		},
	}
	prepareHttpClientResp()

	for _, pathInfo := range destTests {
		pr := NewPathResolver(
			pathInfo.source.Path.String(),
			pathInfo.dest,
			&mockClient{},
		)

		d, err := pr.resolveDest(&pathInfo.source)
		if err != nil {
			if pathInfo.valid {
				t.Errorf(
					"Unexpected error: \"%v\" for valid destination %s",
					err,
					pathInfo.dest,
				)
			}
			return
		}

		if !pathInfo.valid {
			t.Errorf(
				"Error expected for invalid destination %s",
				pathInfo.source.Path.String(),
			)
		}

		assert.Equal(t, d.Name(), pathInfo.resultDest)
	}
}

func TestResolveDestError(t *testing.T) {
	dir, f := initTmpStorage()
	defer os.RemoveAll(dir)

	_ = os.Chmod(f.Name(), 0000)

	testURL, _ := url.ParseRequestURI("http://source.com/test")
	pr := NewPathResolver(testURL.String(), f.Name(), nil)
	_, err := pr.resolveDest(&Source{testURL, 100, ".txt", nil})

	assert.EqualError(
		t,
		err,
		fmt.Sprintf(
			"splitter: path resolver: cannot open file - %s: open %s: permission denied",
			f.Name(),
			f.Name(),
		),
	)
}

func TestResolveDestSourceError(t *testing.T) {
	dir, _ := initTmpStorage()
	defer os.RemoveAll(dir)

	f, err := os.Create(path.Join(dir, "file.txt"))
	if err != nil {
		log.Fatal(err)
	}
	_ = os.Chmod(f.Name(), 0000)

	testURL, _ := url.ParseRequestURI("http://source.com/file.txt")
	pr := NewPathResolver(testURL.String(), dir, nil)
	_, err = pr.resolveDest(&Source{testURL, 100, ".txt", nil})

	assert.EqualError(
		t,
		err,
		fmt.Sprintf(
			"splitter: path resolver: cannot resolve destination source: open %s: permission denied",
			f.Name(),
		),
	)
}

func prepareHttpClientResp() {
	GetHeadFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode:    200,
			Header:        http.Header{"Content-Type": []string{"image/jpeg"}},
			ContentLength: 100,
		}, nil
	}
}

func initTmpStorage() (string, *os.File) {
	dir, err := ioutil.TempDir("", "splitter")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(path.Join(dir, "dest_file.txt"))
	if err != nil {
		log.Fatal(err)
	}

	return dir, f
}
