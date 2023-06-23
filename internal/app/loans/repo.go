package loans

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
func NewRepo(dbConnections *db.Connections) ILoanRepo {
	return &repo{
		base.NewRepo(dbConnections),
	}
}

type ILoanRepo interface {
	Create(ctx context.Context, receiver base.IModel) error
	Update(ctx context.Context, receiver base.IModel, selectiveList ...string) error
	FindByID(ctx context.Context, receiver base.IModel, id int64) error
	Transaction(ctx context.Context, fc func(ctx context.Context) error) error
	List(ctx context.Context, receivers interface{}, offset, limit int) error
}

func (r repo) List(ctx context.Context, receivers interface{}, offset, limit int) error {
	tnt := tenant.From(ctx)
	q := r.Repo.GetConnection(ctx)
	if tnt.IsUser() {
		q = q.Where("user_id = ?", tnt.User().GetID())
	}
	q = q.
		Preload("LoanPayments").
		Offset(offset).
		Limit(limit).
		Order("created_at desc")

	q = q.Find(receivers)

	err := base.GetDBError(q)

	return err
}

func (r repo) FindByID(ctx context.Context, receiver base.IModel, id int64) error {
	q := r.GetConnection(ctx).Where("id = ?", id).Preload("LoanPayments")
	tnt := tenant.From(ctx)
	if tnt.IsUser() {
		q.Where("user_id = ?", tnt.User().GetID())
	}

	q = q.First(receiver)

	return base.GetDBError(q)
}
