package splitter

import (
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Head(url string) (*http.Response, error)
}
