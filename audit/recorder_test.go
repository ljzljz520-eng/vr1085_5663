package audit

import (
	"gatekeeper/store"
	"testing"
	"time"
)

func TestRecorderSummary(t *testing.T) {
	s, _ := store.Open(":memory:")
	defer s.Close()
	r := New(s)
	if e := r.Record("s", "g", "step", "v", 0, time.Unix(1, 0)); e != nil {
		t.Fatal(e)
	}
	v, _ := r.Summary("s")
	if v == "" {
		t.Fatal("empty")
	}
}
