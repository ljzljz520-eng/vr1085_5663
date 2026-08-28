package workflow

import (
	"errors"
	"gatekeeper/crypto"
	"gatekeeper/query"
)

func ValidateWrapped(q *query.Query, id string) error {
	r, e := q.Session(id)
	if e != nil {
		return e
	}
	return crypto.ValidateCiphertext(r.WrappedSecret)
}
func RequireDigest(secret, digest string) error {
	if crypto.Digest(secret) != digest {
		return errors.New("digest mismatch")
	}
	return nil
}
func Readiness(q *query.Query) bool { return q.Health()["store"] == "ready" }
