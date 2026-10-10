package hscontrol

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apiv1 "github.com/arsydoni4326-alt/headscale/hscontrol/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAPIV1UpdateCheck proves the update-check endpoint is served by the v1
// API and returns the expected response shape.
func TestAPIV1UpdateCheck(t *testing.T) {
	h := newAPIV1Harness(t)

	t.Run("returns version info", func(t *testing.T) {
		res := h.callHuma(http.MethodGet, "/api/v1/update-check", nil)

		assert.Equal(t, http.StatusOK, res.status)

		var body map[string]any
		require.NoError(t, json.Unmarshal(res.body, &body))

		current, ok := body["current"].(map[string]any)
		require.True(t, ok, "expected current object in response")
		assert.NotEmpty(t, current["version"])
		assert.NotEmpty(t, current["commit"])
		assert.NotEmpty(t, current["buildTime"])
	})

	t.Run("accepts check=true", func(t *testing.T) {
		res := h.callHuma(http.MethodGet, "/api/v1/update-check?check=true", nil)

		assert.Equal(t, http.StatusOK, res.status)
	})
}

// TestAPIV1UpdateCheckIsPublic proves the update-check endpoint is reachable
// without an API key through the full router, matching its no-auth behavior.
func TestAPIV1UpdateCheckIsPublic(t *testing.T) {
	app := createTestApp(t)
	handler := app.HTTPHandler()

	req := httptest.NewRequestWithContext(
		context.Background(), http.MethodGet, "/api/v1/update-check", nil,
	)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}

// TestAPIV1UpdateCheckInSpec proves the update-check endpoint appears in the
// emitted OpenAPI document alongside the rest of the v1 API.
func TestAPIV1UpdateCheckInSpec(t *testing.T) {
	spec, err := apiv1.Spec()
	require.NoError(t, err)

	assert.Contains(t, string(spec), "/api/v1/update-check")
	assert.Contains(t, string(spec), "UpdateCheckResult")
}
