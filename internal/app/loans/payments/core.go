package payments

import (
	"context"
	"fmt"
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

func (c Core) MarkAsPaid(ctx context.Context, input interface{}) (*LoanPayment, error) {
	loanMarkAsPaidInput := input.(*dtos.LoanMarkAsPaidRequest)

	var loanPayment LoanPayment
	err := c.repo.FindByLoanIDAndSeqNo(ctx, &loanPayment, loanMarkAsPaidInput.LoanID, loanMarkAsPaidInput.SequenceNo)
	if err != nil {
		return nil, err
	}

	if loanPayment.IsPaid() {
		return nil, fmt.Errorf("loan payment already paid for seq %s", loanPayment.SequenceNo)
	}

	if loanPayment.Amount != loanMarkAsPaidInput.Amount {
		return nil, fmt.Errorf("incorrect loan payment amount")
	}

	(&loanPayment).MarkPaid()

	err = c.repo.Update(ctx, &loanPayment)
	if err != nil {
		return nil, err
	}

	return &loanPayment, nil
}

type ILoanPaymentCore interface {
	Create(ctx context.Context, input interface{}) (*LoanPayment, error)
	CreateScheduledPaymentsForLoan(ctx context.Context, input interface{}) ([]*LoanPayment, error)
	MarkAsPaid(ctx context.Context, input interface{}) (*LoanPayment, error)
}
