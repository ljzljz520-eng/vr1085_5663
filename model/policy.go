package model

import "errors"

type SessionPolicy struct {
	MinTTLSeconds        int
	AllowedAlgorithms    []string
	RequireActiveGateway bool
}

func (p SessionPolicy) Validate() error {
	if p.MinTTLSeconds <= 0 {
		return errors.New("ttl policy invalid")
	}
	if len(p.AllowedAlgorithms) == 0 {
		return errors.New("algorithm policy empty")
	}
	return nil
}
func (p SessionPolicy) Allows(algorithm string) bool {
	for _, a := range p.AllowedAlgorithms {
		if a == algorithm {
			return true
		}
	}
	return false
}
func DefaultPolicy() SessionPolicy {
	return SessionPolicy{MinTTLSeconds: 60, AllowedAlgorithms: []string{"xor-sha256"}, RequireActiveGateway: true}
}
