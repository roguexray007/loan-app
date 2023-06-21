package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upLoanPayments, downLoanPayments)
}

func upLoanPayments(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE loan_payments (
		id BIGINT unsigned NOT NULL PRIMARY KEY,
		loan_id BIGINT unsigned DEFAULT NULL,
		user_id BIGINT unsigned DEFAULT NULL,
		amount BIGINT NOT NULL,
		sequence_no int DEFAULT NULL,
		status varchar(40) DEFAULT NULL,
		scheduled_at int NOT NULL,
		created_at int NOT NULL,
		updated_at int NOT NULL
	);`)

	return err
}

func downLoanPayments(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE IF EXISTS loan_payments`)
	return err
}
