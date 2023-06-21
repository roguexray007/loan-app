package loans

import (
	"context"

	"github.com/roguexray007/loan-app/internal/app/base"
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
}
