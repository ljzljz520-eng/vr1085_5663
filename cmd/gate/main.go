package main

import (
	"fmt"
	"gatekeeper/audit"
	"gatekeeper/service"
	"gatekeeper/store"
	"gatekeeper/workflow"
	"os"
	"time"
)

func main() {
	path := "gate.db"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	s, e := store.Open(path)
	if e != nil {
		panic(e)
	}
	defer s.Close()
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	svc := service.New(s, now)
	svc.RegisterGateway("gate-1", "public-key-1", "xor-sha256")
	svc.RegisterParameters("params-1", "daily", "xor-sha256", "v1")
	c := workflow.New(svc, audit.New(s))
	result, e := c.Negotiate("gate-1", "params-1", "session-1", "daily-secret")
	if e != nil {
		panic(e)
	}
	sum, _ := audit.New(s).Summary("session-1")
	fmt.Printf("negotiation=%s audit=%s\n", result, sum)
}
