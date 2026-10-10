package db

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/netip"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/arsydoni4326-alt/headscale/hscontrol/db/sqliteconfig"
	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
	"github.com/arsydoni4326-alt/headscale/hscontrol/util"
	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/rs/zerolog/log"
	"github.com/tailscale/squibble"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

//go:embed schema.sql
var dbSchema string

// sqliteDialect is the dialector name gorm reports via DB.Name() for SQLite.
// Compared to select dialect-specific DDL (SQLite is the schema source of
// truth; other dialects mirror it).
const sqliteDialect = "sqlite"

const headplaneSchemaRepairMigrationID = "202610050900-repair-headplane-schema"

const headplaneIndexNormalizationMigrationID = "202610050930-normalize-headplane-indexes"

func init() {
	schema.RegisterSerializer("text", TextSerialiser{})
}

var errDatabaseNotSupported = errors.New("database type not supported")

var errForeignKeyConstraintsViolated = errors.New("foreign key constraints violated")

const (
	maxIdleConns   = 100
	maxOpenConns   = 100
	contextTimeout = 10 * time.Second
)

type HSDatabase struct {
	DB  *gorm.DB
	cfg *types.Config
}

// NewHeadscaleDatabase creates a new database connection and runs migrations.
//
//nolint:gocyclo // migration closures inflate the count; each is linear
func NewHeadscaleDatabase(cfg *types.Config) (*HSDatabase, error) {
	dbConn, err := openDB(cfg.Database)
	if err != nil {
		return nil, err
	}

	err = checkMinimumMigration(dbConn)
	if err != nil {
		return nil, fmt.Errorf("version check: %w", err)
	}

	err = checkVersionUpgradePath(dbConn)
	if err != nil {
		return nil, fmt.Errorf("version check: %w", err)
	}

	migrations := gormigrate.New(
		dbConn,
		gormigrate.DefaultOptions,
		[]*gormigrate.Migration{
			// New migrations must be added as transactions at the end of this list.
			// Migrations start from v0.25.0. If upgrading from v0.24.x or earlier,
			// you must first upgrade to v0.25.1 before upgrading to this version.

			// v0.25.0
			{
				// Add a constraint to routes ensuring they cannot exist without a node.
				ID: "202501221827",
				Migrate: func(tx *gorm.DB) error {
					// Remove any invalid routes associated with a node that does not exist.
					if tx.Migrator().HasTable(&types.Route{}) && tx.Migrator().HasTable(&types.Node{}) { //nolint:staticcheck // SA1019: Route kept for migrations
						err := tx.Exec("delete from routes where node_id not in (select id from nodes)").Error
						if err != nil {
							return err
						}
					}

					// Remove any invalid routes without a node_id.
					if tx.Migrator().HasTable(&types.Route{}) { //nolint:staticcheck // SA1019: Route kept for migrations
						err := tx.Exec("delete from routes where node_id is null").Error
						if err != nil {
							return err
						}
					}

					err := tx.AutoMigrate(&types.Route{}) //nolint:staticcheck // SA1019: Route kept for migrations
					if err != nil {
						return fmt.Errorf("automigrating types.Route: %w", err)
					}

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			// Add back constraint so you cannot delete preauth keys that
			// is still used by a node.
			{
				ID: "202501311657",
				Migrate: func(tx *gorm.DB) error {
					err := tx.AutoMigrate(&types.PreAuthKey{})
					if err != nil {
						return fmt.Errorf("automigrating types.PreAuthKey: %w", err)
					}

					err = tx.AutoMigrate(&types.Node{})
					if err != nil {
						return fmt.Errorf("automigrating types.Node: %w", err)
					}

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			// Ensure there are no nodes referring to a deleted preauthkey.
			{
				ID: "202502070949",
				Migrate: func(tx *gorm.DB) error {
					if tx.Migrator().HasTable(&types.PreAuthKey{}) {
						err := tx.Exec(`
UPDATE nodes
SET auth_key_id = NULL
WHERE auth_key_id IS NOT NULL
AND auth_key_id NOT IN (
    SELECT id FROM pre_auth_keys
);
							`).Error
						if err != nil {
							return fmt.Errorf("setting auth_key to null on nodes with non-existing keys: %w", err)
						}
					}

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			// v0.26.0
			// Migrate all routes from the Route table to the new field ApprovedRoutes
			// in the Node table. Then drop the Route table.
			{
				ID: "202502131714",
				Migrate: func(tx *gorm.DB) error {
					if !tx.Migrator().HasColumn(&types.Node{}, "approved_routes") {
						err := tx.Migrator().AddColumn(&types.Node{}, "approved_routes")
						if err != nil {
							return fmt.Errorf("adding column types.Node: %w", err)
						}
					}

					nodeRoutes := map[uint64][]netip.Prefix{}

					var routes []types.Route //nolint:staticcheck // SA1019: Route kept for migrations

					err = tx.Find(&routes).Error
					if err != nil {
						return fmt.Errorf("fetching routes: %w", err)
					}

					for _, route := range routes {
						if route.Enabled {
							nodeRoutes[route.NodeID] = append(nodeRoutes[route.NodeID], route.Prefix)
						}
					}

					for nodeID, routes := range nodeRoutes {
						slices.SortFunc(routes, netip.Prefix.Compare)
						routes = slices.Compact(routes)

						data, _ := json.Marshal(routes)

						err = tx.Model(&types.Node{}).Where("id = ?", nodeID).Update("approved_routes", data).Error
						if err != nil {
							return fmt.Errorf("saving approved routes to new column: %w", err)
						}
					}

					// Drop the old table.
					_ = tx.Migrator().DropTable(&types.Route{}) //nolint:staticcheck // SA1019: Route kept for migrations

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			{
				ID: "202502171819",
				Migrate: func(tx *gorm.DB) error {
					// This migration originally removed the last_seen column
					// from the node table, but it was added back in
					// 202505091439.
					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			// Add back last_seen column to node table.
			{
				ID: "202505091439",
				Migrate: func(tx *gorm.DB) error {
					// Add back last_seen column to node table if it does not exist.
					// This is a workaround for the fact that the last_seen column
					// was removed in the 202502171819 migration, but only for some
					// beta testers.
					if !tx.Migrator().HasColumn(&types.Node{}, "last_seen") {
						_ = tx.Migrator().AddColumn(&types.Node{}, "last_seen")
					}

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			// Fix the provider identifier for users that have a double slash in the
			// provider identifier.
			{
				ID: "202505141324",
				Migrate: func(tx *gorm.DB) error {
					users, err := ListUsers(tx, nil)
					if err != nil {
						return fmt.Errorf("listing users: %w", err)
					}

					for _, user := range users {
						cleaned := types.CleanIdentifier(user.ProviderIdentifier.String)

						// Update only the provider_identifier column. Using
						// Save() would write every field on the struct, which
						// breaks when later migrations add new columns (e.g.
						// oidc_groups) that do not exist yet at this point in
						// the migration chain.
						err := tx.Model(&types.User{}).
							Where("id = ?", user.ID).
							Update("provider_identifier", cleaned).Error
						if err != nil {
							return fmt.Errorf("saving user: %w", err)
						}
					}

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			// v0.27.0
			// Schema migration to ensure all tables match the expected schema.
			// This migration recreates all tables to match the exact structure in schema.sql,
			// preserving all data during the process.
			// Only SQLite will be migrated for consistency.
			{
				ID: "202507021200",
				Migrate: func(tx *gorm.DB) error {
					// Only run on SQLite
					if cfg.Database.Type != types.DatabaseSqlite {
						log.Info().Msg("skipping schema migration on non-SQLite database")
						return nil
					}

					log.Info().Msg("starting schema recreation with table renaming")

					// Rename existing tables to _old versions
					tablesToRename := []string{"users", "pre_auth_keys", "api_keys", "nodes", "policies"}

					// Check if routes table exists and drop it (should have been migrated already)
					var routesExists bool

					err := tx.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='routes'").Row().Scan(&routesExists)
					if err == nil && routesExists {
						log.Info().Msg("dropping leftover routes table")

						err := tx.Exec("DROP TABLE routes").Error
						if err != nil {
							return fmt.Errorf("dropping routes table: %w", err)
						}
					}

					// Drop all indexes first to avoid conflicts
					indexesToDrop := []string{
						"idx_users_deleted_at",
						"idx_provider_identifier",
						"idx_name_provider_identifier",
						"idx_name_no_provider_identifier",
						"idx_api_keys_prefix",
						"idx_policies_deleted_at",
					}

					for _, index := range indexesToDrop {
						_ = tx.Exec("DROP INDEX IF EXISTS " + index).Error
					}

					for _, table := range tablesToRename {
						// Check if table exists before renaming
						var exists bool

						err := tx.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Row().Scan(&exists)
						if err != nil {
							return fmt.Errorf("checking if table %s exists: %w", table, err)
						}

						if exists {
							// Drop old table if it exists from previous failed migration
							_ = tx.Exec("DROP TABLE IF EXISTS " + table + "_old").Error

							// Rename current table to _old
							err := tx.Exec("ALTER TABLE " + table + " RENAME TO " + table + "_old").Error
							if err != nil {
								return fmt.Errorf("renaming table %s to %s_old: %w", table, table, err)
							}
						}
					}

					// Create new tables with correct schema
					tableCreationSQL := []string{
						`CREATE TABLE users(
  id integer PRIMARY KEY AUTOINCREMENT,
  name text,
  display_name text,
  email text,
  provider_identifier text,
  provider text,
  profile_pic_url text,
  created_at datetime,
  updated_at datetime,
  deleted_at datetime
)`,
						`CREATE TABLE pre_auth_keys(
  id integer PRIMARY KEY AUTOINCREMENT,
  key text,
  user_id integer,
  reusable numeric,
  ephemeral numeric DEFAULT false,
  used numeric DEFAULT false,
  tags text,
  expiration datetime,
  created_at datetime,
  CONSTRAINT fk_pre_auth_keys_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL
)`,
						`CREATE TABLE api_keys(
  id integer PRIMARY KEY AUTOINCREMENT,
  prefix text,
  hash blob,
  expiration datetime,
  last_seen datetime,
  created_at datetime
)`,
						`CREATE TABLE nodes(
  id integer PRIMARY KEY AUTOINCREMENT,
  machine_key text,
  node_key text,
  disco_key text,
  endpoints text,
  host_info text,
  ipv4 text,
  ipv6 text,
  hostname text,
  given_name varchar(63),
  user_id integer,
  register_method text,
  forced_tags text,
  auth_key_id integer,
  last_seen datetime,
  expiry datetime,
  approved_routes text,
  created_at datetime,
  updated_at datetime,
  deleted_at datetime,
  CONSTRAINT fk_nodes_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_nodes_auth_key FOREIGN KEY(auth_key_id) REFERENCES pre_auth_keys(id)
)`,
						`CREATE TABLE policies(
  id integer PRIMARY KEY AUTOINCREMENT,
  data text,
  created_at datetime,
  updated_at datetime,
  deleted_at datetime
)`,
					}

					for _, createSQL := range tableCreationSQL {
						err := tx.Exec(createSQL).Error
						if err != nil {
							return fmt.Errorf("creating new table: %w", err)
						}
					}

					// Copy data directly using SQL
					dataCopySQL := []string{
						`INSERT INTO users (id, name, display_name, email, provider_identifier, provider, profile_pic_url, created_at, updated_at, deleted_at)
             SELECT id, name, display_name, email, provider_identifier, provider, profile_pic_url, created_at, updated_at, deleted_at
             FROM users_old`,

						`INSERT INTO pre_auth_keys (id, key, user_id, reusable, ephemeral, used, tags, expiration, created_at)
             SELECT id, key, user_id, reusable, ephemeral, used, tags, expiration, created_at
             FROM pre_auth_keys_old`,

						`INSERT INTO api_keys (id, prefix, hash, expiration, last_seen, created_at)
             SELECT id, prefix, hash, expiration, last_seen, created_at
             FROM api_keys_old`,

						`INSERT INTO nodes (id, machine_key, node_key, disco_key, endpoints, host_info, ipv4, ipv6, hostname, given_name, user_id, register_method, forced_tags, auth_key_id, last_seen, expiry, approved_routes, created_at, updated_at, deleted_at)
             SELECT id, machine_key, node_key, disco_key, endpoints, host_info, ipv4, ipv6, hostname, given_name, user_id, register_method, forced_tags, auth_key_id, last_seen, expiry, approved_routes, created_at, updated_at, deleted_at
             FROM nodes_old`,

						`INSERT INTO policies (id, data, created_at, updated_at, deleted_at)
             SELECT id, data, created_at, updated_at, deleted_at
             FROM policies_old`,
					}

					for _, copySQL := range dataCopySQL {
						err := tx.Exec(copySQL).Error
						if err != nil {
							return fmt.Errorf("copying data: %w", err)
						}
					}

					// Create indexes
					indexes := []string{
						"CREATE INDEX idx_users_deleted_at ON users(deleted_at)",
						`CREATE UNIQUE INDEX idx_provider_identifier ON users(
  provider_identifier
) WHERE provider_identifier IS NOT NULL`,
						`CREATE UNIQUE INDEX idx_name_provider_identifier ON users(
  name,
  provider_identifier
)`,
						`CREATE UNIQUE INDEX idx_name_no_provider_identifier ON users(
  name
) WHERE provider_identifier IS NULL`,
						"CREATE UNIQUE INDEX idx_api_keys_prefix ON api_keys(prefix)",
						"CREATE INDEX idx_policies_deleted_at ON policies(deleted_at)",
					}

					for _, indexSQL := range indexes {
						err := tx.Exec(indexSQL).Error
						if err != nil {
							return fmt.Errorf("creating index: %w", err)
						}
					}

					// Drop old tables only after everything succeeds
					for _, table := range tablesToRename {
						err := tx.Exec("DROP TABLE IF EXISTS " + table + "_old").Error
						if err != nil {
							log.Warn().Str("table", table+"_old").Err(err).Msg("failed to drop old table, but migration succeeded")
						}
					}

					log.Info().Msg("schema recreation completed successfully")

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			// v0.27.1
			{
				// Drop all tables that are no longer in use and has existed.
				// They potentially still present from broken migrations in the past.
				ID: "202510311551",
				Migrate: func(tx *gorm.DB) error {
					for _, oldTable := range []string{"namespaces", "machines", "shared_machines", "kvs", "pre_auth_key_acl_tags", "routes"} {
						err := tx.Migrator().DropTable(oldTable)
						if err != nil {
							log.Trace().Str("table", oldTable).
								Err(err).
								Msg("Error dropping old table, continuing...")
						}
					}

					return nil
				},
				Rollback: func(tx *gorm.DB) error {
					return nil
				},
			},
			{
				// Drop all indices that are no longer in use and has existed.
				// They potentially still present from broken migrations in the past.
				// They should all be cleaned up by the db engine, but we are a bit
				// conservative to ensure all our previous mess is cleaned up.
				ID: "202511101554-drop-old-idx",
				Migrate: func(tx *gorm.DB) error {
					for _, oldIdx := range []struct{ name, table string }{
						{"idx_namespaces_deleted_at", "namespaces"},
						{"idx_routes_deleted_at", "routes"},
						{"idx_shared_machines_deleted_at", "shared_machines"},
					} {
						err := tx.Migrator().DropIndex(oldIdx.table, oldIdx.name)
						if err != nil {
							log.Trace().
								Str("index", oldIdx.name).
								Str("table", oldIdx.table).
								Err(err).
								Msg("Error dropping old index, continuing...")
						}
					}

					return nil
				},
				Rollback: func(tx *gorm.DB) error {
					return nil
				},
			},

			// Migrations **above** this points will be REMOVED in version **0.29.0**
			// This is to clean up a lot of old migrations that is seldom used
			// and carries a lot of technical debt.
			// Any new migrations should be added after the comment below and follow
			// the rules it sets out.

			// Migrations start from v0.29.0; older databases are rejected by
			// checkMinimumMigration and must upgrade to the latest 0.29.x first.
			//
			// Rules:
			// - NEVER use gorm.AutoMigrate, write the exact migration steps needed
			// - AutoMigrate depends on the struct staying exactly the same, which it won't over time.
			// - Never write migrations that requires foreign keys to be disabled.
			// - ALL errors in migrations must be handled properly.
			// Shipped in 0.29.1.
			// TODO(kradalby): remove in 0.31, which upgrades only from 0.30.
			{
				// Recover user_id on untagged nodes detached by the earlier
				// version of 202602201200-clear-tagged-node-user-id, which
				// treated tags='null' as tagged and cleared the user. This
				// repairs databases that already upgraded to 0.29.0; databases
				// that took the fixed migration find nothing to repair.
				// Recovery is best-effort: the owner is re-derived from the
				// node's pre-auth key, so nodes registered via CLI/OIDC (no
				// pre-auth key) cannot be recovered and must be reassigned
				// manually.
				// Fixes: https://github.com/arsydoni4326-alt/headscale/issues/3323
				ID: "202606181200-recover-null-tags-node-user-id",
				Migrate: func(tx *gorm.DB) error {
					err := tx.Exec(`
UPDATE nodes
SET user_id = (
	SELECT pak.user_id FROM pre_auth_keys pak WHERE pak.id = nodes.auth_key_id
)
WHERE user_id IS NULL
	AND auth_key_id IS NOT NULL
	AND (tags IS NULL OR tags = '' OR tags = '[]' OR tags = 'null');
						`).Error
					if err != nil {
						return fmt.Errorf("recovering user_id on untagged nodes: %w", err)
					}

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			// 0.30 development: columns and tables that 202609231300 reads.
			// TODO(kradalby): remove in 0.31 with the credentials migration.
			{
				// Add an optional owning user to API keys so the v2 API can
				// create user-owned (untagged) auth keys, mirroring Tailscale's
				// "key owned by the creating identity".
				ID: "202606191500-api-key-user-id",
				Migrate: func(tx *gorm.DB) error {
					if !tx.Migrator().HasColumn(&types.APIKey{}, "user_id") {
						err := tx.Migrator().AddColumn(&types.APIKey{}, "user_id")
						if err != nil {
							return fmt.Errorf("adding user_id to api_keys: %w", err)
						}
					}

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			{
				// Add a free-text description to pre-auth keys, set via the
				// v2 keys API.
				ID: "202606191501-pre-auth-key-description",
				Migrate: func(tx *gorm.DB) error {
					if !tx.Migrator().HasColumn(&types.PreAuthKey{}, "description") {
						err := tx.Migrator().AddColumn(&types.PreAuthKey{}, "description")
						if err != nil {
							return fmt.Errorf("adding description to pre_auth_keys: %w", err)
						}
					}

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			{
				// Add a revoked timestamp to pre-auth keys. The v2 API's DELETE
				// soft-revokes a key (set revoked = now) rather than destroying
				// it; the row is reaped later by the background collector.
				ID: "202606201200-pre-auth-key-revoked",
				Migrate: func(tx *gorm.DB) error {
					if !tx.Migrator().HasColumn(&types.PreAuthKey{}, "revoked") {
						err := tx.Migrator().AddColumn(&types.PreAuthKey{}, "revoked")
						if err != nil {
							return fmt.Errorf("adding revoked to pre_auth_keys: %w", err)
						}
					}

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			{
				// Add the OAuth client + access token tables backing the v2 API's
				// OAuth client-credentials flow. They mirror the api_keys /
				// pre_auth_keys security model: a public id/prefix plus an Argon2id
				// hash of the secret.
				//
				// SQLite uses explicit DDL that matches schema.sql byte-for-byte
				// (the squibble digest is the SQLite source of truth). Postgres,
				// which has no digest and rejects SQLite-isms like AUTOINCREMENT,
				// uses dialect-aware AutoMigrate, mirroring InitSchema's fresh-DB
				// table creation so an existing Postgres deployment can upgrade.
				ID: "202606211200-oauth-clients-and-tokens",
				Migrate: func(tx *gorm.DB) error {
					if tx.Migrator().HasTable(&types.OAuthClient{}) &&
						tx.Migrator().HasTable(&types.OAuthAccessToken{}) {
						return nil
					}

					if tx.Name() != "sqlite" {
						return tx.AutoMigrate(&types.OAuthClient{}, &types.OAuthAccessToken{})
					}

					if !tx.Migrator().HasTable(&types.OAuthClient{}) {
						err := tx.Exec(`CREATE TABLE oauth_clients(
  id integer PRIMARY KEY AUTOINCREMENT,
  client_id text,
  secret_hash blob,
  scopes text,
  tags text,
  description text,
  user_id integer,
  created_at datetime,
  revoked datetime
)`).Error
						if err != nil {
							return fmt.Errorf("creating oauth_clients table: %w", err)
						}

						err = tx.Exec(`CREATE UNIQUE INDEX idx_oauth_clients_client_id ON oauth_clients(client_id)`).Error
						if err != nil {
							return fmt.Errorf("creating oauth_clients index: %w", err)
						}
					}

					if !tx.Migrator().HasTable(&types.OAuthAccessToken{}) {
						err := tx.Exec(`CREATE TABLE oauth_access_tokens(
  id integer PRIMARY KEY AUTOINCREMENT,
  prefix text,
  hash blob,
  client_id text,
  scopes text,
  tags text,
  expiration datetime,
  created_at datetime
)`).Error
						if err != nil {
							return fmt.Errorf("creating oauth_access_tokens table: %w", err)
						}

						err = tx.Exec(`CREATE UNIQUE INDEX idx_oauth_access_tokens_prefix ON oauth_access_tokens(prefix)`).Error
						if err != nil {
							return fmt.Errorf("creating oauth_access_tokens index: %w", err)
						}
					}

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			// Shipped in 0.29.3.
			// TODO(kradalby): remove in 0.31, which upgrades only from 0.30.
			{
				// Clear stale key expiry on tagged nodes. A tagged node is
				// owned by its tags and never expires (KB 1068), but a buggy
				// handleLogout stamped a past expiry on it, leaving it
				// permanently Expired and unable to re-authenticate. The
				// buggy writer is fixed, so this only repairs rows written
				// before the upgrade; a fixed server cannot recreate them.
				// Match the tagged-node predicate of 0.29's
				// clear-tagged-node-user-id migration (a nil tags slice
				// marshals to 'null', so exclude it).
				// Fixes: https://github.com/arsydoni4326-alt/headscale/issues/3371
				ID: "202607241200-clear-tagged-node-expiry",
				Migrate: func(tx *gorm.DB) error {
					err := tx.Exec(`
UPDATE nodes
SET expiry = NULL
WHERE tags IS NOT NULL AND tags != '[]' AND tags != '' AND tags != 'null'
	AND expiry IS NOT NULL;
						`).Error
					if err != nil {
						return fmt.Errorf("clearing expiry on tagged nodes: %w", err)
					}

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			{
				// Add oidc_groups column to users table to store OIDC group
				// memberships. Groups are fetched from the OIDC provider's
				// 'groups' claim during authentication and can be referenced
				// in ACL policies for group-based access control.
				// This enables parity with Tailscale's OIDC group support.
				ID: "202609291402-add-oidc-groups-to-users",
				Migrate: func(tx *gorm.DB) error {
					if !tx.Migrator().HasColumn(&types.User{}, "OIDCGroups") {
						err := tx.Migrator().AddColumn(&types.User{}, "OIDCGroups")
						if err != nil {
							return fmt.Errorf("adding oidc_groups to users: %w", err)
						}
					}

					return nil
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			// 0.30: unified credentials table (InitSchema keeps ensureCredentialsTable).
			// TODO(kradalby): remove in 0.31 with the credentials migration.
			{
				// Create the unified credentials table; the next migration
				// backfills it. Explicit DDL for both dialects (no AutoMigrate).
				ID:       "202609231200-create-credentials",
				Migrate:  ensureCredentialsTable,
				Rollback: func(db *gorm.DB) error { return nil },
			},
			{
				// Move every credential into the unified table and drop the
				// per-kind tables (see migrateToCredentials).
				ID: "202609231300-migrate-to-credentials",
				Migrate: func(tx *gorm.DB) error {
					// Already migrated (e.g. fresh DB via InitSchema): nothing to do.
					if !tx.Migrator().HasTable("pre_auth_keys") &&
						!tx.Migrator().HasTable("api_keys") {
						return nil
					}

					return tx.Transaction(migrateToCredentials)
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			{
				// Add webhooks table for monitoring/alerting integrations.
				// Webhooks can be configured to receive HTTP POST notifications
				// for events like node up/down, health check failures, and alerts.
				ID: "202609301000-create-webhooks-table",
				Migrate: func(tx *gorm.DB) error {
					return ensureWebhooksTable(tx)
				},
				Rollback: func(tx *gorm.DB) error {
					return tx.Migrator().DropTable(&types.Webhook{})
				},
			},
			{
				// Add headplane_users table for multi-user Headplane authentication.
				// This enables multiple users to log in to the Headplane UI with
				// individual credentials (username + password) and role-based access.
				// Migration auto-creates an admin user from the existing config password.
				ID: "202610031721-create-headplane-users",
				Migrate: func(tx *gorm.DB) error {
					return ensureHeadplaneUsersTable(tx, cfg)
				},
				Rollback: func(tx *gorm.DB) error {
					return tx.Migrator().DropTable(&HeadplaneUser{})
				},
			},
			{
				// Migrate headplane_settings from single-user to per-user.
				// Adds user_id column, removes single-row constraint, and migrates
				// existing settings to the admin user.
				ID: "202610041200-per-user-headplane-settings",
				Migrate: func(tx *gorm.DB) error {
					return migrateHeadplaneSettingsToPerUser(tx)
				},
				Rollback: func(tx *gorm.DB) error {
					return nil
				},
			},
			{
				// Restore Headplane tables that were removed to work around the
				// schema validator before they were added to schema.sql.
				ID: headplaneSchemaRepairMigrationID,
				Migrate: func(tx *gorm.DB) error {
					return repairHeadplaneSchema(tx, cfg)
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
			{
				// GORM created the Headplane indexes with quoted identifiers,
				// which does not match the canonical SQLite schema text.
				ID: headplaneIndexNormalizationMigrationID,
				Migrate: func(tx *gorm.DB) error {
					return normalizeHeadplaneIndexes(tx)
				},
				Rollback: func(db *gorm.DB) error { return nil },
			},
		},
	)

	migrations.InitSchema(func(tx *gorm.DB) error {
		// Credentials use the migration's explicit DDL (AutoMigrate cannot
		// express its CHECK constraints), created before Node so the
		// nodes.auth_key_id foreign key to credentials(id) can be created.
		err := tx.AutoMigrate(&types.User{})
		if err != nil {
			return err
		}

		err = ensureCredentialsTable(tx)
		if err != nil {
			return err
		}

		err = tx.AutoMigrate(&types.Node{}, &types.Policy{})
		if err != nil {
			return err
		}

		// Webhooks use explicit DDL so their indexes match schema.sql exactly;
		// AutoMigrate would emit backticked index DDL that fails validation.
		err = ensureWebhooksTable(tx)
		if err != nil {
			return err
		}

		err = ensureHeadplaneUsersTable(tx, cfg)
		if err != nil {
			return err
		}

		// Drop all indexes (both GORM-created and potentially pre-existing ones)
		// to ensure we can recreate them in the correct format
		dropIndexes := []string{
			`DROP INDEX IF EXISTS "idx_users_deleted_at"`,
			`DROP INDEX IF EXISTS "idx_policies_deleted_at"`,
			`DROP INDEX IF EXISTS "idx_provider_identifier"`,
			`DROP INDEX IF EXISTS "idx_name_provider_identifier"`,
			`DROP INDEX IF EXISTS "idx_name_no_provider_identifier"`,
			`DROP INDEX IF EXISTS "idx_nodes_auth_key_id"`,
		}

		for _, dropSQL := range dropIndexes {
			err := tx.Exec(dropSQL).Error
			if err != nil {
				return err
			}
		}

		// Recreate indexes without backticks to match schema.sql format
		indexes := []string{
			`CREATE INDEX idx_users_deleted_at ON users(deleted_at)`,
			`CREATE INDEX idx_policies_deleted_at ON policies(deleted_at)`,
			`CREATE UNIQUE INDEX idx_provider_identifier ON users(provider_identifier) WHERE provider_identifier IS NOT NULL`,
			`CREATE UNIQUE INDEX idx_name_provider_identifier ON users(name, provider_identifier)`,
			`CREATE UNIQUE INDEX idx_name_no_provider_identifier ON users(name) WHERE provider_identifier IS NULL`,
			`CREATE INDEX idx_nodes_auth_key_id ON nodes(auth_key_id)`,
		}

		for _, indexSQL := range indexes {
			err := tx.Exec(indexSQL).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	err = runMigrations(cfg.Database, dbConn, migrations)
	if err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	// Store the current version in the database after migrations succeed.
	// Dev builds skip this to preserve the stored version for the next
	// real versioned binary.
	currentVersion := types.GetVersionInfo().Version
	if !isDev(currentVersion) {
		err = setDatabaseVersion(dbConn, currentVersion)
		if err != nil {
			return nil, fmt.Errorf(
				"storing database version: %w",
				err,
			)
		}
	}

	// Validate that the schema ends up in the expected state.
	// This is currently only done on sqlite as squibble does not
	// support Postgres and we use our sqlite schema as our source of
	// truth.
	if cfg.Database.Type == types.DatabaseSqlite {
		sqlConn, err := dbConn.DB()
		if err != nil {
			return nil, fmt.Errorf("getting DB from gorm: %w", err)
		}

		// or else it blocks...
		sqlConn.SetMaxIdleConns(maxIdleConns)

		sqlConn.SetMaxOpenConns(maxOpenConns)
		defer sqlConn.SetMaxIdleConns(1)
		defer sqlConn.SetMaxOpenConns(1)

		ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
		defer cancel()

		opts := squibble.DigestOptions{
			IgnoreTables: []string{
				// Litestream tables, these are inserted by
				// litestream and not part of our schema
				// https://litestream.io/how-it-works
				"_litestream_lock",
				"_litestream_seq",
			},
		}

		if err := squibble.Validate(ctx, sqlConn, dbSchema, &opts); err != nil { //nolint:noinlineerr
			return nil, fmt.Errorf("validating schema: %w", err)
		}
	}

	db := HSDatabase{
		DB:  dbConn,
		cfg: cfg,
	}

	return &db, err
}

func openDB(cfg types.DatabaseConfig) (*gorm.DB, error) {
	// TODO(kradalby): Integrate this with zerolog
	var dbLogger logger.Interface
	if cfg.Debug {
		dbLogger = util.NewDBLogWrapper(&log.Logger, cfg.Gorm.SlowThreshold, cfg.Gorm.SkipErrRecordNotFound, cfg.Gorm.ParameterizedQueries)
	} else {
		dbLogger = logger.Default.LogMode(logger.Silent)
	}

	switch cfg.Type {
	case types.DatabaseSqlite:
		dir := filepath.Dir(cfg.Sqlite.Path)

		err := util.EnsureDir(dir)
		if err != nil {
			return nil, fmt.Errorf("creating directory for sqlite: %w", err)
		}

		log.Info().
			Str("database", types.DatabaseSqlite).
			Str("path", cfg.Sqlite.Path).
			Msg("Opening database")

		// Build SQLite configuration with pragmas set at connection time
		sqliteConfig := sqliteconfig.Default(cfg.Sqlite.Path)
		if cfg.Sqlite.WriteAheadLog {
			sqliteConfig.JournalMode = sqliteconfig.JournalModeWAL
			sqliteConfig.WALAutocheckpoint = cfg.Sqlite.WALAutoCheckPoint
		}

		connectionURL, err := sqliteConfig.ToURL()
		if err != nil {
			return nil, fmt.Errorf("building sqlite connection URL: %w", err)
		}

		db, err := gorm.Open(
			sqlite.Open(connectionURL),
			&gorm.Config{
				PrepareStmt: cfg.Gorm.PrepareStmt,
				Logger:      dbLogger,
			},
		)

		// The pure Go SQLite library does not handle locking in
		// the same way as the C based one and we can't use the gorm
		// connection pool as of 2022/02/23.
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(1)
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetConnMaxIdleTime(time.Hour)

		return db, err

	case types.DatabasePostgres:
		dbString := fmt.Sprintf(
			"host=%s dbname=%s user=%s",
			cfg.Postgres.Host,
			cfg.Postgres.Name,
			cfg.Postgres.User,
		)

		log.Info().
			Str("database", types.DatabasePostgres).
			Str("path", dbString).
			Msg("Opening database")

		if sslEnabled, err := strconv.ParseBool(cfg.Postgres.Ssl); err == nil { //nolint:noinlineerr
			if !sslEnabled {
				dbString += " sslmode=disable"
			}
		} else {
			dbString += " sslmode=" + cfg.Postgres.Ssl
		}

		if cfg.Postgres.Port != 0 {
			dbString += fmt.Sprintf(" port=%d", cfg.Postgres.Port)
		}

		if cfg.Postgres.Pass != "" {
			dbString += " password=" + cfg.Postgres.Pass
		}

		db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{
			Logger: dbLogger,
		})
		if err != nil {
			return nil, err
		}

		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(cfg.Postgres.MaxIdleConnections)
		sqlDB.SetMaxOpenConns(cfg.Postgres.MaxOpenConnections)
		sqlDB.SetConnMaxIdleTime(
			time.Duration(cfg.Postgres.ConnMaxIdleTimeSecs) * time.Second,
		)

		return db, nil
	}

	return nil, fmt.Errorf(
		"database of type %s is not supported: %w",
		cfg.Type,
		errDatabaseNotSupported,
	)
}

func runMigrations(cfg types.DatabaseConfig, dbConn *gorm.DB, migrations *gormigrate.Gormigrate) error {
	if cfg.Type == types.DatabaseSqlite {
		if err := migrations.Migrate(); err != nil { //nolint:noinlineerr
			return err
		}

		// Check for constraint violations at the end
		type constraintViolation struct {
			Table           string
			RowID           int
			Parent          string
			ConstraintIndex int
		}

		var violatedConstraints []constraintViolation

		rows, err := dbConn.Raw("PRAGMA foreign_key_check").Rows()
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var violation constraintViolation

			err := rows.Scan(&violation.Table, &violation.RowID, &violation.Parent, &violation.ConstraintIndex)
			if err != nil {
				return err
			}

			violatedConstraints = append(violatedConstraints, violation)
		}

		if err := rows.Err(); err != nil { //nolint:noinlineerr
			return err
		}

		if len(violatedConstraints) > 0 {
			for _, violation := range violatedConstraints {
				log.Error().
					Str("table", violation.Table).
					Int("row_id", violation.RowID).
					Str("parent", violation.Parent).
					Msg("Foreign key constraint violated")
			}

			return errForeignKeyConstraintsViolated
		}
	} else {
		// PostgreSQL can run all migrations in one block - no foreign key issues
		err := migrations.Migrate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (hsdb *HSDatabase) PingDB(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	sqlDB, err := hsdb.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}

func (hsdb *HSDatabase) Close() error {
	db, err := hsdb.DB.DB()
	if err != nil {
		return err
	}

	if hsdb.cfg.Database.Type == types.DatabaseSqlite && hsdb.cfg.Database.Sqlite.WriteAheadLog {
		db.Exec("VACUUM") //nolint:errcheck,noctx
	}

	return db.Close()
}

func (hsdb *HSDatabase) Read(fn func(rx *gorm.DB) error) error {
	rx := hsdb.DB.Begin()
	defer rx.Rollback()

	return fn(rx)
}

func Read[T any](db *gorm.DB, fn func(rx *gorm.DB) (T, error)) (T, error) {
	rx := db.Begin()
	defer rx.Rollback()

	ret, err := fn(rx)
	if err != nil {
		var no T
		return no, err
	}

	return ret, nil
}

func (hsdb *HSDatabase) Write(fn func(tx *gorm.DB) error) error {
	tx := hsdb.DB.Begin()
	defer tx.Rollback()

	err := fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit().Error
}

func Write[T any](db *gorm.DB, fn func(tx *gorm.DB) (T, error)) (T, error) {
	tx := db.Begin()
	defer tx.Rollback()

	ret, err := fn(tx)
	if err != nil {
		var no T
		return no, err
	}

	return ret, tx.Commit().Error
}

// ensureHeadplaneUsersTable creates the headplane_users table and migrates
// the existing single-user password to a default admin user.
func ensureHeadplaneUsersTable(tx *gorm.DB, cfg *types.Config) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		usersTableMissing := !tx.Migrator().HasTable(headplaneUsersTableName)
		if err := EnsureHeadplaneTables(tx); err != nil {
			return err
		}

		// Phase 13c: No longer create default admin user from config password.
		// Legacy tables are preserved for rollback, but runtime user creation
		// and password authentication are retired.
		_ = usersTableMissing

		return nil
	})
}

// migrateHeadplaneSettingsToPerUser migrates the headplane_settings table
// from single-user to per-user by adding user_id column and migrating
// existing settings to the admin user (id=1).
func migrateHeadplaneSettingsToPerUser(tx *gorm.DB) error {
	// Check if the table exists
	if !tx.Migrator().HasTable("headplane_settings") {
		// Table doesn't exist yet, nothing to migrate
		return nil
	}

	// Check if user_id column already exists (migration already run)
	if tx.Migrator().HasColumn("headplane_settings", "user_id") {
		return nil
	}

	// Dialect-specific migration
	dialect := tx.Name()

	if dialect == sqliteDialect {
		// SQLite: Need to recreate table since ALTER TABLE has limitations

		// 1. Read existing settings if any
		type OldSettings struct {
			ID              int    `gorm:"column:id"`
			APIKeyEncrypted string `gorm:"column:api_key_encrypted"`
			APIKeyNonce     string `gorm:"column:api_key_nonce"`
			APIKeySalt      string `gorm:"column:api_key_salt"`
			Theme           string `gorm:"column:theme"`
			ProfileName     string `gorm:"column:profile_name"`
		}

		var oldSettings OldSettings
		hasExisting := false
		err := tx.Table("headplane_settings").First(&oldSettings).Error
		if err == nil {
			hasExisting = true
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("reading existing settings: %w", err)
		}

		// 2. Drop old table
		if err := tx.Exec("DROP TABLE IF EXISTS headplane_settings").Error; err != nil {
			return fmt.Errorf("dropping old headplane_settings table: %w", err)
		}

		// 3. Create the canonical per-user settings table.
		if err := createHeadplaneSettingsTable(tx); err != nil {
			return err
		}

		// 4. Migrate existing settings to admin user (id=1) if any existed
		if hasExisting {
			// Get admin user ID (should be 1, but verify)
			var adminUser HeadplaneUser
			err := tx.Where("role = ?", "admin").Order("id ASC").First(&adminUser).Error
			if err != nil {
				log.Warn().Msg("No admin user found, skipping settings migration")
			} else {
				if err := tx.Exec(`
					INSERT INTO headplane_settings (user_id, api_key_encrypted, api_key_nonce, api_key_salt, theme, profile_name, updated_at)
					VALUES (?, ?, ?, ?, ?, ?, datetime('now'))
				`, adminUser.ID, oldSettings.APIKeyEncrypted, oldSettings.APIKeyNonce, oldSettings.APIKeySalt, oldSettings.Theme, oldSettings.ProfileName).Error; err != nil {
					return fmt.Errorf("migrating settings to admin user: %w", err)
				}
				log.Info().Uint("user_id", adminUser.ID).Msg("Migrated existing settings to admin user")
			}
		}
	} else {
		// PostgreSQL: Can use ALTER TABLE

		// 1. Add user_id column (nullable first)
		if err := tx.Exec("ALTER TABLE headplane_settings ADD COLUMN user_id INTEGER").Error; err != nil {
			return fmt.Errorf("adding user_id column: %w", err)
		}

		// 2. Get admin user ID
		var adminUser HeadplaneUser
		err := tx.Where("role = ?", "admin").Order("id ASC").First(&adminUser).Error
		if err != nil {
			log.Warn().Msg("No admin user found for settings migration")
		} else {
			// 3. Set user_id to admin for existing rows
			if err := tx.Exec("UPDATE headplane_settings SET user_id = ? WHERE user_id IS NULL", adminUser.ID).Error; err != nil {
				return fmt.Errorf("setting user_id for existing settings: %w", err)
			}
			log.Info().Uint("user_id", adminUser.ID).Msg("Migrated existing settings to admin user")
		}

		// 4. Make user_id NOT NULL and add unique constraint
		if err := tx.Exec("ALTER TABLE headplane_settings ALTER COLUMN user_id SET NOT NULL").Error; err != nil {
			return fmt.Errorf("making user_id NOT NULL: %w", err)
		}

		if err := tx.Exec("CREATE UNIQUE INDEX idx_headplane_settings_user_id ON headplane_settings(user_id)").Error; err != nil {
			return fmt.Errorf("adding unique index on user_id: %w", err)
		}

		// 5. Drop old id constraint if it exists (CHECK id = 1)
		// PostgreSQL doesn't have easy way to drop CHECK, and it won't cause issues
	}

	return nil
}
