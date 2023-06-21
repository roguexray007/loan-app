package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/app/common/enum"
	"github.com/roguexray007/loan-app/internal/app/dtos"
	users "github.com/roguexray007/loan-app/internal/app/users"
)

type UserV1 struct {
	userService *users.Service
}

var UserService UserV1

func NewUserController(userService *users.Service) {
	UserService = UserV1{
		userService: userService,
	}
}

func (controller *UserV1) CreateUser(ctx *gin.Context) (interface{}, error, int) {
	createUser := dtos.GetRequestBuilder(enum.UserCreateRequest)
	err := createUser.Build(ctx)

	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	response, ierr := controller.userService.Create(ctx.Request.Context(), createUser)

	if ierr != nil {
		return nil, ierr, http.StatusBadRequest
	}

	return response, nil, http.StatusOK
}
