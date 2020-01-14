package splitter

import (
    "io/ioutil"
    "log"
    "os"
    "path"
    "testing"
)

type PathResolverTest struct {
    source, dest string
    valid        bool
}

func TestPathResolver_PathInfo(t *testing.T) {
    dir, f := initTmpStorage()

    defer os.RemoveAll(dir)

    prt := PathResolverTest{"http://source.com", f.Name(), true}

    pr := NewPathResolver(prt.source, prt.dest)
    pi, err := pr.PathInfo()
    if err != nil {
        t.Errorf(
            "Unexpected error: \"%v\" for valid source %s and dest %s",
            err,
            prt.source,
            prt.dest,
        )
    }

    if pi.Source.String() != prt.source {
        t.Errorf(
            "Resolved source actual - %s, expected - %s.",
            pi.Source.String(),
            prt.source,
        )
    }

    if pi.Dest.Name() != f.Name() {
        t.Errorf(
            "Resolved destination actual - %s, expected - %s.",
            pi.Dest.Name(),
            f.Name(),
        )
    }
}

var sourceTests = []PathResolverTest{
    {"http://source.com", "/tmp", true},
    {"https://source.com", "/tmp", true},
    {"", "/tmp", false},
    {"source.com", "/tmp", false},
    {"https://source", "/tmp", false},
}

func TestResolveSource(t *testing.T) {
    for _, pathInfo := range sourceTests {
        pr := NewPathResolver(pathInfo.source, pathInfo.dest)
        s, err := pr.resolveSource()

        if err != nil {
            if pathInfo.valid {
                t.Errorf(
                    "Unexpected error: \"%v\" for valid source %s",
                    err,
                    pathInfo.source,
                )
            }
            return
        }

        if !pathInfo.valid {
            t.Errorf(
                "Error expected for invalid source %s",
                pathInfo.source,
            )
        }

        if s.String() != pathInfo.source {
            t.Errorf(
                "Resolved source actual - %s, expected - %s.",
                s.String(),
                pathInfo.source,
            )
        }
    }
}

func TestResolveDest(t *testing.T) {
    dir, f := initTmpStorage()
    defer os.RemoveAll(dir)

    destTests := []struct {
        source, dest, resultDest string
        valid                    bool
    }{
        {"http://source.com/file.txt", dir, path.Join(dir, "file.txt"), true},
        {"http://source.com", f.Name(), f.Name(), true},
        {"http://source.com", "aaa", "aaa", false},
    }

    for _, pathInfo := range destTests {
        pr := NewPathResolver(pathInfo.source, pathInfo.dest)
        d, err := pr.resolveDest()
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
                pathInfo.source,
            )
        }

        if d.Name() != pathInfo.resultDest {
            t.Errorf(
                "Resolved destination actual - %s, expected - %s",
                d.Name(),
                pathInfo.resultDest,
            )
        }
    }
}

func initTmpStorage() (string, *os.File) {
    dir, err := ioutil.TempDir("/tmp", "splitter")
    if err != nil {
        log.Fatal(err)
    }
    f, err := os.Create(path.Join(dir, "dest_file.txt"))
    if err != nil {
        log.Fatal(err)
    }

    return dir, f
}
