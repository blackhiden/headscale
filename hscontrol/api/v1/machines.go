package apiv1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/arsydoni4326-alt/headscale/hscontrol/scope"
	"github.com/arsydoni4326-alt/headscale/hscontrol/util/zlog/zf"
	"github.com/rs/zerolog/log"
)

func init() {
	registrations = append(registrations, registerMachines)
}

// requiredApprovalScope is the OAuth scope a scope-limited (non-admin) caller
// must hold to approve a machine. An admin API key is all-access and bypasses
// the scope check.
const requiredApprovalScope = scope.DevicesCore

var (
	// errMachineNotFound is wrapped with the requested node ID in the
	// single-approval 404 response.
	errMachineNotFound = errors.New("machine not found")
	// errNoNodeIDs rejects a bulk request with an empty nodeIds array.
	errNoNodeIDs = errors.New("nodeIds array cannot be empty")
)

// actorFromContext names the caller for audit logging: "admin_api_key" for an
// all-access key or locally-trusted request, "oauth_token" for a scope-limited
// bearer token. It never discloses key material.
func actorFromContext(ctx context.Context) string {
	if isAdmin(ctx) {
		return "admin_api_key"
	}

	if _, ok := principalScopes(ctx); ok {
		return "oauth_token"
	}

	return "unknown"
}

// checkApprovalPermission enforces RBAC for the approval operations: the caller
// must be an all-access admin (API key or locally-trusted request) or hold an
// OAuth token granting [requiredApprovalScope]. The auth middleware has already
// validated the credential and, for a token, rejected one missing the scope
// declared via [requireScope]; this is the defence-in-depth check handlers run
// themselves so the rule holds even if an operation's scope metadata is
// missing.
func checkApprovalPermission(ctx context.Context) error {
	if isAdmin(ctx) {
		return nil
	}

	scopes, _ := principalScopes(ctx)
	if scope.Grants(scope.Parse(scopes), requiredApprovalScope) {
		return nil
	}

	return huma.Error403Forbidden("approving machines requires admin rights or the " +
		string(requiredApprovalScope) + " scope")
}

// ApproveMachineRequestBody is the request body for approving a single machine.
type ApproveMachineRequestBody struct {
	NodeID string `doc:"Node ID to approve"         example:"123"    json:"nodeId"`
	UserID string `doc:"Optional user ID to assign" example:"user-1" json:"userId,omitempty"`
}

// ApproveMachinesRequestBody is the request body for bulk machine approval.
type ApproveMachinesRequestBody struct {
	NodeIDs []string `doc:"List of node IDs to approve"             example:"[\"123\",\"456\"]" json:"nodeIds"`
	UserID  string   `doc:"Optional user ID to assign to all nodes" example:"user-1"            json:"userId,omitempty"`
}

// ApproveMachineOutputBody is the response for approving a single machine.
type ApproveMachineOutputBody struct {
	Success bool  `json:"success"`
	Node    *Node `json:"node"`
}

// BulkApprovalResult represents the result of a single approval in bulk operation.
type BulkApprovalResult struct {
	NodeID  string `json:"nodeId"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// ApproveMachinesOutputBody is the response for bulk machine approval.
type ApproveMachinesOutputBody struct {
	Success  bool                 `json:"success"`
	Approved []string             `json:"approved"`
	Failed   []string             `json:"failed"`
	Errors   map[string]string    `json:"errors"`
	Results  []BulkApprovalResult `json:"results"`
}

type (
	approveMachineInput struct {
		NodeID string `doc:"Node ID" path:"id"`
		Body   ApproveMachineRequestBody
	}
	approveMachineOutput struct {
		Body ApproveMachineOutputBody
	}
)

type (
	approveMachinesInput struct {
		Body ApproveMachinesRequestBody
	}
	approveMachinesOutput struct {
		Body ApproveMachinesOutputBody
	}
)

func registerMachines(api huma.API, b Backend) {
	huma.Register(api, requireScope(huma.Operation{
		OperationID: "approveMachine",
		Method:      http.MethodPost,
		Path:        "/api/v1/machines/{id}/approve",
		Summary:     "Approve a single machine",
		Description: "Approves a pending machine by clearing its expiry, allowing it to join the network.",
		Tags:        []string{"Machines"},
		Security:    bearerAuth,
	}, requiredApprovalScope), func(ctx context.Context, in *approveMachineInput) (*approveMachineOutput, error) {
		err := checkApprovalPermission(ctx)
		if err != nil {
			return nil, err
		}

		// Parse node ID
		nodeID, err := parseNodeID(in.NodeID)
		if err != nil {
			return nil, err
		}

		// Get the node to check if it exists
		node, exists := b.State.GetNodeByID(nodeID)
		if !exists {
			return nil, huma.Error404NotFound("node not found",
				fmt.Errorf("%w: %d", errMachineNotFound, nodeID))
		}

		// Approve the machine by clearing its expiry (set to nil = never expires)
		approvedNode, c, err := b.State.SetNodeExpiry(nodeID, nil)
		if err != nil {
			return nil, mapError("approving machine", err)
		}

		// Emit change notification
		b.Change(c)

		// Audit log: structured logging with all required fields
		log.Info().
			Str(zf.AuditAction, "machine_approve").
			Str(zf.AuditActor, actorFromContext(ctx)).
			Str(zf.AuditResource, fmt.Sprintf("node:%d", nodeID)).
			Uint64(zf.NodeID, nodeID.Uint64()).
			Str(zf.NodeName, node.Hostname()).
			Time(zf.Timestamp, time.Now()).
			Msg("Machine approved via API")

		out := &approveMachineOutput{}
		out.Body.Success = true
		out.Body.Node = new(nodeFromView(approvedNode))

		return out, nil
	})

	registerBulkApproval(api, b)
}

func registerBulkApproval(api huma.API, b Backend) {
	huma.Register(api, requireScope(huma.Operation{
		OperationID: "approveMachines",
		Method:      http.MethodPost,
		Path:        "/api/v1/machines/approve",
		Summary:     "Approve multiple machines",
		Description: "Approves multiple pending machines in a single operation. Returns success/failure status for each machine.",
		Tags:        []string{"Machines"},
		Security:    bearerAuth,
	}, requiredApprovalScope), func(ctx context.Context, in *approveMachinesInput) (*approveMachinesOutput, error) {
		err := checkApprovalPermission(ctx)
		if err != nil {
			return nil, err
		}

		if len(in.Body.NodeIDs) == 0 {
			return nil, huma.Error400BadRequest("bulk approval", errNoNodeIDs)
		}

		out := &approveMachinesOutput{}
		out.Body.Approved = []string{}
		out.Body.Failed = []string{}
		out.Body.Errors = make(map[string]string)
		out.Body.Results = []BulkApprovalResult{}

		actor := actorFromContext(ctx)

		// Process each node independently: a failure for one machine must not
		// abort the rest, so per-node errors are collected and reported.
		for _, nodeIDStr := range in.Body.NodeIDs {
			result := BulkApprovalResult{NodeID: nodeIDStr, Success: false}

			nodeErr := approveOneMachine(b, nodeIDStr, actor)
			if nodeErr != nil {
				result.Error = nodeErr.Error()

				out.Body.Failed = append(out.Body.Failed, nodeIDStr)
				out.Body.Errors[nodeIDStr] = result.Error
				out.Body.Results = append(out.Body.Results, result)

				continue
			}

			result.Success = true

			out.Body.Approved = append(out.Body.Approved, nodeIDStr)
			out.Body.Results = append(out.Body.Results, result)
		}

		// Overall success if at least one approval succeeded.
		out.Body.Success = len(out.Body.Approved) > 0

		return out, nil
	})
}

// approveOneMachine parses nodeIDStr, approves that machine, emits the change
// notification, and writes an audit log entry. It returns a descriptive error
// when the machine is unknown or cannot be approved, which the bulk handler
// records against the node id.
func approveOneMachine(b Backend, nodeIDStr, actor string) error {
	nodeID, err := parseNodeID(nodeIDStr)
	if err != nil {
		return fmt.Errorf("%w: %w", errMachineNotFound, err)
	}

	node, exists := b.State.GetNodeByID(nodeID)
	if !exists {
		return fmt.Errorf("%w: %d", errMachineNotFound, nodeID)
	}

	_, c, err := b.State.SetNodeExpiry(nodeID, nil)
	if err != nil {
		return fmt.Errorf("approving machine: %w", err)
	}

	b.Change(c)

	log.Info().
		Str(zf.AuditAction, "machine_approve_bulk").
		Str(zf.AuditActor, actor).
		Str(zf.AuditResource, fmt.Sprintf("node:%d", nodeID)).
		Uint64(zf.NodeID, nodeID.Uint64()).
		Str(zf.NodeName, node.Hostname()).
		Time(zf.Timestamp, time.Now()).
		Msg("Machine approved via bulk API")

	return nil
}
