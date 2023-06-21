package payments

import "github.com/roguexray007/loan-app/internal/app/base"

const (
	TableLoanPayment  = "loan_payments"
	EntityLoanPayment = "loan_payment"
)

type LoanPayment struct {
	base.Model
	LoanID      int64  `json:"loan_id"`
	UserID      int64  `json:"user_id"`
	Amount      int64  `json:"amount"`
	Status      string `json:"status"`
	ScheduledAt int64  `json:"scheduled_at"`
	SequenceNo  int    `json:"sequence_no"`
}

// TableName sets the insert table name
func (l *LoanPayment) TableName() string {
	return TableLoanPayment
}

// EntityName gives the display name of the entity
func (l *LoanPayment) EntityName() string {
	return EntityLoanPayment
}

func (l *LoanPayment) MarkPending() *LoanPayment {
	l.Status = Pending
	return l
}

type ILoanPayment interface {
	GetID() int64
}
