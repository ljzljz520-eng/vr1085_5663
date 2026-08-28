package audit

import (
	"gatekeeper/model"
	"strings"
)

func EventNames(events []model.AuditEntry) []string {
	out := make([]string, 0, len(events))
	for _, e := range events {
		out = append(out, e.Event)
	}
	return out
}
func JoinValues(events []model.AuditEntry) string {
	v := make([]string, 0, len(events))
	for _, e := range events {
		v = append(v, e.Value)
	}
	return strings.Join(v, "|")
}
func HasSequence(events []model.AuditEntry, n int) bool {
	for _, e := range events {
		if e.Sequence == n {
			return true
		}
	}
	return false
}
