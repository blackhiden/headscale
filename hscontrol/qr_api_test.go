package hscontrol

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blackhiden/headscale/hscontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistrationQRPayloadHandler_Success(t *testing.T) {
	t.Parallel()

	app := createTestApp(t)
	
	// Create a pending registration
	authID := types.MustAuthID()
	regData := &types.RegistrationData{
		Hostname: "test-node",
	}
	authReq := types.NewRegisterAuthRequest(regData)
	authReq.SetExpiry(time.Now().Add(15 * time.Minute))
	app.state.SetAuthCacheEntry(authID, authReq)

	// Make request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/registration/qr/"+authID.String(), nil)
	rec := httptest.NewRecorder()

	app.HTTPHandler().ServeHTTP(rec, req)

	// Assert response
	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	// Parse response
	var payload struct {
		Type      string    `json:"type"`
		Version   string    `json:"version"`
		AuthID    string    `json:"auth_id"`
		ServerURL string    `json:"server_url"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	err := json.NewDecoder(rec.Body).Decode(&payload)
	require.NoError(t, err)

	assert.Equal(t, "headscale-registration", payload.Type)
	assert.Equal(t, "1", payload.Version)
	assert.Equal(t, authID.String(), payload.AuthID)
	assert.Equal(t, app.cfg.ServerURL, payload.ServerURL)
	assert.WithinDuration(t, authReq.ExpiresAt(), payload.ExpiresAt, time.Second)
}

func TestRegistrationQRPayloadHandler_NotFound(t *testing.T) {
	t.Parallel()

	app := createTestApp(t)
	
	authID := types.MustAuthID()

	// Make request for non-existent auth_id
	req := httptest.NewRequest(http.MethodGet, "/api/v1/registration/qr/"+authID.String(), nil)
	rec := httptest.NewRecorder()

	app.HTTPHandler().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestRegistrationQRPayloadHandler_Expired(t *testing.T) {
	t.Parallel()

	app := createTestApp(t)
	
	// Create an expired pending registration
	authID := types.MustAuthID()
	regData := &types.RegistrationData{
		Hostname: "test-node",
	}
	authReq := types.NewRegisterAuthRequest(regData)
	authReq.SetExpiry(time.Now().Add(-1 * time.Minute)) // Expired
	app.state.SetAuthCacheEntry(authID, authReq)

	// Make request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/registration/qr/"+authID.String(), nil)
	rec := httptest.NewRecorder()

	app.HTTPHandler().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusGone, rec.Code)
}

func TestRegistrationQRPayloadHandler_InvalidAuthID(t *testing.T) {
	t.Parallel()

	app := createTestApp(t)
	
	// Make request with invalid auth_id
	req := httptest.NewRequest(http.MethodGet, "/api/v1/registration/qr/invalid-auth-id", nil)
	rec := httptest.NewRecorder()

	app.HTTPHandler().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegistrationQRPayloadHandler_NotRegistrationRequest(t *testing.T) {
	t.Parallel()

	app := createTestApp(t)
	
	// Create a non-registration auth request (e.g., SSH check)
	authID := types.MustAuthID()
	authReq := types.NewAuthRequest() // Not a registration request
	authReq.SetExpiry(time.Now().Add(15 * time.Minute))
	app.state.SetAuthCacheEntry(authID, authReq)

	// Make request
	req := httptest.NewRequest(http.MethodGet, "/api/v1/registration/qr/"+authID.String(), nil)
	rec := httptest.NewRecorder()

	app.HTTPHandler().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
