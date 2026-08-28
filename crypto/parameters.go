package crypto

import "errors"

type Parameter struct {
	ID        string
	Algorithm string
	KeySize   int
	Version   string
}

func (p Parameter) Valid() bool {
	return p.ID != "" && p.Algorithm != "" && p.KeySize >= 128 && p.Version != ""
}
func ParseParameter(id, algorithm, version string) (Parameter, error) {
	p := Parameter{ID: id, Algorithm: algorithm, Version: version, KeySize: 256}
	if !p.Valid() {
		return p, errors.New("invalid parameter")
	}
	return p, nil
}
func Compatible(a, b Parameter) bool { return a.Valid() && b.Valid() && a.Algorithm == b.Algorithm }
