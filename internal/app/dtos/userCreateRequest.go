package dtos

import "github.com/gin-gonic/gin"

// UserCreateRequest is to capture path parameters and validate them
type UserCreateRequest struct {
	Username string `json:"username"`
	Pass     string `json:"pass"`
}

func (ucr *UserCreateRequest) Build(ctx *gin.Context) error {
	err := ctx.ShouldBindJSON(ucr)

	if err != nil {
		return err
	}

	return nil
}
