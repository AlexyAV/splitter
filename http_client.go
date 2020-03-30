package splitter

import (
	"net/http"
)

// HTTPClient is the interface that wraps the basic http methods from standard
// http.Client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	Get(url string) (resp *http.Response, err error)
}
