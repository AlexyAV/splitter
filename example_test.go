package splitter_test

import (
    "context"
    "fmt"
    "log"
    splitter "splitter/src"
)

func Example() {
    // With absolute destination path
    pr := splitter.NewPathResolver("https://picsum.photos/200", "/tmp/")
    pi, err := pr.PathInfo()
    if err != nil {
        log.Fatal(err)
    }

    // Create Splitter instance with new PathInfo and 10 chunks
    s := splitter.NewSplitter(pi, 10, context.Background())

    // Start file download
    err = s.Download()
    if err != nil {
        log.Fatal(err)
    }
}

func ExampleNewPathResolver() {
    // With absolute destination path
    pr := splitter.NewPathResolver("https://picsum.photos/200", "/tmp/")

    // Current directory will be used as destination
    // splitter.NewPathResolver("https://picsum.photos/200", ".")

    pi, err := pr.PathInfo()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%T\n", pi.Source)
    fmt.Printf("%T\n", pi.Dest)
    // Output:
    // *url.URL
    // *os.File
}

func ExampleRangeBuilder_NextRange() {
    contentLength := 55
    chunkCount := 6
    rb := splitter.NewRangeBuilder(contentLength, chunkCount)

    for {
        r, err := rb.NextRange()

        if err == splitter.EOR {
            break
        }

        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf(
            "Start - %d; End %d; Range header - %s\n",
            r.Start,
            r.End,
            r.BuildRangeHeader(),
        )
    }
    // Output:
    // Start - 0; End 10; Range header - bytes=0-9
    // Start - 10; End 19; Range header - bytes=10-18
    // Start - 19; End 28; Range header - bytes=19-27
    // Start - 28; End 37; Range header - bytes=28-36
    // Start - 37; End 46; Range header - bytes=37-45
    // Start - 46; End 55; Range header - bytes=46-54
}
