package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upUsers, downUsers)
}

func upUsers(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE users (
		id BIGINT unsigned NOT NULL PRIMARY KEY,
		username text DEFAULT NULL,
		pass  text DEFAULT NULL,
		created_at int NOT NULL,
		updated_at int NOT NULL
	);`)

	return err
}

func downUsers(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS users`)
	return err
}
