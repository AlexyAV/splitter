package splitter
//
//import (
//    "fmt"
//    "testing"
//)
//
//func TestNextRange(t *testing.T)  {
//    rb := NewRangeBuilder(100, 10)
//
//    for i := 0; i < 10; i++ {
//        r, err := rb.NextRange()
//
//        if err != nil {
//            t.Error(err)
//        }
//
//        if r.Start == r.End {
//            t.Errorf("DownloadRange start %d = end %d", r.Start, r.End)
//        }
//
//        brh := r.BuildRangeHeader()
//
//        if brh == fmt.Sprintf("bytes=%s-%s", string(r.Start), string(r.End)) {
//            t.Errorf("BuildRangeHeader for start %d and end %d: %s", r.Start, r.End, brh)
//        }
//    }
//
//    _, err := rb.NextRange()
//    if err == nil {
//        t.Error("NextRange() unexpected result end of range EOD")
//    }
//}
