package model

import "time"

type SessionEvent struct {
	Name      string
	SessionID string
	GatewayID string
	Payload   string
	At        time.Time
}

func (e SessionEvent) Valid() bool { return e.Name != "" && e.SessionID != "" && e.GatewayID != "" }
func (e SessionEvent) IsSecurityEvent() bool {
	return e.Name == "wrapped-secret" || e.Name == "opened-session"
}
func NewEvent(name, session, gateway, payload string, at time.Time) SessionEvent {
	return SessionEvent{Name: name, SessionID: session, GatewayID: gateway, Payload: payload, At: at}
}
