package query

import (
	"fmt"
	"gatekeeper/model"
	"time"
)

func RenderSession(r model.SessionRecord) string {
	return fmt.Sprintf("%s|%s|%s|%s", r.ID, r.GatewayID, r.State, r.ExpiresAt.Format(time.RFC3339))
}
func RenderGateway(k model.GatewayKey) string {
	return fmt.Sprintf("%s|%s|%t", k.GatewayID, k.Fingerprint, k.Active)
}
func FilterEvents(events []model.AuditEntry, name string) []model.AuditEntry {
	out := make([]model.AuditEntry, 0)
	for _, e := range events {
		if name == "" || e.Event == name {
			out = append(out, e)
		}
	}
	return out
}
