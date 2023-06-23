package dtos

import "github.com/gin-gonic/gin"

type LoanPayRequest struct {
	LoanID int64 `json:"loan_id"`
	Amount int64 `json:"amount"`
}

func (lcr *LoanPayRequest) Build(ctx *gin.Context) error {
	err := ctx.ShouldBindJSON(lcr)

	if err != nil {
		return err
	}

	return nil
}
