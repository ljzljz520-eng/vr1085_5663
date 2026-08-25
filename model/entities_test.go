package model

import (
	"testing"
	"time"
)

func TestEntityRules(t *testing.T) {
	r := SessionRecord{State: "open", ExpiresAt: time.Now().Add(time.Hour)}
	if !r.IsOpen() || r.IsExpired(time.Now()) {
		t.Fatal("entity rule")
	}
	if !DefaultPolicy().Allows("xor-sha256") {
		t.Fatal("policy")
	}
}
