package store

import (
	"database/sql"
	"gatekeeper/model"
	"time"
)

func ScanSession(row *sql.Row) (model.SessionRecord, error) {
	var r model.SessionRecord
	var a, b string
	e := row.Scan(&r.ID, &r.GatewayID, &r.ParameterSetID, &r.WrappedSecret, &r.State, &a, &b)
	r.CreatedAt, _ = time.Parse(time.RFC3339Nano, a)
	r.ExpiresAt, _ = time.Parse(time.RFC3339Nano, b)
	return r, e
}
func (s *Store) SessionsByGateway(gateway string) ([]model.SessionRecord, error) {
	rows, e := s.db.Query("select id,gateway_id,parameter_id,wrapped_secret,state,created_at,expires_at from sessions where gateway_id=?", gateway)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []model.SessionRecord
	for rows.Next() {
		var r model.SessionRecord
		var a, b string
		if e = rows.Scan(&r.ID, &r.GatewayID, &r.ParameterSetID, &r.WrappedSecret, &r.State, &a, &b); e != nil {
			return nil, e
		}
		r.CreatedAt, _ = time.Parse(time.RFC3339Nano, a)
		r.ExpiresAt, _ = time.Parse(time.RFC3339Nano, b)
		out = append(out, r)
	}
	return out, rows.Err()
}
