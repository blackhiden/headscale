package hscontrol

import (
	"encoding/json"
	"net/http"
	"testing"

	apiv1 "github.com/arsydoni4326-alt/headscale/hscontrol/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tailscale.com/tailcfg"
)

// TestAPIV1Derp proves the read-only GET /api/v1/derp endpoint returns the
// configured DERP map with its regions and nodes.
func TestAPIV1Derp(t *testing.T) {
	h := newAPIV1Harness(t)

	t.Run("unconfigured derp map", func(t *testing.T) {
		res := h.callHuma(http.MethodGet, "/api/v1/derp", nil)

		assert.Equal(t, http.StatusOK, res.status)

		var body map[string]any
		require.NoError(t, json.Unmarshal(res.body, &body))

		configured, ok := body["configured"].(bool)
		require.True(t, ok, "expected configured field")
		assert.False(t, configured)

		totalRegions, ok := body["totalRegions"].(float64)
		require.True(t, ok, "expected totalRegions field")
		assert.Equal(t, float64(0), totalRegions)
	})

	t.Run("configured derp map", func(t *testing.T) {
		h.app.state.SetDERPMap(&tailcfg.DERPMap{
			Regions: map[tailcfg.DERPRegionID]*tailcfg.DERPRegion{
				900: {
					RegionID:   900,
					RegionCode: "test",
					RegionName: "Test Region",
					Nodes: []*tailcfg.DERPNode{{
						Name:     "test0",
						RegionID: 900,
						HostName: "127.0.0.1",
						IPv4:     "127.0.0.1",
						DERPPort: 8766,
						STUNPort: 3478,
					}},
				},
			},
		})

		res := h.callHuma(http.MethodGet, "/api/v1/derp", nil)

		assert.Equal(t, http.StatusOK, res.status)

		var body map[string]any
		require.NoError(t, json.Unmarshal(res.body, &body))

		configured, ok := body["configured"].(bool)
		require.True(t, ok, "expected configured field")
		assert.True(t, configured)

		totalRegions, ok := body["totalRegions"].(float64)
		require.True(t, ok, "expected totalRegions field")
		assert.Equal(t, float64(1), totalRegions)

		regions, ok := body["regions"].([]any)
		require.True(t, ok, "expected regions array")
		require.Len(t, regions, 1)

		region, ok := regions[0].(map[string]any)
		require.True(t, ok, "expected region object")
		assert.Equal(t, float64(900), region["regionId"])
		assert.Equal(t, "test", region["regionCode"])
		assert.Equal(t, "Test Region", region["regionName"])

		nodes, ok := region["nodes"].([]any)
		require.True(t, ok, "expected nodes array")
		require.Len(t, nodes, 1)

		node, ok := nodes[0].(map[string]any)
		require.True(t, ok, "expected node object")
		assert.Equal(t, "test0", node["name"])
		assert.Equal(t, "127.0.0.1", node["hostName"])
		assert.Equal(t, float64(8766), node["derpPort"])
		assert.Equal(t, float64(3478), node["stunPort"])
		assert.Equal(t, "127.0.0.1", node["ipv4"])
	})
}

// TestAPIV1DerpInSpec proves the derp endpoint appears in the emitted OpenAPI
// document alongside the rest of the v1 API, so the fork's API surface stays
// fully documented.
func TestAPIV1DerpInSpec(t *testing.T) {
	spec, err := apiv1.Spec()
	require.NoError(t, err)

	assert.Contains(t, string(spec), "/api/v1/derp")
	assert.Contains(t, string(spec), "DERPResponseBody")
	assert.Contains(t, string(spec), "getDerp")
}
