package store

import (
	"gatekeeper/model"
	"testing"
	"time"
)

func TestStoreRoundTrip(t *testing.T) {
	s, e := Open(":memory:")
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	p := model.ParameterSet{ID: "p", Algorithm: "xor", Enabled: true, CreatedAt: time.Now()}
	if e = s.SaveParameters(p); e != nil {
		t.Fatal(e)
	}
	if _, e = s.GetParameters("p"); e != nil {
		t.Fatal(e)
	}
}
