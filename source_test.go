package splitter

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"testing"
)

var (
	GetDoFunc  func(req *http.Request) (*http.Response, error)
	GetGetFunc func(url string) (resp *http.Response, err error)
)

type mockClient struct {
	mock.Mock
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func (m *mockClient) Get(url string) (resp *http.Response, err error) {
	return GetGetFunc(url)
}

func TestNewSource(t *testing.T) {
	httpClient := &mockClient{}
	testUrl, _ := url.Parse("http://test-url.com/image/source.jpg")

	GetGetFunc = func(url string) (resp *http.Response, err error) {
		return &http.Response{
			StatusCode:    200,
			Header:        http.Header{"Content-Type": []string{"image/jpeg"}},
			ContentLength: 100,
		}, nil
	}

	s, err := NewSource(testUrl, httpClient)

	assert.Nil(t, err)
	assert.Equal(t, 100, s.Size)
	assert.Equal(t, testUrl, s.Path)
	assert.Contains(t, []string{".jpeg", ".jpg"}, s.Ext)
}

func TestNewSourceRequestError(t *testing.T) {
	httpClient := &mockClient{}
	testUrl, _ := url.Parse("http://test-url.com/image/source.jpg")

	GetGetFunc = func(url string) (*http.Response, error) {
		return nil, errors.New("request failed")
	}

	_, err := NewSource(testUrl, httpClient)

	assert.EqualError(
		t,
		err,
		"splitter: source: cannot fetch source info: request failed",
	)
}

func TestNewSourceContentLengthError(t *testing.T) {
	httpClient := &mockClient{}
	testUrl, _ := url.Parse("http://test-url.com/image/source.jpg")

	GetGetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Header:     http.Header{"Content-Type": []string{"image/jpeg"}},
		}, nil
	}

	_, err := NewSource(testUrl, httpClient)

	assert.EqualError(
		t,
		err,
		"splitter: source: cannot fetch content length: <nil>",
	)
}

func TestNewSourceContentTypeError(t *testing.T) {
	httpClient := &mockClient{}
	testUrl, _ := url.Parse("http://test-url.com/image/source.jpg")

	GetGetFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode:    200,
			ContentLength: 100,
		}, nil
	}

	_, err := NewSource(testUrl, httpClient)

	assert.EqualError(
		t,
		err,
		"splitter: source: cannot fetch content type: mime: no media type",
	)
}
