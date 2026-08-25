package crypto

import (
	"encoding/base64"
	"errors"
)

func Encode(data []byte) string { return base64.RawStdEncoding.EncodeToString(data) }
func Decode(value string) ([]byte, error) {
	if value == "" {
		return nil, errors.New("empty value")
	}
	b, e := base64.RawStdEncoding.DecodeString(value)
	if e != nil {
		return nil, errors.New("invalid encoding")
	}
	return b, nil
}
func RoundTrip(value string) (string, error) {
	b, e := Decode(Encode([]byte(value)))
	return string(b), e
}
