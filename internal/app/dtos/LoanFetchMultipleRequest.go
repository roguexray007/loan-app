package dtos

import "github.com/gin-gonic/gin"

type LoanFetchMultipleRequest struct {
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

func (lcr *LoanFetchMultipleRequest) Build(ctx *gin.Context) error {
	err := ctx.ShouldBindQuery(&lcr)

	if err != nil {
		return err
	}

	return nil
}
