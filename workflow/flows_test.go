package workflow

import (
	"gatekeeper/audit"
	"gatekeeper/service"
	"gatekeeper/store"
	"testing"
	"time"
)

func TestWorkflowNegotiation(t *testing.T) {
	s, _ := store.Open(":memory:")
	defer s.Close()
	x := service.New(s, time.Unix(100, 0))
	x.RegisterGateway("g", "key", "a")
	x.RegisterParameters("p", "p", "a", "1")
	c := New(x, audit.New(s))
	if _, e := c.Negotiate("g", "p", "s", "secret"); e != nil {
		t.Fatal(e)
	}
}
func TestWorkflowRecovery(t *testing.T) {
	s, _ := store.Open(":memory:")
	defer s.Close()
	x := service.New(s, time.Unix(100, 0))
	x.RegisterGateway("g", "key", "a")
	x.RegisterParameters("p", "p", "a", "1")
	c := New(x, audit.New(s))
	c.Negotiate("g", "p", "s", "secret")
	if v, e := c.Recover("s"); e != nil || v != "secret" {
		t.Fatal(e, v)
	}
}
