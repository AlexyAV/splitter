package splitter

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

// PathResolverError represent error message and context for path resolver.
type PathResolverError struct {
	context string
	err     error
}

func (pr *PathResolverError) Error() string {
	return fmt.Sprintf("splitter: path resolver: %s: %v", pr.context, pr.err)
}

// A PathInfo is simple storage for source and destination paths.
type PathInfo struct {
	Source *Source
	Dest   *os.File
}

// A PathResolver allows to resolve source path and destination path. Provides
// single public method to resolve both values as *url.URL and *os.File
// respectively.
type PathResolver struct {
	Source string
	Dest   string
	client HTTPClient
}

// NewPathResolver creates new PathResolver instance.
func NewPathResolver(source string, dest string, client HTTPClient) *PathResolver {
	destAbs, _ := filepath.Abs(dest)
	return &PathResolver{Source: source, Dest: destAbs, client: client}
}

// PathInfo resolves provided source and dest path and creates PathInfo instance
// with resolved source as *url.URL and dest as *os.File.
func (pr *PathResolver) PathInfo() (*PathInfo, error) {
	rawSource, err := pr.resolveSource()
	if err != nil {
		return nil, err
	}

	s, err := NewSource(rawSource, pr.client)
	if err != nil {
		return nil, err
	}

	d, err := pr.resolveDest(s)
	if err != nil {
		return nil, err
	}

	return &PathInfo{Source: s, Dest: d}, nil
}

// resolveSource resolves provided source path and create *url.URL instance
// or return error in case of invalid path.
func (pr *PathResolver) resolveSource() (*url.URL, error) {
	uri, err := url.ParseRequestURI(pr.Source)
	if err != nil {
		return nil, &PathResolverError{context: "invalid source path", err: err}
	}

	return uri, nil
}

// resolveDest resolves provided destination path and create *os.File instance
// or return error in case of invalid path or lack of permissions. It accepts
// full path with file extension as well as dir path. In last case file name
// from source will be used.
func (pr *PathResolver) resolveDest(s *Source) (*os.File, error) {
	if _, err := os.Stat(pr.Dest); os.IsNotExist(err) {
		return nil, err
	}

	if extProvided(pr.Dest) {
		f, err := os.OpenFile(pr.Dest, os.O_RDWR, 0666)
		if err != nil {
			return nil, &PathResolverError{
				context: fmt.Sprintf("cannot open file - %s", pr.Dest),
				err:     err,
			}
		}

		return f, nil
	}

	basePath := path.Base(s.Path.Path)

	if !extProvided(pr.Source) {
		basePath += s.Ext
	}

	d, err := os.Create(path.Join(pr.Dest, basePath))
	if err != nil {
		return nil, &PathResolverError{
			context: "cannot resolve destination source",
			err:     err,
		}
	}

	return d, nil
}

// extProvided checks if the path contains an extension part.
func extProvided(p string) bool {
	return len(filepath.Ext(filepath.Base(p))) != 0
}
