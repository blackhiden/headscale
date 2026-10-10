package db

import (
	"fmt"
	"strings"

	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
	"gorm.io/gorm"
)

// headplaneUsersDDLSQLite and headplaneSettingsDDLSQLite match schema.sql.
// They are used for both fresh databases and migrations.
//
//nolint:gosec // DDL, not a credential
const headplaneUsersDDLSQLite = `CREATE TABLE headplane_users(
  id integer PRIMARY KEY AUTOINCREMENT,
  username text NOT NULL,
  password_hash text NOT NULL,
  role text DEFAULT "user",
  created_at datetime,
  updated_at datetime
)`

//nolint:gosec // DDL, not a credential
const headplaneSettingsDDLSQLite = `CREATE TABLE headplane_settings(
  id integer PRIMARY KEY AUTOINCREMENT,
  user_id integer NOT NULL,
  api_key_encrypted text,
  api_key_nonce text,
  api_key_salt text,
  theme text DEFAULT "light",
  profile_name text,
  updated_at datetime
)`

//nolint:gosec // DDL, not a credential
const headplaneUsersDDLPostgres = `CREATE TABLE headplane_users(
  id bigserial PRIMARY KEY,
  username text NOT NULL,
  password_hash text NOT NULL,
  role text DEFAULT 'user',
  created_at timestamptz,
  updated_at timestamptz
)`

//nolint:gosec // DDL, not a credential
const headplaneSettingsDDLPostgres = `CREATE TABLE headplane_settings(
  id bigserial PRIMARY KEY,
  user_id bigint NOT NULL,
  api_key_encrypted text,
  api_key_nonce text,
  api_key_salt text,
  theme text DEFAULT 'light',
  profile_name text,
  updated_at timestamptz
)`

const (
	headplaneUsersUsernameIndex  = "idx_headplane_users_username"
	headplaneSettingsUserIDIndex = "idx_headplane_settings_user_id"
	headplaneUsersTableName      = "headplane_users"
	headplaneSettingsTableName   = "headplane_settings"
)

var headplaneUserIndexes = []string{
	`CREATE UNIQUE INDEX idx_headplane_users_username ON headplane_users(username)`,
}

var headplaneSettingsIndexes = []string{
	`CREATE UNIQUE INDEX idx_headplane_settings_user_id ON headplane_settings(user_id)`,
}

type headplaneIndex struct {
	name string
	ddl  string
}

var headplaneIndexes = []headplaneIndex{
	{name: headplaneUsersUsernameIndex, ddl: headplaneUserIndexes[0]},
	{name: headplaneSettingsUserIDIndex, ddl: headplaneSettingsIndexes[0]},
}

// EnsureHeadplaneTables creates the Headplane tables and indexes when missing.
// Both new installations and upgrades use this so SQLite validation sees one
// canonical schema.
func EnsureHeadplaneTables(tx *gorm.DB) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		if !tx.Migrator().HasTable(headplaneUsersTableName) {
			if err := createHeadplaneUsersTable(tx); err != nil {
				return err
			}
		} else if !tx.Migrator().HasIndex(headplaneUsersTableName, headplaneUsersUsernameIndex) {
			if err := tx.Exec(headplaneUserIndexes[0]).Error; err != nil {
				return fmt.Errorf("creating headplane users index: %w", err)
			}
		}

		if !tx.Migrator().HasTable(headplaneSettingsTableName) {
			if err := createHeadplaneSettingsTable(tx); err != nil {
				return err
			}
		} else if tx.Migrator().HasColumn(headplaneSettingsTableName, "user_id") &&
			!tx.Migrator().HasIndex(headplaneSettingsTableName, headplaneSettingsUserIDIndex) {
			if err := tx.Exec(headplaneSettingsIndexes[0]).Error; err != nil {
				return fmt.Errorf("creating headplane settings index: %w", err)
			}
		}

		return nil
	})
}

func repairHeadplaneSchema(tx *gorm.DB, cfg *types.Config) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		if err := ensureHeadplaneUsersTable(tx, cfg); err != nil {
			return err
		}

		if tx.Name() != sqliteDialect {
			return nil
		}

		for _, table := range []string{headplaneUsersTableName, headplaneSettingsTableName} {
			if err := rebuildLegacyHeadplaneTable(tx, table); err != nil {
				return err
			}
		}

		return nil
	})
}

func normalizeHeadplaneIndexes(tx *gorm.DB) error {
	if tx.Name() != sqliteDialect {
		return nil
	}

	return tx.Transaction(func(tx *gorm.DB) error {
		for _, index := range headplaneIndexes {
			var current string
			if err := tx.Raw(
				`SELECT sql FROM sqlite_master WHERE type = 'index' AND name = ?`,
				index.name,
			).Scan(&current).Error; err != nil {
				return fmt.Errorf("reading %s definition: %w", index.name, err)
			}

			if current == index.ddl {
				continue
			}

			if err := tx.Exec(`DROP INDEX IF EXISTS ` + index.name).Error; err != nil {
				return fmt.Errorf("dropping %s: %w", index.name, err)
			}
			if err := tx.Exec(index.ddl).Error; err != nil {
				return fmt.Errorf("creating %s: %w", index.name, err)
			}
		}

		return nil
	})
}

func rebuildLegacyHeadplaneTable(tx *gorm.DB, table string) error {
	var ddl string
	if err := tx.Raw(`SELECT sql FROM sqlite_master WHERE type = 'table' AND name = ?`, table).
		Scan(&ddl).Error; err != nil {
		return fmt.Errorf("reading %s schema: %w", table, err)
	}

	if !strings.Contains(ddl, "DEFAULT '") {
		return nil
	}

	sequence, err := sqliteSequence(tx, table)
	if err != nil {
		return err
	}

	legacyTable := table + "_legacy"
	if err := tx.Exec(`ALTER TABLE ` + table + ` RENAME TO ` + legacyTable).Error; err != nil {
		return fmt.Errorf("renaming %s table: %w", table, err)
	}

	if table == headplaneUsersTableName {
		if err := tx.Exec(`DROP INDEX IF EXISTS ` + headplaneUsersUsernameIndex).Error; err != nil {
			return fmt.Errorf("dropping headplane users index: %w", err)
		}
		if err := createHeadplaneUsersTable(tx); err != nil {
			return err
		}
		if err := tx.Exec(`INSERT INTO headplane_users (id, username, password_hash, role, created_at, updated_at)
SELECT id, username, password_hash, role, created_at, updated_at FROM headplane_users_legacy`).Error; err != nil {
			return fmt.Errorf("copying headplane users: %w", err)
		}
	} else {
		if err := tx.Exec(`DROP INDEX IF EXISTS ` + headplaneSettingsUserIDIndex).Error; err != nil {
			return fmt.Errorf("dropping headplane settings index: %w", err)
		}
		if err := createHeadplaneSettingsTable(tx); err != nil {
			return err
		}
		if err := tx.Exec(`INSERT INTO headplane_settings
(id, user_id, api_key_encrypted, api_key_nonce, api_key_salt, theme, profile_name, updated_at)
SELECT id, user_id, api_key_encrypted, api_key_nonce, api_key_salt, theme, profile_name, updated_at
FROM headplane_settings_legacy`).Error; err != nil {
			return fmt.Errorf("copying headplane settings: %w", err)
		}
	}

	if err := tx.Exec(`DROP TABLE ` + legacyTable).Error; err != nil {
		return fmt.Errorf("dropping legacy %s table: %w", table, err)
	}

	return raiseSQLiteSequence(tx, table, sequence)
}

func createHeadplaneUsersTable(tx *gorm.DB) error {
	ddl := headplaneUsersDDLSQLite
	if tx.Name() != sqliteDialect {
		ddl = headplaneUsersDDLPostgres
	}

	if err := tx.Exec(ddl).Error; err != nil {
		return fmt.Errorf("creating headplane users table: %w", err)
	}

	for _, stmt := range headplaneUserIndexes {
		if err := tx.Exec(stmt).Error; err != nil {
			return fmt.Errorf("creating headplane users index: %w", err)
		}
	}

	return nil
}

func createHeadplaneSettingsTable(tx *gorm.DB) error {
	ddl := headplaneSettingsDDLSQLite
	if tx.Name() != sqliteDialect {
		ddl = headplaneSettingsDDLPostgres
	}

	if err := tx.Exec(ddl).Error; err != nil {
		return fmt.Errorf("creating headplane settings table: %w", err)
	}

	for _, stmt := range headplaneSettingsIndexes {
		if err := tx.Exec(stmt).Error; err != nil {
			return fmt.Errorf("creating headplane settings index: %w", err)
		}
	}

	return nil
}
