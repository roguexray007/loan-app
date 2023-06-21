package loans

import (
	"github.com/roguexray007/loan-app/internal/app/base"
	"github.com/roguexray007/loan-app/internal/app/loans/payments"
)

const (
	TableLoan  = "loans"
	EntityLoan = "loan"
)

type Loan struct {
	base.Model
	UserID    int64  `json:"user_id"`
	Amount    int64  `json:"amount"`
	Terms     int    `json:"terms"`
	TermsPaid int    `json:"terms_paid"`
	Status    string `json:"status"`

	LoanPayments []*payments.LoanPayment `json:"loanPayments" gorm:"foreignKey:loan_id"`
}

// TableName sets the insert table name
func (l *Loan) TableName() string {
	return TableLoan
}

// EntityName gives the display name of the entity
func (l *Loan) EntityName() string {
	return EntityLoan
}

func (l *Loan) MarkPending() *Loan {
	l.Status = Pending
	return l
}

type ILoan interface {
	GetID() int64
}
