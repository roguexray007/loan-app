package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upLoan, downLoan)
}

func upLoan(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE loans (
		id BIGINT unsigned NOT NULL PRIMARY KEY,
		user_id BIGINT unsigned DEFAULT NULL,
		amount BIGINT NOT NULL,
		terms int DEFAULT NULL,
		status varchar(40) DEFAULT NULL,
		created_at int NOT NULL,
		updated_at int NOT NULL
	);`)

	return err
}

func downLoan(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS loans`)
	return err
}
