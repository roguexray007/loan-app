package base

import (
	goErr "errors"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

const (
	PQCodeUniqueViolation = "unique_violation"

	errDBError                   = "db_error"
	errNoRowAffected             = "no_row_affected"
	errRecordNotFound            = "record_not_found"
	errUniqueConstraintViolation = "unique_constraint_violation"
)

// GetDBError accepts db instance and the details
// creates appropriate error based on the type of query result
// if there is no error then returns nil
func GetDBError(d *gorm.DB) error {
	if d.Error == nil {
		return nil
	}

	// check of error is specific to dialect
	if de, ok := DialectError(d); ok {
		// is the specific error is captured then return it
		// else try construct further errors
		if err := de.ConstructError(); err != nil {
			return err
		}
	}

	// Construct error based on type of db operation
	err := func() error {
		switch true {
		case d.RecordNotFound():
			return goErr.New(errRecordNotFound + ": " + d.Error.Error())

		default:
			return goErr.New(errDBError + ": " + d.Error.Error())
		}
	}()

	// add specific details of error
	return err
}

// DialectError returns true if the error is from dialect
func DialectError(d *gorm.DB) (IDialectError, bool) {
	switch d.Error.(type) {
	case *pq.Error:
		return pqError{d.Error.(*pq.Error)}, true
	default:
		return nil, false
	}
}

// IDialectError interface to handler dialect related errors
type IDialectError interface {
	ConstructError() error
}

// pqError holds the error occurred by postgres
type pqError struct {
	err *pq.Error
}

// ConstructError will create appropriate error based on dialect
func (pqe pqError) ConstructError() error {
	switch pqe.err.Code.Name() {
	case PQCodeUniqueViolation:
		return goErr.
			New(errUniqueConstraintViolation)
	default:
		return nil
	}
}
