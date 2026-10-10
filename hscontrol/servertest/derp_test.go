package servertest

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tailscale.com/tailcfg"
)

// TestDERPEndpoint verifies the read-only GET /api/v1/derp endpoint returns
// the configured DERP map with its regions and nodes.
func TestDERPEndpoint(t *testing.T) {
	srv := NewServer(t, WithRealListener())

	user := srv.CreateUser(t, "derp-test")
	apiKey := srv.CreateAPIKey(t, user)

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL+"/api/v1/derp", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var derp struct {
		Configured   bool `json:"configured"`
		TotalRegions int  `json:"totalRegions"`
		Regions      []struct {
			RegionID   int    `json:"regionId"`
			RegionName string `json:"regionName"`
			RegionCode string `json:"regionCode"`
			Nodes      []struct {
				Name     string `json:"name"`
				HostName string `json:"hostName"`
				DERPPort int    `json:"derpPort"`
				STUNPort int    `json:"stunPort"`
				IPv4     string `json:"ipv4"`
				IPv6     string `json:"ipv6"`
			} `json:"nodes"`
		} `json:"regions"`
	}
	require.NoError(t, json.Unmarshal(body, &derp))

	assert.True(t, derp.Configured)
	assert.Equal(t, 1, derp.TotalRegions)
	require.Len(t, derp.Regions, 1)

	region := derp.Regions[0]
	assert.Equal(t, 900, region.RegionID)
	assert.Equal(t, "Test Region", region.RegionName)
	assert.Equal(t, "test", region.RegionCode)
	require.Len(t, region.Nodes, 1)

	node := region.Nodes[0]
	assert.Equal(t, "test0", node.Name)
	assert.Equal(t, "127.0.0.1", node.HostName)
	assert.Equal(t, "127.0.0.1", node.IPv4)
	assert.Equal(t, -1, node.DERPPort)
}

// TestDERPEndpointRequiresAuth verifies the endpoint rejects unauthenticated
// requests.
func TestDERPEndpointRequiresAuth(t *testing.T) {
	srv := NewServer(t, WithRealListener())

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL+"/api/v1/derp", nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// TestDERPEndpointEmptyMap verifies the endpoint reports an unconfigured DERP
// map without erroring.
func TestDERPEndpointEmptyMap(t *testing.T) {
	srv := NewServer(t, WithRealListener())

	user := srv.CreateUser(t, "derp-empty")
	apiKey := srv.CreateAPIKey(t, user)

	// Replace the default DERP map with an empty one.
	srv.State().SetDERPMap(&tailcfg.DERPMap{Regions: map[tailcfg.DERPRegionID]*tailcfg.DERPRegion{}})

	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, srv.URL+"/api/v1/derp", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var derp struct {
		Configured   bool `json:"configured"`
		TotalRegions int  `json:"totalRegions"`
	}
	require.NoError(t, json.Unmarshal(body, &derp))

	assert.True(t, derp.Configured)
	assert.Equal(t, 0, derp.TotalRegions)
}
