package v2

import (
	"net/netip"
	"testing"

	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
	"github.com/arsydoni4326-alt/headscale/hscontrol/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"tailscale.com/tailcfg"
)

// srcIPsContain reports whether any rule's SrcIPs cover the given IP.
// SrcIPs may be single IPs, CIDRs, or ranges (e.g. 100.64.0.1-100.64.0.2).
func srcIPsContain(rules []tailcfg.FilterRule, ip string) bool {
	addr := netip.MustParseAddr(ip)

	for _, rule := range rules {
		for _, srcIP := range rule.SrcIPs {
			ipSet, err := util.ParseIPSet(srcIP, nil)
			if err != nil {
				continue
			}

			if ipSet.Contains(addr) {
				return true
			}
		}
	}

	return false
}

// dstPortsContain reports whether any rule's DstPorts cover the given IP.
func dstPortsContain(rules []tailcfg.FilterRule, ip string) bool {
	addr := netip.MustParseAddr(ip)

	for _, rule := range rules {
		for _, dp := range rule.DstPorts {
			ipSet, err := util.ParseIPSet(dp.IP, nil)
			if err != nil {
				continue
			}

			if ipSet.Contains(addr) {
				return true
			}
		}
	}

	return false
}

// TestOIDCGroupInACL verifies that a group referenced in an ACL that is
// not defined in the policy file but exists as an OIDC group on a user
// is accepted and resolves to that user's nodes.
func TestOIDCGroupInACL(t *testing.T) {
	users := types.Users{
		{
			ID:         1,
			Name:       "alice",
			Email:      "alice@example.com",
			OIDCGroups: []string{"group:engineering"},
		},
		{
			ID:         2,
			Name:       "bob",
			Email:      "bob@example.com",
			OIDCGroups: []string{"group:engineering"},
		},
		{
			ID:    3,
			Name:  "carol",
			Email: "carol@example.com",
		},
	}

	nodes := types.Nodes{
		node("alice-laptop", "100.64.0.1", "fd7a:115c:a1e0::1", users[0]),
		node("bob-laptop", "100.64.0.2", "fd7a:115c:a1e0::2", users[1]),
		node("carol-laptop", "100.64.0.3", "fd7a:115c:a1e0::3", users[2]),
	}

	// OIDC group referenced in ACL src without being defined in policy.
	policy := `{
		"acls": [
			{
				"action": "accept",
				"src": ["group:engineering"],
				"dst": ["100.64.0.3:80"]
			}
		]
	}`

	pm, err := NewPolicyManager([]byte(policy), users, nodes.ViewSlice())
	require.NoError(t, err, "policy with OIDC group should parse and validate")

	filter, _ := pm.Filter()
	require.NotNil(t, filter, "filter should not be nil")

	// The filter should contain the IPs of alice and bob (OIDC group members).
	assert.True(t, srcIPsContain(filter, "100.64.0.1"), "alice's IP should be in the filter for OIDC group")
	assert.True(t, srcIPsContain(filter, "100.64.0.2"), "bob's IP should be in the filter for OIDC group")
	assert.False(t, srcIPsContain(filter, "100.64.0.3"), "carol's IP should NOT be in the filter (not in OIDC group)")
}

// TestOIDCGroupInACLDst verifies OIDC groups work as ACL destinations.
func TestOIDCGroupInACLDst(t *testing.T) {
	users := types.Users{
		{
			ID:         1,
			Name:       "alice",
			Email:      "alice@example.com",
			OIDCGroups: []string{"group:engineering"},
		},
	}

	nodes := types.Nodes{
		node("alice-laptop", "100.64.0.1", "fd7a:115c:a1e0::1", users[0]),
	}

	policy := `{
		"acls": [
			{
				"action": "accept",
				"src": ["autogroup:member"],
				"dst": ["group:engineering:*"]
			}
		]
	}`

	pm, err := NewPolicyManager([]byte(policy), users, nodes.ViewSlice())
	require.NoError(t, err, "policy with OIDC group destination should parse and validate")

	filter, _ := pm.Filter()
	require.NotNil(t, filter, "filter should not be nil")

	assert.True(t, dstPortsContain(filter, "100.64.0.1"), "alice's IP should be in the filter destinations for OIDC group")
}

// TestOIDCGroupInTagOwners verifies OIDC groups work as tag owners.
func TestOIDCGroupInTagOwners(t *testing.T) {
	users := types.Users{
		{
			ID:         1,
			Name:       "alice",
			Email:      "alice@example.com",
			OIDCGroups: []string{"group:infra"},
		},
	}

	nodes := types.Nodes{
		node("alice-laptop", "100.64.0.1", "fd7a:115c:a1e0::1", users[0]),
	}

	policy := `{
		"tagOwners": {
			"tag:server": ["group:infra"]
		}
	}`

	pm, err := NewPolicyManager([]byte(policy), users, nodes.ViewSlice())
	require.NoError(t, err, "policy with OIDC group tag owner should parse and validate")

	// Verify the tag owner resolves to alice's node.
	canHaveTag := pm.NodeCanHaveTag(nodes[0].View(), "tag:server")
	assert.True(t, canHaveTag, "alice should be able to own tag:server via OIDC group")
}

// TestOIDCGroupInAutoApprovers verifies OIDC groups work as auto-approvers.
func TestOIDCGroupInAutoApprovers(t *testing.T) {
	users := types.Users{
		{
			ID:         1,
			Name:       "alice",
			Email:      "alice@example.com",
			OIDCGroups: []string{"group:network"},
		},
	}

	nodes := types.Nodes{
		node("alice-laptop", "100.64.0.1", "fd7a:115c:a1e0::1", users[0]),
	}

	policy := `{
		"autoApprovers": {
			"routes": {
				"10.0.0.0/16": ["group:network"]
			}
		}
	}`

	pm, err := NewPolicyManager([]byte(policy), users, nodes.ViewSlice())
	require.NoError(t, err, "policy with OIDC group auto-approver should parse and validate")

	// Verify the auto-approver resolves to alice's node.
	canApprove := pm.NodeCanApproveRoute(nodes[0].View(), netip.MustParsePrefix("10.0.0.0/16"))
	assert.True(t, canApprove, "alice should be able to approve routes via OIDC group")
}

// TestOIDCGroupUndefined verifies that a group that is neither defined in
// the policy nor present as an OIDC group on any user is rejected.
func TestOIDCGroupUndefined(t *testing.T) {
	users := types.Users{
		{
			ID:    1,
			Name:  "alice",
			Email: "alice@example.com",
		},
	}

	nodes := types.Nodes{
		node("alice-laptop", "100.64.0.1", "fd7a:115c:a1e0::1", users[0]),
	}

	policy := `{
		"acls": [
			{
				"action": "accept",
				"src": ["group:notdefined"],
				"dst": ["100.64.0.1:80"]
			}
		]
	}`

	_, err := NewPolicyManager([]byte(policy), users, nodes.ViewSlice())
	require.Error(t, err, "policy with undefined group should fail validation")
	assert.ErrorIs(t, err, ErrGroupNotDefined)
}

// TestOIDCGroupPolicyDefinedTakesPrecedence verifies that when a group is
// both defined in the policy and present as an OIDC group, the policy
// definition is used (existing behavior preserved).
func TestOIDCGroupPolicyDefinedTakesPrecedence(t *testing.T) {
	users := types.Users{
		{
			ID:         1,
			Name:       "alice",
			Email:      "alice@example.com",
			OIDCGroups: []string{"group:engineering"},
		},
		{
			ID:         2,
			Name:       "bob",
			Email:      "bob@example.com",
			OIDCGroups: []string{"group:engineering"},
		},
	}

	nodes := types.Nodes{
		node("alice-laptop", "100.64.0.1", "fd7a:115c:a1e0::1", users[0]),
		node("bob-laptop", "100.64.0.2", "fd7a:115c:a1e0::2", users[1]),
	}

	// Policy defines group:engineering as only bob, but both users have it
	// as an OIDC group. The policy definition should take precedence.
	policy := `{
		"groups": {
			"group:engineering": ["bob@example.com"]
		},
		"acls": [
			{
				"action": "accept",
				"src": ["group:engineering"],
				"dst": ["100.64.0.1:80"]
			}
		]
	}`

	pm, err := NewPolicyManager([]byte(policy), users, nodes.ViewSlice())
	require.NoError(t, err)

	filter, _ := pm.Filter()
	require.NotNil(t, filter)

	// Policy-defined group only includes bob, so alice should NOT be in the filter.
	assert.False(t, srcIPsContain(filter, "100.64.0.1"), "alice should not be in the filter (policy group only includes bob)")
	assert.True(t, srcIPsContain(filter, "100.64.0.2"), "bob should be in the filter (policy group includes bob)")
}
