package fixtures

import (
	"context"

	users "github.com/roguexray007/loan-app/internal/app/users"
)

func UserCreate(entities ...*users.User) {
	for _, entity := range entities {
		err := users.GetRepo().Create(context.TODO(), entity)
		if err != nil {
			panic(err)
		}
	}
}
