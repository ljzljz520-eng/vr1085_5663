package service

import (
	"errors"
	"fmt"
	"gatekeeper/crypto"
	"gatekeeper/model"
	"gatekeeper/store"
	"time"
)

type SessionService struct {
	Store *store.Store
	Clock time.Time
}

func New(s *store.Store, now time.Time) *SessionService { return &SessionService{Store: s, Clock: now} }
func (x *SessionService) RegisterGateway(id, key, algorithm string) (model.GatewayKey, error) {
	if id == "" || key == "" {
		return model.GatewayKey{}, errors.New("gateway credentials required")
	}
	k := model.GatewayKey{GatewayID: id, PublicKey: key, Algorithm: algorithm, Fingerprint: crypto.Fingerprint(key), Active: true}
	return k, x.Store.SaveGateway(k)
}
func (x *SessionService) RegisterParameters(id, name, algorithm, version string) (model.ParameterSet, error) {
	if id == "" || algorithm == "" {
		return model.ParameterSet{}, errors.New("parameter set required")
	}
	p := model.ParameterSet{ID: id, Name: name, Algorithm: algorithm, Version: version, Enabled: true, CreatedAt: x.Clock}
	return p, x.Store.SaveParameters(p)
}
func (x *SessionService) OpenSession(id, gateway, param, secret string, ttl time.Duration) (model.SessionRecord, error) {
	k, e := x.Store.GetGateway(gateway)
	if e != nil {
		return model.SessionRecord{}, fmt.Errorf("gateway: %w", e)
	}
	p, e := x.Store.GetParameters(param)
	if e != nil {
		return model.SessionRecord{}, fmt.Errorf("parameters: %w", e)
	}
	if !p.Supports(k.Algorithm) {
		return model.SessionRecord{}, errors.New("parameter algorithm mismatch")
	}
	w, e := crypto.WrapSecret(k.PublicKey, secret, p.ID)
	if e != nil {
		return model.SessionRecord{}, e
	}
	r := model.SessionRecord{ID: id, GatewayID: gateway, ParameterSetID: param, WrappedSecret: w, State: "open", CreatedAt: x.Clock, ExpiresAt: x.Clock.Add(ttl)}
	return r, x.Store.SaveSession(r)
}
func (x *SessionService) UnwrapSession(id string) (model.SessionSecret, error) {
	r, e := x.Store.GetSession(id)
	if e != nil {
		return model.SessionSecret{}, e
	}
	if r.IsExpired(x.Clock) {
		return model.SessionSecret{}, errors.New("session expired")
	}
	k, e := x.Store.GetGateway(r.GatewayID)
	if e != nil {
		return model.SessionSecret{}, e
	}
	p, e := x.Store.GetParameters(r.ParameterSetID)
	if e != nil {
		return model.SessionSecret{}, e
	}
	plain, e := crypto.UnwrapSecret(k.PublicKey, r.WrappedSecret, p.ID)
	if e != nil {
		return model.SessionSecret{}, e
	}
	return model.SessionSecret{SessionID: id, Plaintext: plain, Digest: crypto.Digest(plain), ValidUntil: r.ExpiresAt}, nil
}
func (x *SessionService) CloseSession(id string) error {
	r, e := x.Store.GetSession(id)
	if e != nil {
		return e
	}
	r.State = "closed"
	return x.Store.SaveSession(r)
}
func (x *SessionService) Describe(id string) (string, error) {
	r, e := x.Store.GetSession(id)
	if e != nil {
		return "", e
	}
	return fmt.Sprintf("session=%s gateway=%s state=%s", r.ID, r.GatewayID, r.State), nil
}
