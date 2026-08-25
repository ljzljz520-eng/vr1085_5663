package query

import (
	"gatekeeper/store"
	"testing"
)

func TestHealth(t *testing.T) {
	s, _ := store.Open(":memory:")
	defer s.Close()
	if New(s).Health()["store"] != "ready" {
		t.Fatal("health")
	}
}
