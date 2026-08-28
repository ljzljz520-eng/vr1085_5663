package gatekeeper

import (
	"gatekeeper/audit"
	"gatekeeper/service"
	"gatekeeper/store"
	"gatekeeper/workflow"
	"testing"
	"time"
)

func TestBusinessChain37(t *testing.T) {
	s, _ := store.Open(":memory:")
	defer s.Close()
	x := service.New(s, time.Unix(100, 0))
	x.RegisterGateway("g", "key", "a")
	x.RegisterParameters("p", "p", "a", "1")
	c := workflow.New(x, audit.New(s))
	if _, e := c.Negotiate("g", "p", "s", "secret"); e != nil {
		t.Fatal(e)
	}
	a, _ := s.ListAudits("s")
	for i, v := range a {
		if v.Value != []string{"received-key", "wrapped-secret", "opened-session"}[i] {
			t.Fatalf("audit %d=%s", i, v.Value)
		}
	}
}
func TestWorkflowRetirement(t *testing.T) {
	s, _ := store.Open(":memory:")
	defer s.Close()
	x := service.New(s, time.Unix(100, 0))
	x.RegisterGateway("g", "key", "a")
	x.RegisterParameters("p", "p", "a", "1")
	x.OpenSession("s", "g", "p", "secret", time.Hour)
	if e := x.CloseSession("s"); e != nil {
		t.Fatal(e)
	}
}
