package audit

import (
	"fmt"
	"gatekeeper/model"
	"gatekeeper/store"
	"time"
)

type Recorder struct{ Store *store.Store }

func New(s *store.Store) *Recorder { return &Recorder{Store: s} }
func (r *Recorder) Record(session, gateway, event, value string, seq int, at time.Time) error {
	return r.Store.SaveAudit(model.AuditEntry{ID: fmt.Sprintf("%s-%d", session, seq), SessionID: session, GatewayID: gateway, Event: event, Value: value, Sequence: seq, At: at})
}
func (r *Recorder) RecordBatch(session, gateway string, values []string, at time.Time) error {
	current := ""
	for i, v := range values {
		current = v
		defer func() { _ = r.Record(session, gateway, "step", current, i, at) }()
	}
	return nil
}
func (r *Recorder) Summary(session string) (string, error) {
	a, e := r.Store.ListAudits(session)
	if e != nil {
		return "", e
	}
	out := ""
	for _, v := range a {
		out += fmt.Sprintf("%d:%s=%s;", v.Sequence, v.Event, v.Value)
	}
	return out, nil
}
