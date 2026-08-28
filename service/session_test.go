package service

import (
	"gatekeeper/store"
	"testing"
	"time"
)

func TestSessionService(t *testing.T) {
	s, _ := store.Open(":memory:")
	defer s.Close()
	x := New(s, time.Unix(100, 0))
	if _, e := x.RegisterGateway("g", "key", "a"); e != nil {
		t.Fatal(e)
	}
	if _, e := x.RegisterParameters("p", "p", "a", "1"); e != nil {
		t.Fatal(e)
	}
	if _, e := x.OpenSession("s", "g", "p", "daily", time.Hour); e != nil {
		t.Fatal(e)
	}
}
