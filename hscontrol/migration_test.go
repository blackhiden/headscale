package hscontrol

import (
	"testing"
)

// Phase 13c: Migration tests removed as part of Headplane local-auth retirement.
// TestMigrationFromSingleToMultiUser, TestDatabaseSchema, and TestBackwardCompatibility
// were testing the now-retired local password authentication and user management system.
// Legacy tables (headplane_users, headplane_settings) are preserved but no longer
// actively managed by Headscale.
//
// Normal Headscale API key authentication and node management are unaffected and
// tested elsewhere in the test suite.

func TestPlaceholder(t *testing.T) {
	// Placeholder to keep the file valid during Phase 13c rollback window.
	// This file can be removed in a future cleanup.
	t.Skip("Phase 13c: Headplane local-auth migration tests retired")
}

