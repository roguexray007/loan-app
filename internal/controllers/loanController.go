package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/app/common/enum"
	"github.com/roguexray007/loan-app/internal/app/dtos"
	"github.com/roguexray007/loan-app/internal/app/loans"
)

type LoanV1 struct {
	loanService *loans.Service
}

var LoanService LoanV1

func NewLoanController(loanService *loans.Service) {
	LoanService = LoanV1{
		loanService: loanService,
	}
}

func (controller *LoanV1) CreateLoan(ctx *gin.Context) (interface{}, error, int) {
	createLoan := dtos.GetRequestBuilder(enum.LoanCreateRequest)
	err := createLoan.Build(ctx)

	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	response, ierr := controller.loanService.Create(ctx.Request.Context(), createLoan)

	if ierr != nil {
		return nil, ierr, http.StatusBadRequest
	}

	return response, nil, http.StatusOK
}
