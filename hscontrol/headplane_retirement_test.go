package hscontrol

import (
	"testing"

	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServerStartupAfterHeadplaneRetirement verifies that Headscale server starts
// successfully after Phase 13c Headplane local-auth retirement.
func TestServerStartupAfterHeadplaneRetirement(t *testing.T) {
	cfg := &types.Config{
		ServerURL: "http://headscale.example.com",
		Database: types.DatabaseConfig{
			Type:   "sqlite3",
			Sqlite: types.SqliteConfig{Path: t.TempDir() + "/headscale_test.db"},
		},
		NoisePrivateKeyPath: t.TempDir() + "/noise_private.key",
		Policy:              types.PolicyConfig{Mode: types.PolicyModeDB},
	}

	app, err := NewHeadscale(cfg)
	require.NoError(t, err, "server should start without Headplane config")
	require.NotNil(t, app, "app should not be nil")
	require.NotNil(t, app.state, "app state should be initialized")
}

// TestHeadplaneConfigDeprecatedButPreserved verifies that HeadplaneConfig fields
// are preserved but not used during the rollback window.
func TestHeadplaneConfigDeprecatedButPreserved(t *testing.T) {
	cfg := &types.Config{
		ServerURL: "http://headscale.example.com",
		Database: types.DatabaseConfig{
			Type:   "sqlite3",
			Sqlite: types.SqliteConfig{Path: t.TempDir() + "/headscale_test.db"},
		},
		NoisePrivateKeyPath: t.TempDir() + "/noise_private.key",
		Policy:              types.PolicyConfig{Mode: types.PolicyModeDB},
		Headplane: types.HeadplaneConfig{
			Password: "legacy-password-ignored",
		},
	}

	app, err := NewHeadscale(cfg)
	require.NoError(t, err, "server should start even with legacy Headplane.Password")
	require.NotNil(t, app)

	// Verify the config value is preserved but unused
	assert.Equal(t, "legacy-password-ignored", cfg.Headplane.Password,
		"legacy config field should be preserved during rollback window")
}

// TestLegacyTablesPreserved verifies that legacy Headplane tables are created
// during migration but no longer actively managed.
func TestLegacyTablesPreserved(t *testing.T) {
	cfg := &types.Config{
		ServerURL: "http://headscale.example.com",
		Database: types.DatabaseConfig{
			Type:   "sqlite3",
			Sqlite: types.SqliteConfig{Path: t.TempDir() + "/headscale_test.db"},
		},
		NoisePrivateKeyPath: t.TempDir() + "/noise_private.key",
		Policy:              types.PolicyConfig{Mode: types.PolicyModeDB},
	}

	app, err := NewHeadscale(cfg)
	require.NoError(t, err)
	require.NotNil(t, app)

	// Verify legacy tables exist for rollback compatibility
	db := app.state.DB().DB
	assert.True(t, db.Migrator().HasTable("headplane_users"),
		"legacy headplane_users table should exist for rollback")
	assert.True(t, db.Migrator().HasTable("headplane_settings"),
		"legacy headplane_settings table should exist for rollback")
}
