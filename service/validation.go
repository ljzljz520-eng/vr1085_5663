package service

import (
	"errors"
	"gatekeeper/model"
	"time"
)

func ValidateGateway(k model.GatewayKey) error {
	if !k.IsUsable() {
		return errors.New("gateway key unusable")
	}
	return nil
}
func ValidateTTL(ttl time.Duration) error {
	if ttl <= 0 {
		return errors.New("ttl must be positive")
	}
	if ttl > 24*time.Hour {
		return errors.New("ttl exceeds day")
	}
	return nil
}
func ValidateSecret(secret string) error {
	if len(secret) < 4 {
		return errors.New("secret too short")
	}
	return nil
}
