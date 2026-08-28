package model

import "time"

type SessionRecord struct {
	ID, GatewayID, ParameterSetID, WrappedSecret, State string
	CreatedAt                                           time.Time
	ExpiresAt                                           time.Time
}
type GatewayKey struct {
	GatewayID, PublicKey, Algorithm, Fingerprint string
	Active                                       bool
}
type AuditEntry struct {
	ID, SessionID, GatewayID, Event, Value string
	At                                     time.Time
	Sequence                               int
}
type ParameterSet struct {
	ID, Name, Algorithm, Version string
	Enabled                      bool
	CreatedAt                    time.Time
}
type SessionSecret struct {
	SessionID, Plaintext, Digest string
	ValidUntil                   time.Time
}

func (s SessionRecord) IsExpired(now time.Time) bool  { return !now.Before(s.ExpiresAt) }
func (s SessionRecord) IsOpen() bool                  { return s.State == "open" }
func (k GatewayKey) IsUsable() bool                   { return k.Active && k.PublicKey != "" && k.Algorithm != "" }
func (p ParameterSet) Supports(algorithm string) bool { return p.Enabled && p.Algorithm == algorithm }
