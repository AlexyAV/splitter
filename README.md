[![Build Status](https://travis-ci.org/AlexyAV/splitter.svg?branch=master)](https://travis-ci.org/AlexyAV/splitter)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlexyAV/splitter)](https://goreportcard.com/report/github.com/AlexyAV/splitter)
[![GoDoc](https://godoc.org/github.com/AlexyAV/splitter?status.svg)](https://godoc.org/github.com/AlexyAV/splitter)
# Splitter 
Simple Go package for file chunk async download that made just for fun. Choose a file, slice it up, download it. Please do not set 20 chunks for a 10kb file.

# Installation
Fetch package
```
go get github.com/AlexyAV/splitter
```
<pre lang="go">
import (
  splitter "github.com/AlexyAV/splitter"
)
</pre>

# Usage
<pre lang="go">
func main() {
  // With absolute destination path
  pr := splitter.NewPathResolver("https://picsum.photos/200", "/tmp/")
  pi, err := pr.PathInfo()
  if err != nil {
    log.Fatal(err)
  }

  // Create Splitter instance with new PathInfo and 10 chunks
  s := splitter.NewSplitter(context.Background(), pi, 10)

  // Start file download
  err = s.Download()
  if err != nil {
    log.Fatal(err)
  }
}
</pre>
