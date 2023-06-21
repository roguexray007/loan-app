package loans

import (
	"context"

	"github.com/roguexray007/loan-app/internal/app/dtos"
	"github.com/roguexray007/loan-app/internal/app/loans/payments"
	"github.com/roguexray007/loan-app/internal/routing/tenant"
)

type Core struct {
	repo ILoanRepo
}

func NewCore(repo ILoanRepo) ILoanCore {
	return Core{
		repo: repo,
	}
}

func (c Core) Create(ctx context.Context, input interface{}) (interface{}, error) {
	loanCreateInput := input.(*dtos.LoanCreateRequest)

	tnt := tenant.From(ctx)

	loan := Loan{
		UserID:    tnt.User().GetID(),
		Amount:    loanCreateInput.Amount,
		Terms:     loanCreateInput.Terms,
		TermsPaid: 0,
	}
	loan.MarkPending()

	err := c.repo.Transaction(ctx, func(ctx context.Context) error {
		err := c.repo.Create(ctx, &loan)
		if err != nil {
			return err
		}

		loanPayments, err := payments.GetCore().CreateScheduledPaymentsForLoan(ctx, &dtos.LoanPaymentRequest{
			Amount: loan.Amount,
			Terms:  loan.Terms,
			LoanID: loan.ID,
		})
		if err != nil {
			return err
		}
		loan.LoanPayments = loanPayments

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &loan, nil
}

type ILoanCore interface {
	Create(ctx context.Context, input interface{}) (interface{}, error)
}
