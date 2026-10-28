package hscontrol

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	apiv1 "github.com/blackhiden/headscale/hscontrol/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAPIV1ServerInfo proves the read-only GET /api/v1/server-info endpoint
// returns the expected server technical information.
func TestAPIV1ServerInfo(t *testing.T) {
	h := newAPIV1Harness(t)

	res := h.callHuma(http.MethodGet, "/api/v1/server-info", nil)

	assert.Equal(t, http.StatusOK, res.status)

	var body map[string]any
	require.NoError(t, json.Unmarshal(res.body, &body))

	// Version fields
	version, ok := body["version"].(string)
	require.True(t, ok, "expected version field")
	assert.NotEmpty(t, version)

	commit, ok := body["commit"].(string)
	require.True(t, ok, "expected commit field")
	assert.NotEmpty(t, commit)

	buildTime, ok := body["buildTime"].(string)
	require.True(t, ok, "expected buildTime field")
	assert.NotEmpty(t, buildTime)

	// Go info
	goInfo, ok := body["go"].(map[string]any)
	require.True(t, ok, "expected go field")
	assert.NotEmpty(t, goInfo["version"])
	assert.NotEmpty(t, goInfo["os"])
	assert.NotEmpty(t, goInfo["arch"])

	// OS/Arch
	os, ok := body["os"].(string)
	require.True(t, ok, "expected os field")
	assert.NotEmpty(t, os)

	arch, ok := body["arch"].(string)
	require.True(t, ok, "expected arch field")
	assert.NotEmpty(t, arch)

	// CPU
	cpu, ok := body["cpu"].(map[string]any)
	require.True(t, ok, "expected cpu field")
	cpuCount, ok := cpu["count"].(float64)
	require.True(t, ok, "expected cpu.count field")
	assert.Greater(t, cpuCount, float64(0))

	// Memory
	memory, ok := body["memory"].(map[string]any)
	require.True(t, ok, "expected memory field")
	_, ok = memory["totalBytes"].(float64)
	require.True(t, ok, "expected memory.totalBytes field")

	// Storage
	storage, ok := body["storage"].(map[string]any)
	require.True(t, ok, "expected storage field")
	_, ok = storage["totalBytes"].(float64)
	require.True(t, ok, "expected storage.totalBytes field")

	// Deployment
	deployment, ok := body["deployment"].(map[string]any)
	require.True(t, ok, "expected deployment field")
	mode, ok := deployment["mode"].(string)
	require.True(t, ok, "expected deployment.mode field")
	assert.Contains(t, []string{"native", "docker"}, mode)

	// Public IP
	_, ok = body["publicIp"].(string)
	require.True(t, ok, "expected publicIp field")

	// Tailnet
	tailnet, ok := body["tailnet"].(map[string]any)
	require.True(t, ok, "expected tailnet field")
	_, ok = tailnet["ipRangeV4"].(string)
	require.True(t, ok, "expected tailnet.ipRangeV4 field")
	_, ok = tailnet["ipRangeV6"].(string)
	require.True(t, ok, "expected tailnet.ipRangeV6 field")
	_, ok = tailnet["magicDnsDomain"].(string)
	require.True(t, ok, "expected tailnet.magicDnsDomain field")
	_, ok = tailnet["derpRegions"].(float64)
	require.True(t, ok, "expected tailnet.derpRegions field")
}

// TestAPIV1ServerInfoInSpec proves the server-info endpoint appears in the
// emitted OpenAPI document alongside the rest of the v1 API.
func TestAPIV1ServerInfoInSpec(t *testing.T) {
	spec, err := apiv1.Spec()
	require.NoError(t, err)

	assert.Contains(t, string(spec), "/api/v1/server-info")
	assert.Contains(t, string(spec), "ServerInfoResponseBody")
	assert.Contains(t, string(spec), "getServerInfo")
}

// TestAPIV1ServerInfoRequiresAuth proves the server-info endpoint requires
// API key authentication.
func TestAPIV1ServerInfoRequiresAuth(t *testing.T) {
	app := createTestApp(t)
	handler := app.HTTPHandler()

	expiry := time.Now().Add(time.Hour)
	valid, _, err := app.state.CreateAPIKey(&expiry)
	require.NoError(t, err)

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{name: "missing bearer", authHeader: "", wantStatus: http.StatusUnauthorized},
		{name: "no bearer prefix", authHeader: "tskey-invalid", wantStatus: http.StatusUnauthorized},
		{name: "invalid bearer token", authHeader: "Bearer tskey-invalid", wantStatus: http.StatusUnauthorized},
		{name: "valid api key", authHeader: "Bearer " + valid, wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := callHandlerAuth(handler, http.MethodGet, "/api/v1/server-info", tt.authHeader, nil)
			assert.Equalf(t, tt.wantStatus, res.status, "body: %s", res.body)
		})
	}
}
