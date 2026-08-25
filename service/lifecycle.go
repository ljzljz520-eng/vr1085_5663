package service

import (
	"errors"
	"gatekeeper/model"
	"gatekeeper/store"
)

type Lifecycle struct{ Store *store.Store }

func NewLifecycle(s *store.Store) *Lifecycle { return &Lifecycle{Store: s} }
func (l *Lifecycle) Transition(id, state string) (model.SessionRecord, error) {
	r, e := l.Store.GetSession(id)
	if e != nil {
		return r, e
	}
	if state != "open" && state != "closed" {
		return r, errors.New("unknown state")
	}
	if r.State == "closed" && state == "open" {
		return r, errors.New("closed session immutable")
	}
	r.State = state
	return r, l.Store.SaveSession(r)
}
func (l *Lifecycle) IsClosed(id string) (bool, error) {
	r, e := l.Store.GetSession(id)
	return e == nil && r.State == "closed", e
}
func (l *Lifecycle) RequireOpen(id string) error {
	r, e := l.Store.GetSession(id)
	if e != nil {
		return e
	}
	if !r.IsOpen() {
		return errors.New("session not open")
	}
	return nil
}
