package store

import "database/sql"

func (s *Store) InTransaction(fn func(*sql.Tx) error) error {
	tx, e := s.db.Begin()
	if e != nil {
		return e
	}
	if e = fn(tx); e != nil {
		tx.Rollback()
		return e
	}
	return tx.Commit()
}
func (s *Store) Ping() error   { return s.db.Ping() }
func (s *Store) Vacuum() error { _, e := s.db.Exec("vacuum"); return e }
func (s *Store) DeleteSession(id string) error {
	_, e := s.db.Exec("delete from sessions where id=?", id)
	return e
}
