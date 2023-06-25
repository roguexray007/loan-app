package loans

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/roguexray007/loan-app/internal/app/base"
	"github.com/roguexray007/loan-app/internal/app/dtos"
	"github.com/roguexray007/loan-app/internal/app/loans/payments"
	"github.com/roguexray007/loan-app/internal/provider"
	"github.com/roguexray007/loan-app/internal/routing/tenant"
)

type Core struct {
	repo ILoanRepo
	base.Core
}

func NewCore(repo ILoanRepo) ILoanCore {
	return Core{
		repo: repo,
		Core: base.Core{
			provider.GetMutex(nil),
		},
	}
}

func (c Core) Create(ctx context.Context, input interface{}) (*Loan, error) {
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

func (c Core) FetchLoans(ctx context.Context, input interface{}) ([]Loan, error) {
	var Loans []Loan

	loanFetchInput := input.(*dtos.LoanFetchMultipleRequest)

	limit := 0
	if loanFetchInput.Limit == 0 {
		limit = 10
	}

	err := c.repo.List(ctx, &Loans, loanFetchInput.Offset, limit)
	if err != nil {
		return nil, err
	}

	return Loans, nil
}

func (c Core) ApproveLoan(ctx context.Context, input interface{}) (*Loan, error) {
	loanApproveInput := input.(*dtos.LoanApproveRequest)

	var loan Loan
	// take lock before approving loan for handling concurrent req on same loanID
	mu, err := c.AcquireResource(ctx, strconv.FormatInt(loanApproveInput.LoanID, 10), 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer c.ReleaseResource(ctx, mu)

	err = c.repo.FindByID(ctx, &loan, loanApproveInput.LoanID)
	if err != nil {
		return nil, err
	}

	// if loan is already approved , return
	if loan.IsApproved() {
		return &loan, nil
	}

	if loan.IsPaid() {
		return nil, fmt.Errorf("can't move loan from paid to approve state")
	}

	(&loan).MarkApproved()
	err = c.repo.Update(ctx, &loan, "status")
	if err != nil {
		return nil, err
	}

	return &loan, nil
}

func (c Core) PayLoan(ctx context.Context, input interface{}) (*Loan, error) {
	loanPayInput := input.(*dtos.LoanPayRequest)

	var loan Loan
	// take lock before approving loan for handling concurrent req on same loanID
	mu, err := c.AcquireResource(ctx, strconv.FormatInt(loanPayInput.LoanID, 10), 30*time.Second)
	if err != nil {
		return nil, err
	}
	defer c.ReleaseResource(ctx, mu)

	err = c.repo.Transaction(ctx, func(ctx context.Context) error {
		err = c.repo.FindByID(ctx, &loan, loanPayInput.LoanID)
		if err != nil {
			return err
		}

		if loan.IsPending() {
			return fmt.Errorf("loan not approved by admin")
		}

		if loan.IsPaid() {
			return fmt.Errorf("loan already paid")
		}

		var loanPayments []*payments.LoanPayment
		loanPayment, err := payments.GetCore().MarkAsPaid(ctx, &dtos.LoanMarkAsPaidRequest{
			Amount:     loanPayInput.Amount,
			SequenceNo: loan.TermsPaid + 1,
			LoanID:     loan.ID,
		})
		for _, lp := range loan.LoanPayments {
			if lp.SequenceNo == loanPayment.SequenceNo {
				loanPayments = append(loanPayments, loanPayment)
			} else {
				loanPayments = append(loanPayments, lp)
			}
		}
		loan.LoanPayments = loanPayments

		loan.TermsPaid = loan.TermsPaid + 1
		if loan.TermsPaid == loan.Terms {
			(&loan).MarkPaid()
		}

		err = c.repo.Update(ctx, &loan, "status", "terms_paid")
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &loan, nil
}

type ILoanCore interface {
	Create(ctx context.Context, input interface{}) (*Loan, error)
	FetchLoans(ctx context.Context, input interface{}) ([]Loan, error)
	ApproveLoan(ctx context.Context, input interface{}) (*Loan, error)
	PayLoan(ctx context.Context, input interface{}) (*Loan, error)
}
