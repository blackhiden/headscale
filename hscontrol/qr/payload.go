package qr

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	errEmptyAuthID         = errors.New("auth_id cannot be empty")
	errEmptyServerURL      = errors.New("server_url cannot be empty")
	errZeroExpiry          = errors.New("expires_at cannot be zero")
	errExpiryNotInFuture   = errors.New("expires_at must be in the future")
	errInvalidPayloadType  = errors.New("invalid payload type")
	errUnsupportedVersion  = errors.New("unsupported payload version")
	errRegistrationExpired = errors.New("payload has expired")
)

// RegistrationPayload represents the JSON payload encoded in a QR code
// for headscale node registration.
type RegistrationPayload struct {
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	AuthID    string    `json:"auth_id"`
	ServerURL string    `json:"server_url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewRegistrationPayload creates and validates a new registration payload.
// It enforces non-empty fields and ensures the expiry is in the future.
func NewRegistrationPayload(
	authID, serverURL string,
	expiresAt time.Time,
) (*RegistrationPayload, error) {
	if authID == "" {
		return nil, errEmptyAuthID
	}

	if serverURL == "" {
		return nil, errEmptyServerURL
	}

	if expiresAt.IsZero() {
		return nil, errZeroExpiry
	}

	if expiresAt.Before(time.Now()) {
		return nil, errExpiryNotInFuture
	}

	return &RegistrationPayload{
		Type:      "headscale-registration",
		Version:   "1",
		AuthID:    authID,
		ServerURL: serverURL,
		ExpiresAt: expiresAt,
	}, nil
}

// ToJSON serializes the payload to JSON bytes.
func (p *RegistrationPayload) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

// ParseRegistrationPayload deserializes and validates a registration payload
// from JSON bytes.
func ParseRegistrationPayload(data []byte) (*RegistrationPayload, error) {
	var payload RegistrationPayload

	err := json.Unmarshal(data, &payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	// Validate the parsed payload
	if payload.Type != "headscale-registration" {
		return nil, fmt.Errorf("%w: %s", errInvalidPayloadType, payload.Type)
	}

	if payload.Version != "1" {
		return nil, fmt.Errorf("%w: %s", errUnsupportedVersion, payload.Version)
	}

	if payload.AuthID == "" {
		return nil, errEmptyAuthID
	}

	if payload.ServerURL == "" {
		return nil, errEmptyServerURL
	}

	if payload.ExpiresAt.IsZero() {
		return nil, errZeroExpiry
	}

	if payload.ExpiresAt.Before(time.Now()) {
		return nil, errRegistrationExpired
	}

	return &payload, nil
}
