package splitter

import (
	"net/http"
)

// HTTPClient is the interface that wraps the basic http methods from standard
// http.Client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Head(url string) (*http.Response, error)
}
