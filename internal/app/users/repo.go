package users

import (
	"context"

	"github.com/roguexray007/loan-app/internal/app/base"
	"github.com/roguexray007/loan-app/pkg/db"
)

type repo struct {
	base.Repo
}

var repoObj IUserRepo

// NewRepo creates new repo instance for the model
func NewRepo(dbConnections *db.Connections) IUserRepo {
	repoObj = &repo{
		base.NewRepo(dbConnections),
	}
	return repoObj
}

func GetRepo() IUserRepo {
	return repoObj
}

func (r repo) FindByUsername(ctx context.Context, receiver base.IModel, username string) error {
	q := r.GetConnection(ctx).Where("username = ?", username).First(receiver)

	return base.GetDBError(q)
}

type IUserRepo interface {
	Create(ctx context.Context, receiver base.IModel) error
	Update(ctx context.Context, receiver base.IModel, selectiveList ...string) error
	FindByID(ctx context.Context, receiver base.IModel, id int64) error
	FindByUsername(ctx context.Context, receiver base.IModel, username string) error
}
