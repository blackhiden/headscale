package qr

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegistrationPayload_Valid(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(24 * time.Hour)

	payload, err := NewRegistrationPayload(authID, serverURL, expiresAt)

	require.NoError(t, err)
	assert.Equal(t, "headscale-registration", payload.Type)
	assert.Equal(t, "1", payload.Version)
	assert.Equal(t, authID, payload.AuthID)
	assert.Equal(t, serverURL, payload.ServerURL)
	assert.Equal(t, expiresAt, payload.ExpiresAt)
}

func TestNewRegistrationPayload_EmptyAuthID(t *testing.T) {
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(24 * time.Hour)

	payload, err := NewRegistrationPayload("", serverURL, expiresAt)

	require.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "auth_id cannot be empty")
}

func TestNewRegistrationPayload_EmptyServerURL(t *testing.T) {
	authID := "hskey-abc123"
	expiresAt := time.Now().Add(24 * time.Hour)

	payload, err := NewRegistrationPayload(authID, "", expiresAt)

	require.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "server_url cannot be empty")
}

func TestNewRegistrationPayload_ZeroExpiry(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"

	payload, err := NewRegistrationPayload(authID, serverURL, time.Time{})

	require.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "expires_at cannot be zero")
}

func TestNewRegistrationPayload_PastExpiry(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(-24 * time.Hour)

	payload, err := NewRegistrationPayload(authID, serverURL, expiresAt)

	require.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "expires_at must be in the future")
}

func TestPayload_ToJSON(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	payload, err := NewRegistrationPayload(authID, serverURL, expiresAt)
	require.NoError(t, err)

	jsonData, err := payload.ToJSON()
	require.NoError(t, err)

	// Verify JSON structure
	var decoded map[string]any

	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "headscale-registration", decoded["type"])
	assert.Equal(t, "1", decoded["version"])
	assert.Equal(t, authID, decoded["auth_id"])
	assert.Equal(t, serverURL, decoded["server_url"])
	assert.Equal(t, "2027-01-01T00:00:00Z", decoded["expires_at"])
}

func TestParseRegistrationPayload_Valid(t *testing.T) {
	futureTime := time.Now().Add(24 * time.Hour)
	jsonData := []byte(fmt.Sprintf(`{
		"type": "headscale-registration",
		"version": "1",
		"auth_id": "hskey-abc123",
		"server_url": "https://headscale.example.com",
		"expires_at": "%s"
	}`, futureTime.Format(time.RFC3339)))

	payload, err := ParseRegistrationPayload(jsonData)

	require.NoError(t, err)
	assert.Equal(t, "headscale-registration", payload.Type)
	assert.Equal(t, "1", payload.Version)
	assert.Equal(t, "hskey-abc123", payload.AuthID)
	assert.Equal(t, "https://headscale.example.com", payload.ServerURL)
	assert.WithinDuration(t, futureTime, payload.ExpiresAt, time.Second)
}

func TestParseRegistrationPayload_InvalidJSON(t *testing.T) {
	jsonData := []byte(`{invalid json}`)

	payload, err := ParseRegistrationPayload(jsonData)

	require.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "failed to parse payload")
}

func TestParseRegistrationPayload_InvalidType(t *testing.T) {
	jsonData := []byte(`{
		"type": "wrong-type",
		"version": "1",
		"auth_id": "hskey-abc123",
		"server_url": "https://headscale.example.com",
		"expires_at": "2026-10-08T07:00:00Z"
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	require.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "invalid payload type")
}

func TestParseRegistrationPayload_InvalidVersion(t *testing.T) {
	jsonData := []byte(`{
		"type": "headscale-registration",
		"version": "2",
		"auth_id": "hskey-abc123",
		"server_url": "https://headscale.example.com",
		"expires_at": "2026-10-08T07:00:00Z"
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	require.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "unsupported payload version")
}

func TestParseRegistrationPayload_EmptyAuthID(t *testing.T) {
	jsonData := []byte(`{
		"type": "headscale-registration",
		"version": "1",
		"auth_id": "",
		"server_url": "https://headscale.example.com",
		"expires_at": "2026-10-08T07:00:00Z"
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	require.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "auth_id cannot be empty")
}

func TestParseRegistrationPayload_EmptyServerURL(t *testing.T) {
	jsonData := []byte(`{
		"type": "headscale-registration",
		"version": "1",
		"auth_id": "hskey-abc123",
		"server_url": "",
		"expires_at": "2026-10-08T07:00:00Z"
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	require.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "server_url cannot be empty")
}

func TestParseRegistrationPayload_ZeroExpiry(t *testing.T) {
	jsonData := []byte(`{
		"type": "headscale-registration",
		"version": "1",
		"auth_id": "hskey-abc123",
		"server_url": "https://headscale.example.com",
		"expires_at": "0001-01-01T00:00:00Z"
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	require.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "expires_at cannot be zero")
}

func TestParseRegistrationPayload_ExpiredAtParseTime(t *testing.T) {
	// Simulate a QR code that was generated in the past and is now expired
	jsonData := []byte(`{
		"type": "headscale-registration",
		"version": "1",
		"auth_id": "hskey-abc123",
		"server_url": "https://headscale.example.com",
		"expires_at": "2020-01-01T00:00:00Z"
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	require.Error(t, err)
	assert.Nil(t, payload)
	assert.Contains(t, err.Error(), "payload has expired")
}

func TestPayload_RoundTrip(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	// Create original payload
	original, err := NewRegistrationPayload(authID, serverURL, expiresAt)
	require.NoError(t, err)

	// Serialize to JSON
	jsonData, err := original.ToJSON()
	require.NoError(t, err)

	// Parse back
	parsed, err := ParseRegistrationPayload(jsonData)
	require.NoError(t, err)

	// Verify all fields match
	assert.Equal(t, original.Type, parsed.Type)
	assert.Equal(t, original.Version, parsed.Version)
	assert.Equal(t, original.AuthID, parsed.AuthID)
	assert.Equal(t, original.ServerURL, parsed.ServerURL)
	assert.Equal(t, original.ExpiresAt, parsed.ExpiresAt)
}

func TestParseRegistrationPayload_MalformedJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData []byte
		errMsg   string
	}{
		{
			name:     "invalid JSON",
			jsonData: []byte(`{invalid json`),
			errMsg:   "failed to parse payload",
		},
		{
			name:     "empty JSON",
			jsonData: []byte(``),
			errMsg:   "failed to parse payload",
		},
		{
			name:     "null JSON",
			jsonData: []byte(`null`),
			errMsg:   "invalid payload type",
		},
		{
			name:     "array instead of object",
			jsonData: []byte(`[]`),
			errMsg:   "failed to parse payload",
		},
		{
			name:     "missing type field",
			jsonData: []byte(`{"version":"1","auth_id":"hskey-abc","server_url":"https://example.com","expires_at":"2027-01-01T00:00:00Z"}`),
			errMsg:   "invalid payload type",
		},
		{
			name:     "missing version field",
			jsonData: []byte(`{"type":"headscale-registration","auth_id":"hskey-abc","server_url":"https://example.com","expires_at":"2027-01-01T00:00:00Z"}`),
			errMsg:   "unsupported payload version",
		},
		{
			name:     "invalid time format",
			jsonData: []byte(`{"type":"headscale-registration","version":"1","auth_id":"hskey-abc","server_url":"https://example.com","expires_at":"not-a-time"}`),
			errMsg:   "failed to parse payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := ParseRegistrationPayload(tt.jsonData)

			require.Error(t, err)
			assert.Nil(t, payload)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestNewRegistrationPayload_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name      string
		authID    string
		serverURL string
		shouldErr bool
	}{
		{
			name: "unicode in auth ID",
			//nolint:gosmopolitan // Test QR payload round-tripping Unicode data.
			authID:    "hskey-测试-😀",
			serverURL: "https://headscale.example.com",
			shouldErr: false,
		},
		{
			name:      "special chars in URL",
			serverURL: "https://headscale.example.com/path?query=value&foo=bar#fragment",
			authID:    "hskey-abc123",
			shouldErr: false,
		},
		{
			name:      "URL with port",
			serverURL: "https://headscale.example.com:8443",
			authID:    "hskey-abc123",
			shouldErr: false,
		},
		{
			name:      "very long auth ID",
			authID:    strings.Repeat("a", 1000),
			serverURL: "https://headscale.example.com",
			shouldErr: false,
		},
		{
			name:      "very long URL",
			serverURL: "https://headscale.example.com/" + strings.Repeat("path/", 100),
			authID:    "hskey-abc123",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expiresAt := time.Now().Add(24 * time.Hour)
			payload, err := NewRegistrationPayload(tt.authID, tt.serverURL, expiresAt)

			if tt.shouldErr {
				require.Error(t, err)
				assert.Nil(t, payload)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.authID, payload.AuthID)
				assert.Equal(t, tt.serverURL, payload.ServerURL)

				// Ensure it can round-trip
				jsonData, err := payload.ToJSON()
				require.NoError(t, err)

				parsed, err := ParseRegistrationPayload(jsonData)
				require.NoError(t, err)
				assert.Equal(t, payload.AuthID, parsed.AuthID)
				assert.Equal(t, payload.ServerURL, parsed.ServerURL)
			}
		})
	}
}

func TestParseRegistrationPayload_ExtraFields(t *testing.T) {
	// JSON with extra unknown fields should still parse successfully
	jsonData := []byte(`{
		"type": "headscale-registration",
		"version": "1",
		"auth_id": "hskey-abc123",
		"server_url": "https://headscale.example.com",
		"expires_at": "2027-01-01T00:00:00Z",
		"extra_field": "should be ignored",
		"another_field": 12345
	}`)

	payload, err := ParseRegistrationPayload(jsonData)

	require.NoError(t, err)
	assert.Equal(t, "headscale-registration", payload.Type)
	assert.Equal(t, "1", payload.Version)
	assert.Equal(t, "hskey-abc123", payload.AuthID)
	assert.Equal(t, "https://headscale.example.com", payload.ServerURL)
}
