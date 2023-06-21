package payments

import (
	"context"
	"math"
	"time"

	"github.com/roguexray007/loan-app/internal/app/dtos"
	"github.com/roguexray007/loan-app/internal/routing/tenant"
)

type Core struct {
	repo ILoanPaymentRepo
}

var loanPaymentCore ILoanPaymentCore

func NewCore(repo ILoanPaymentRepo) ILoanPaymentCore {
	loanPaymentCore = &Core{
		repo: repo,
	}

	return loanPaymentCore
}

func GetCore() ILoanPaymentCore {
	return loanPaymentCore
}

func (c *Core) Create(ctx context.Context, input interface{}) (*LoanPayment, error) {
	loanPaymentCreateInput := input.(*dtos.LoanPaymentCreateRequest)

	tnt := tenant.From(ctx)

	loanPayment := LoanPayment{
		LoanID:      loanPaymentCreateInput.LoanID,
		UserID:      tnt.User().GetID(),
		Amount:      loanPaymentCreateInput.Amount,
		ScheduledAt: loanPaymentCreateInput.ScheduledAt,
		SequenceNo:  loanPaymentCreateInput.SequenceNo,
	}
	loanPayment.MarkPending()

	err := c.repo.Create(ctx, &loanPayment)

	if err != nil {
		return nil, err
	}

	return &loanPayment, nil
}

func (c *Core) CreateScheduledPaymentsForLoan(ctx context.Context, input interface{}) ([]*LoanPayment, error) {
	loanPaymentCreateInput := input.(*dtos.LoanPaymentRequest)
	amount := int64(math.Round(float64(loanPaymentCreateInput.Amount) / float64(loanPaymentCreateInput.Terms)))
	var loanPayments []*LoanPayment

	for seq := 1; seq <= loanPaymentCreateInput.Terms; seq++ {
		loanPaymentReq := &dtos.LoanPaymentCreateRequest{
			LoanID:      loanPaymentCreateInput.LoanID,
			Amount:      amount,
			ScheduledAt: time.Now().AddDate(0, 0, 7*seq).Unix(),
			SequenceNo:  seq,
		}
		loanPayment, err := c.Create(ctx, loanPaymentReq)
		if err != nil {
			return nil, err
		}

		loanPayments = append(loanPayments, loanPayment)
	}
	return loanPayments, nil
}

type ILoanPaymentCore interface {
	Create(ctx context.Context, input interface{}) (*LoanPayment, error)
	CreateScheduledPaymentsForLoan(ctx context.Context, input interface{}) ([]*LoanPayment, error)
}
