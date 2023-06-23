package payments

import (
	"context"

	"github.com/roguexray007/loan-app/internal/app/base"
	"github.com/roguexray007/loan-app/internal/routing/tenant"
	"github.com/roguexray007/loan-app/pkg/db"
)

type repo struct {
	base.Repo
}

// NewRepo creates new repo instance for the model
func NewRepo(dbConnections *db.Connections) ILoanPaymentRepo {
	return &repo{
		base.NewRepo(dbConnections),
	}
}

type ILoanPaymentRepo interface {
	Create(ctx context.Context, receiver base.IModel) error
	Update(ctx context.Context, receiver base.IModel, selectiveList ...string) error
	FindByID(ctx context.Context, receiver base.IModel, id int64) error
	FindByLoanIDAndSeqNo(ctx context.Context, receiver base.IModel, id int64, sequenceNo int) error
}

func (r repo) FindByLoanIDAndSeqNo(ctx context.Context, receiver base.IModel, loanID int64, sequenceNo int) error {
	q := r.GetConnection(ctx).Where("loan_id = ?", loanID).
		Where("sequence_no = ?", sequenceNo)

	tnt := tenant.From(ctx)
	if tnt.IsUser() {
		q.Where("user_id = ?", tnt.User().GetID())
	}

	q = q.First(receiver)

	return base.GetDBError(q)
}
