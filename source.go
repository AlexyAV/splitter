package splitter

import (
	"fmt"
	"mime"
	"net/url"
)

// SourceError represent error message and context for target source.
type SourceError struct {
	context string
	err     error
}

func (pr *SourceError) Error() string {
	return fmt.Sprintf("splitter: source: %s: %v", pr.context, pr.err)
}

// TODO add description
type Source struct {
	Path   *url.URL
	Size   int
	Ext    string
	client HTTPClient
}

// NewSource creates new Source instance.
func NewSource(source *url.URL, client HTTPClient) (*Source, error) {
	var err error

	s := &Source{Path: source, client: client}
	err = s.enrichSourceInfo()

	return s, err
}

// TODO add description
// fetchContentLength issues a HEAD to the provided source and
// retrieves content length. Converts content length to int value.
func (s *Source) enrichSourceInfo() error {
	headResponse, err := s.client.Head(s.Path.String())
	if err != nil {
		return &SourceError{
			context: "cannot fetch source info",
			err:     err,
		}
	}

	s.Size = int(headResponse.ContentLength)
	if s.Size == 0 {
		return &SourceError{context: "cannot fetch content length"}
	}

	ct, err := mime.ExtensionsByType(headResponse.Header.Get("Content-Type"))
	if len(ct) == 0 || err != nil {
		return &SourceError{
			context: "cannot fetch content type",
			err:     err,
		}
	}

	s.Ext = ct[0]

	return nil
}
