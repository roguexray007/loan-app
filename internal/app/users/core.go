package loans

import (
	"context"

	"github.com/roguexray007/loan-app/internal/app/dtos"
)

type Core struct {
	repo IUserRepo
}

func NewCore(repo IUserRepo) IUserCore {
	return Core{
		repo: repo,
	}
}

func (c Core) Create(ctx context.Context, input interface{}) (interface{}, error) {
	userCreateInput := input.(*dtos.UserCreateRequest)

	//TODO: store pass using a hash
	user := User{
		Username: userCreateInput.Username,
		Pass:     userCreateInput.Pass, // currently no storing hash or encrypting the pass
	}

	err := c.repo.Create(ctx, &user)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

type IUserCore interface {
	Create(ctx context.Context, input interface{}) (interface{}, error)
}
