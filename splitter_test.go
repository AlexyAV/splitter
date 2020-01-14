package splitter

import (
	"context"
	"fmt"
	"testing"
)

func TestNewChunkRequest(t *testing.T) {
	pr := NewPathResolver("https://picsum.photos/200", "/tmp")
	pi, err := pr.PathInfo()
	if err != nil {
		t.Fatal(err)
	}

	s := NewSplitter(pi, 10, context.Background())
	dr := DownloadRange{Start: 10, End: 20}
	expected := "bytes=10-19"

	r, err := s.newChunkRequest(dr)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(r.Header.Get("Range"))

	if r.Header.Get("Range") != expected {
		t.Errorf("Range - %s, expected - %s.", r.Header.Get("Range"), expected)
	}
}
