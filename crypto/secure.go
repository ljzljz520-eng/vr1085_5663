package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

func Fingerprint(key string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(key))) }
func WrapSecret(publicKey, secret, parameter string) (string, error) {
	if publicKey == "" || secret == "" || parameter == "" {
		return "", errors.New("missing encryption input")
	}
	sum := sha256.Sum256([]byte(publicKey + "|" + parameter))
	b := make([]byte, len(secret))
	for i := range []byte(secret) {
		b[i] = secret[i] ^ sum[i%len(sum)]
	}
	return base64.RawStdEncoding.EncodeToString(b), nil
}
func UnwrapSecret(publicKey, wrapped, parameter string) (string, error) {
	if publicKey == "" || wrapped == "" || parameter == "" {
		return "", errors.New("missing decryption input")
	}
	b, e := base64.RawStdEncoding.DecodeString(wrapped)
	if e != nil {
		return "", errors.New("bad ciphertext")
	}
	sum := sha256.Sum256([]byte(publicKey + "|" + parameter))
	for i := range b {
		b[i] ^= sum[i%len(sum)]
	}
	return string(b), nil
}
func Digest(secret string) string { return Fingerprint(secret) }
func ValidateCiphertext(wrapped string) error {
	if strings.TrimSpace(wrapped) == "" {
		return errors.New("empty ciphertext")
	}
	if _, e := base64.RawStdEncoding.DecodeString(wrapped); e != nil {
		return errors.New("bad ciphertext")
	}
	return nil
}
