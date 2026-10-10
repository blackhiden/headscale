package db

import (
	"fmt"

	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
	"gorm.io/gorm"
)

// webhooksDDLSQLite matches schema.sql byte-for-byte (the squibble digest is the
// SQLite source of truth). The same statement is used by the migration and by
// InitSchema, so a fresh database and an upgraded one land in the same shape.
//
//nolint:gosec // DDL, not a credential
const webhooksDDLSQLite = `CREATE TABLE webhooks(
  id integer PRIMARY KEY AUTOINCREMENT,
  created_at datetime,
  updated_at datetime,
  deleted_at datetime,
  name text NOT NULL,
  url text NOT NULL,
  events text NOT NULL,
  headers text,
  secret text,
  enabled numeric NOT NULL DEFAULT true,
  timeout_seconds integer NOT NULL DEFAULT 10
)`

// webhooksDDLPostgres is the Postgres form of [webhooksDDLSQLite].
//
//nolint:gosec // DDL, not a credential
const webhooksDDLPostgres = `CREATE TABLE webhooks(
  id bigserial PRIMARY KEY,
  created_at timestamptz,
  updated_at timestamptz,
  deleted_at timestamptz,
  name text NOT NULL,
  url text NOT NULL,
  events text NOT NULL,
  headers text,
  secret text,
  enabled boolean NOT NULL DEFAULT true,
  timeout_seconds integer NOT NULL DEFAULT 10
)`

// webhookIndexes are created with the table.
var webhookIndexes = []string{
	`CREATE UNIQUE INDEX idx_webhooks_name ON webhooks(name)`,
	`CREATE INDEX idx_webhooks_deleted_at ON webhooks(deleted_at)`,
}

// ensureWebhooksTable creates the webhooks table and its indexes in one
// transaction and is a no-op once the table exists. Both InitSchema and the
// create-webhooks-table migration call it, so the table's DDL is identical on a
// fresh and an upgraded database — AutoMigrate emits backticked index DDL that
// would not match schema.sql.
func ensureWebhooksTable(tx *gorm.DB) error {
	if tx.Migrator().HasTable(&types.Webhook{}) {
		return nil
	}

	return tx.Transaction(createWebhooksTable)
}

// createWebhooksTable creates the webhooks table and its indexes.
func createWebhooksTable(tx *gorm.DB) error {
	ddl := webhooksDDLSQLite
	if tx.Name() != sqliteDialect {
		ddl = webhooksDDLPostgres
	}

	err := tx.Exec(ddl).Error
	if err != nil {
		return fmt.Errorf("creating webhooks table: %w", err)
	}

	for _, stmt := range webhookIndexes {
		err := tx.Exec(stmt).Error
		if err != nil {
			return fmt.Errorf("creating webhooks index: %w", err)
		}
	}

	return nil
}
