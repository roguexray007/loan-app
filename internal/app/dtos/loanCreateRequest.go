package dtos

import "github.com/gin-gonic/gin"

// LoanCreateRequest is to capture path parameters and validate them
type LoanCreateRequest struct {
	Amount int64 `json:"amount"`
	Terms  int   `json:"terms"`
}

func (lcr *LoanCreateRequest) Build(ctx *gin.Context) error {
	err := ctx.ShouldBindJSON(lcr)

	if err != nil {
		return err
	}

	return nil
}
