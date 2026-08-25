package query

import (
	"gatekeeper/model"
	"gatekeeper/store"
	"sort"
)

type Query struct{ Store *store.Store }

func New(s *store.Store) *Query                                 { return &Query{Store: s} }
func (q *Query) Session(id string) (model.SessionRecord, error) { return q.Store.GetSession(id) }
func (q *Query) Gateway(id string) (model.GatewayKey, error)    { return q.Store.GetGateway(id) }
func (q *Query) Audits(id string) ([]model.AuditEntry, error) {
	a, e := q.Store.ListAudits(id)
	sort.Slice(a, func(i, j int) bool { return a[i].Sequence < a[j].Sequence })
	return a, e
}
func (q *Query) Active(id string) (bool, error) {
	r, e := q.Store.GetSession(id)
	return e == nil && r.IsOpen(), e
}
func (q *Query) Health() map[string]string {
	return map[string]string{"store": "ready", "persistence": "sqlite"}
}
