package workflow

import (
	"errors"
	"fmt"
	"gatekeeper/audit"
	"gatekeeper/service"
	"time"
)

type Coordinator struct {
	Sessions *service.SessionService
	Audit    *audit.Recorder
}

func New(s *service.SessionService, a *audit.Recorder) *Coordinator {
	return &Coordinator{Sessions: s, Audit: a}
}
func (c *Coordinator) Negotiate(gateway, param, session, secret string) (string, error) {
	if gateway == "" || param == "" {
		return "", errors.New("missing negotiation identity")
	}
	r, e := c.Sessions.OpenSession(session, gateway, param, secret, time.Hour)
	if e != nil {
		return "", e
	}
	if e = c.Audit.RecordBatch(session, gateway, []string{"received-key", "wrapped-secret", "opened-session"}, c.Sessions.Clock); e != nil {
		return "", e
	}
	return fmt.Sprintf("%s:%s", r.ID, r.WrappedSecret), nil
}
func (c *Coordinator) Recover(session string) (string, error) {
	s, e := c.Sessions.UnwrapSession(session)
	if e != nil {
		return "", e
	}
	return s.Plaintext, nil
}
func (c *Coordinator) Retire(session string) error { return c.Sessions.CloseSession(session) }
func (c *Coordinator) Verify(session string) (bool, error) {
	s, e := c.Sessions.UnwrapSession(session)
	return e == nil && s.Plaintext != "", e
}
