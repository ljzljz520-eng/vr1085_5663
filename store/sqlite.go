package store

import (
	"database/sql"
	"errors"
	"fmt"
	"gatekeeper/model"
	_ "modernc.org/sqlite"
	"time"
)

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	db, e := sql.Open("sqlite", path)
	if e != nil {
		return nil, e
	}
	s := &Store{db: db}
	if e = s.migrate(); e != nil {
		db.Close()
		return nil, e
	}
	return s, nil
}
func (s *Store) migrate() error {
	_, e := s.db.Exec(`create table if not exists gateways(id text primary key, public_key text, algorithm text, fingerprint text, active integer); create table if not exists parameters(id text primary key,name text,algorithm text,version text,enabled integer,created_at text); create table if not exists sessions(id text primary key,gateway_id text,parameter_id text,wrapped_secret text,state text,created_at text,expires_at text); create table if not exists audits(id text primary key,session_id text,gateway_id text,event text,value text,at text,sequence integer)`)
	return e
}
func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
func (s *Store) SaveGateway(k model.GatewayKey) error {
	_, e := s.db.Exec(`insert or replace into gateways values(?,?,?,?,?)`, k.GatewayID, k.PublicKey, k.Algorithm, k.Fingerprint, k.Active)
	return e
}
func (s *Store) GetGateway(id string) (model.GatewayKey, error) {
	var k model.GatewayKey
	var a int
	e := s.db.QueryRow(`select id,public_key,algorithm,fingerprint,active from gateways where id=?`, id).Scan(&k.GatewayID, &k.PublicKey, &k.Algorithm, &k.Fingerprint, &a)
	k.Active = a == 1
	return k, e
}
func (s *Store) SaveParameters(p model.ParameterSet) error {
	_, e := s.db.Exec(`insert or replace into parameters values(?,?,?,?,?,?)`, p.ID, p.Name, p.Algorithm, p.Version, p.Enabled, p.CreatedAt.Format(time.RFC3339Nano))
	return e
}
func (s *Store) GetParameters(id string) (model.ParameterSet, error) {
	var p model.ParameterSet
	var a int
	var t string
	e := s.db.QueryRow(`select id,name,algorithm,version,enabled,created_at from parameters where id=?`, id).Scan(&p.ID, &p.Name, &p.Algorithm, &p.Version, &a, &t)
	p.Enabled = a == 1
	p.CreatedAt, _ = time.Parse(time.RFC3339Nano, t)
	return p, e
}
func (s *Store) SaveSession(r model.SessionRecord) error {
	if r.ID == "" {
		return errors.New("session id required")
	}
	_, e := s.db.Exec(`insert or replace into sessions values(?,?,?,?,?,?,?)`, r.ID, r.GatewayID, r.ParameterSetID, r.WrappedSecret, r.State, r.CreatedAt.Format(time.RFC3339Nano), r.ExpiresAt.Format(time.RFC3339Nano))
	return e
}
func (s *Store) GetSession(id string) (model.SessionRecord, error) {
	var r model.SessionRecord
	var c, x string
	e := s.db.QueryRow(`select id,gateway_id,parameter_id,wrapped_secret,state,created_at,expires_at from sessions where id=?`, id).Scan(&r.ID, &r.GatewayID, &r.ParameterSetID, &r.WrappedSecret, &r.State, &c, &x)
	r.CreatedAt, _ = time.Parse(time.RFC3339Nano, c)
	r.ExpiresAt, _ = time.Parse(time.RFC3339Nano, x)
	return r, e
}
func (s *Store) SaveAudit(a model.AuditEntry) error {
	_, e := s.db.Exec(`insert or replace into audits values(?,?,?,?,?,?,?)`, a.ID, a.SessionID, a.GatewayID, a.Event, a.Value, a.At.Format(time.RFC3339Nano), a.Sequence)
	return e
}
func (s *Store) ListAudits(session string) ([]model.AuditEntry, error) {
	rows, e := s.db.Query(`select id,session_id,gateway_id,event,value,at,sequence from audits where session_id=? order by sequence`, session)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []model.AuditEntry
	for rows.Next() {
		var a model.AuditEntry
		var t string
		if e = rows.Scan(&a.ID, &a.SessionID, &a.GatewayID, &a.Event, &a.Value, &t, &a.Sequence); e != nil {
			return nil, e
		}
		a.At, _ = time.Parse(time.RFC3339Nano, t)
		out = append(out, a)
	}
	return out, rows.Err()
}
func (s *Store) Count(table string) (int, error) {
	if table != "gateways" && table != "parameters" && table != "sessions" && table != "audits" {
		return 0, fmt.Errorf("invalid table")
	}
	var n int
	e := s.db.QueryRow("select count(*) from " + table).Scan(&n)
	return n, e
}
