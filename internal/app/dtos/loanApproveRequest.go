package dtos

import "github.com/gin-gonic/gin"

type LoanApproveRequest struct {
	LoanID int64 `json:"loan_id"`
}

func (lcr *LoanApproveRequest) Build(ctx *gin.Context) error {
	err := ctx.ShouldBindJSON(lcr)

	if err != nil {
		return err
	}

	return nil
}
