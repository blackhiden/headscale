package hscontrol

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	apiv1 "github.com/arsydoni4326-alt/headscale/hscontrol/api/v1"
	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAPIV1ApproveMachine tests single machine approval.
func TestAPIV1ApproveMachine(t *testing.T) {
	t.Run("approve expired node", func(t *testing.T) {
		h := newAPIV1Harness(t)
		nodeID := newNodeSeed("alice", "node-a").register(t, h.app)

		// Expire the node
		expiry := time.Now().Add(-1 * time.Hour)
		_, _, err := h.app.state.SetNodeExpiry(nodeID, &expiry)
		require.NoError(t, err)

		// Verify node is expired
		node, exists := h.app.state.GetNodeByID(nodeID)
		require.True(t, exists)
		assert.True(t, node.IsExpired(), "node should be expired before approval")

		// Approve via API
		res := h.callHuma(http.MethodPost, fmt.Sprintf("/api/v1/machines/%d/approve", nodeID),
			[]byte(fmt.Sprintf(`{"nodeId":"%d"}`, nodeID)))

		require.Equal(t, http.StatusOK, res.status)
		assert.Contains(t, string(res.body), `"success":true`)

		// Verify node is no longer expired
		approvedNode, exists := h.app.state.GetNodeByID(nodeID)
		require.True(t, exists)
		assert.False(t, approvedNode.IsExpired(), "node should not be expired after approval")
	})

	t.Run("not found", func(t *testing.T) {
		h := newAPIV1Harness(t)
		res := h.callHuma(http.MethodPost, "/api/v1/machines/99999/approve",
			[]byte(`{"nodeId":"99999"}`))
		assertStatus(t, res, http.StatusNotFound)
	})

	t.Run("invalid id", func(t *testing.T) {
		h := newAPIV1Harness(t)
		res := h.callHuma(http.MethodPost, "/api/v1/machines/invalid/approve",
			[]byte(`{"nodeId":"invalid"}`))
		assertStatus(t, res, http.StatusBadRequest)
	})
}

// TestAPIV1ApproveMachines tests bulk machine approval.
func TestAPIV1ApproveMachines(t *testing.T) {
	t.Run("bulk approve success", func(t *testing.T) {
		h := newAPIV1Harness(t)

		// Create multiple expired nodes
		node1 := newNodeSeed("alice", "node-1").register(t, h.app)
		node2 := newNodeSeed("bob", "node-2").register(t, h.app)
		node3 := newNodeSeed("carol", "node-3").register(t, h.app)

		expiry := time.Now().Add(-1 * time.Hour)
		for _, nodeID := range []types.NodeID{node1, node2, node3} {
			_, _, err := h.app.state.SetNodeExpiry(nodeID, &expiry)
			require.NoError(t, err)
		}

		// Approve all via bulk API
		body := fmt.Sprintf(`{"nodeIds":["%d","%d","%d"]}`, node1, node2, node3)
		res := h.callHuma(http.MethodPost, "/api/v1/machines/approve", []byte(body))

		require.Equal(t, http.StatusOK, res.status)
		assert.Contains(t, string(res.body), `"success":true`)
		assert.Contains(t, string(res.body), fmt.Sprintf(`"%d"`, node1))
		assert.Contains(t, string(res.body), fmt.Sprintf(`"%d"`, node2))
		assert.Contains(t, string(res.body), fmt.Sprintf(`"%d"`, node3))

		// Verify all nodes are no longer expired
		for _, nodeID := range []types.NodeID{node1, node2, node3} {
			approvedNode, exists := h.app.state.GetNodeByID(nodeID)
			require.True(t, exists)
			assert.False(t, approvedNode.IsExpired(), "node %d should not be expired", nodeID)
		}
	})

	t.Run("partial failure", func(t *testing.T) {
		h := newAPIV1Harness(t)
		node1 := newNodeSeed("alice", "node-1").register(t, h.app)

		expiry := time.Now().Add(-1 * time.Hour)
		_, _, err := h.app.state.SetNodeExpiry(node1, &expiry)
		require.NoError(t, err)

		// Mix valid and invalid IDs
		body := fmt.Sprintf(`{"nodeIds":["%d","99999","invalid"]}`, node1)
		res := h.callHuma(http.MethodPost, "/api/v1/machines/approve", []byte(body))

		require.Equal(t, http.StatusOK, res.status)
		assert.Contains(t, string(res.body), `"success":true`)
		assert.Contains(t, string(res.body), fmt.Sprintf(`"%d"`, node1))

		// Valid node should be approved
		approvedNode, exists := h.app.state.GetNodeByID(node1)
		require.True(t, exists)
		assert.False(t, approvedNode.IsExpired())
	})

	t.Run("empty array", func(t *testing.T) {
		h := newAPIV1Harness(t)
		res := h.callHuma(http.MethodPost, "/api/v1/machines/approve", []byte(`{"nodeIds":[]}`))
		assertStatus(t, res, http.StatusBadRequest)
	})

	t.Run("all failed", func(t *testing.T) {
		h := newAPIV1Harness(t)
		res := h.callHuma(http.MethodPost, "/api/v1/machines/approve",
			[]byte(`{"nodeIds":["99998","99999","invalid"]}`))

		require.Equal(t, http.StatusOK, res.status)
		assert.Contains(t, string(res.body), `"success":false`)
	})
}

// TestAPIV1ApproveMachineRBAC exercises the RBAC layer on the real router
// (no local-trust bypass): an admin API key may approve, an OAuth token with
// the devices:core scope may approve, a read-only token is 403, a missing
// credential is 401, and an invalid key is 401.
func TestAPIV1ApproveMachineRBAC(t *testing.T) {
	app := createTestApp(t)
	handler := app.HTTPHandler()

	nodeID := newNodeSeed("alice", "node-a").register(t, app)
	body := []byte(fmt.Sprintf(`{"nodeId":"%d"}`, nodeID))
	path := fmt.Sprintf("/api/v1/machines/%d/approve", nodeID)

	// Admin API key (all-access).
	expiry := time.Now().Add(time.Hour)
	adminKey, _, err := app.state.CreateAPIKey(&expiry)
	require.NoError(t, err)

	// OAuth client with the write scope, minting a usable access token.
	_, scopedClient, err := app.state.CreateOAuthClient(
		[]string{"devices:core"}, []string{"tag:ci"}, "approver", nil,
	)
	require.NoError(t, err)

	tokExpiry := time.Now().Add(time.Hour)
	scopedToken, _, err := app.state.MintAccessToken(
		scopedClient.ClientID, []string{"devices:core"}, nil, &tokExpiry,
	)
	require.NoError(t, err)

	// OAuth client whose token is read-only: cannot approve.
	_, readonlyClient, err := app.state.CreateOAuthClient(
		[]string{"devices:core:read"}, []string{"tag:ci"}, "reader", nil,
	)
	require.NoError(t, err)

	readonlyToken, _, err := app.state.MintAccessToken(
		readonlyClient.ClientID, []string{"devices:core:read"}, nil, &tokExpiry,
	)
	require.NoError(t, err)

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{name: "missing credential", authHeader: "", wantStatus: http.StatusUnauthorized},
		{name: "invalid key", authHeader: "Bearer hskey-api-nope-nope", wantStatus: http.StatusUnauthorized},
		{name: "admin api key allowed", authHeader: "Bearer " + adminKey, wantStatus: http.StatusOK},
		{name: "scoped token allowed", authHeader: "Bearer " + scopedToken, wantStatus: http.StatusOK},
		{name: "read-only token forbidden", authHeader: "Bearer " + readonlyToken, wantStatus: http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := callHandlerAuth(handler, http.MethodPost, path, tt.authHeader, body)
			assert.Equalf(t, tt.wantStatus, res.status, "body: %s", res.body)
		})
	}
}

// TestAPIV1ApproveMachinesRBAC proves the bulk endpoint enforces the same RBAC
// rule: a read-only token is rejected before any node is touched.
func TestAPIV1ApproveMachinesRBAC(t *testing.T) {
	app := createTestApp(t)
	handler := app.HTTPHandler()

	nodeID := newNodeSeed("alice", "node-a").register(t, app)
	expiry := time.Now().Add(-time.Hour)
	_, _, err := app.state.SetNodeExpiry(nodeID, &expiry)
	require.NoError(t, err)

	body := []byte(fmt.Sprintf(`{"nodeIds":["%d"]}`, nodeID))

	_, readonlyClient, err := app.state.CreateOAuthClient(
		[]string{"devices:core:read"}, []string{"tag:ci"}, "reader", nil,
	)
	require.NoError(t, err)

	tokExpiry := time.Now().Add(time.Hour)
	readonlyToken, _, err := app.state.MintAccessToken(
		readonlyClient.ClientID, []string{"devices:core:read"}, nil, &tokExpiry,
	)
	require.NoError(t, err)

	res := callHandlerAuth(handler, http.MethodPost, "/api/v1/machines/approve",
		"Bearer "+readonlyToken, body)
	assert.Equalf(t, http.StatusForbidden, res.status, "body: %s", res.body)

	// The node must remain expired: the forbidden request changed nothing.
	node, ok := app.state.GetNodeByID(nodeID)
	require.True(t, ok)
	assert.True(t, node.IsExpired(), "read-only token must not have approved the node")
}

// TestAPIV1ApproveMachinePersists proves the approval is written through to the
// database, not just the in-memory NodeStore: after the API call the persisted
// node row has its expiry cleared, so a server restart keeps the machine
// approved.
func TestAPIV1ApproveMachinePersists(t *testing.T) {
	h := newAPIV1Harness(t)
	nodeID := newNodeSeed("alice", "node-a").register(t, h.app)

	expiry := time.Now().Add(-1 * time.Hour)
	_, _, err := h.app.state.SetNodeExpiry(nodeID, &expiry)
	require.NoError(t, err)

	// The expiry must be persisted before approval.
	var before types.Node
	require.NoError(t, h.app.state.DB().DB.First(&before, nodeID.Uint64()).Error)
	require.NotNil(t, before.Expiry, "node should carry a persisted expiry before approval")

	res := h.callHuma(http.MethodPost, fmt.Sprintf("/api/v1/machines/%d/approve", nodeID),
		[]byte(fmt.Sprintf(`{"nodeId":"%d"}`, nodeID)))
	require.Equal(t, http.StatusOK, res.status)

	// The approval must be persisted: the DB expiry is cleared.
	var after types.Node
	require.NoError(t, h.app.state.DB().DB.First(&after, nodeID.Uint64()).Error)
	assert.Nil(t, after.Expiry, "approval must clear the persisted expiry")
}

// TestAPIV1ApproveMachineInSpec proves both approval operations appear in the
// emitted OpenAPI document alongside the rest of the v1 API, so the fork's API
// surface stays fully documented, and that the required OAuth scope is
// published as an x-required-scope extension.
func TestAPIV1ApproveMachineInSpec(t *testing.T) {
	spec, err := apiv1.Spec()
	require.NoError(t, err)

	s := string(spec)
	assert.Contains(t, s, "/api/v1/machines/{id}/approve")
	assert.Contains(t, s, "/api/v1/machines/approve")
	assert.Contains(t, s, "approveMachine")
	assert.Contains(t, s, "approveMachines")
	assert.Contains(t, s, "ApproveMachineOutputBody")
	assert.Contains(t, s, "ApproveMachinesOutputBody")
	assert.Contains(t, s, "x-required-scope")
}
