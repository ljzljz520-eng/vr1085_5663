package gatekeeper

import (
	"gatekeeper/model"
	"gatekeeper/store"
	"testing"
	"time"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/gate.db"
	s, e := store.Open(path)
	if e != nil {
		t.Fatal(e)
	}
	s.SaveGateway(model.GatewayKey{GatewayID: "g", PublicKey: "k", Algorithm: "a", Active: true})
	s.Close()
	s, e = store.Open(path)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	if _, e = s.GetGateway("g"); e != nil {
		t.Fatal(e)
	}
	_ = time.Now()
}
