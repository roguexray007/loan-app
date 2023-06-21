package dtos

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/app/base"
	"github.com/roguexray007/loan-app/internal/app/common/enum"
)

type IResponseBuilder interface {
	Build(ctx context.Context, entity base.IEntity)
}

type IRequestBuilder interface {
	Build(ctx *gin.Context) error
}

func GetResponseBuilder(builderType enum.BuilderType) IResponseBuilder {
	switch builderType {

	default:
		return nil
	}
}

func GetRequestBuilder(builderType enum.BuilderType) IRequestBuilder {
	switch builderType {

	case enum.LoanCreateRequest:
		return &LoanCreateRequest{}
	case enum.UserCreateRequest:
		return &UserCreateRequest{}
	case enum.LoanPaymentCreateRequest:
		return &LoanPaymentCreateRequest{}
	case enum.LoanPaymentRequest:
		return &LoanPaymentRequest{}
	default:
		return nil
	}
}
