package dtos

import "github.com/gin-gonic/gin"

// LoanPaymentCreateRequest is to capture path parameters and validate them
type LoanPaymentCreateRequest struct {
	Amount      int64 `json:"amount"`
	Terms       int   `json:"terms"`
	LoanID      int64 `json:"loan_id"`
	ScheduledAt int64 `json:"scheduled_at"`
	SequenceNo  int   `json:"sequence_no"`
}

func (lcr *LoanPaymentCreateRequest) Build(ctx *gin.Context) error {
	err := ctx.ShouldBindJSON(lcr)

	if err != nil {
		return err
	}

	return nil
}

type LoanPaymentRequest struct {
	Amount int64 `json:"amount"`
	Terms  int   `json:"terms"`
	LoanID int64 `json:"loan_id"`
}

func (lcr *LoanPaymentRequest) Build(ctx *gin.Context) error {
	err := ctx.ShouldBindJSON(lcr)

	if err != nil {
		return err
	}

	return nil
}

type LoanMarkAsPaidRequest struct {
	Amount     int64 `json:"amount"`
	SequenceNo int   `json:"sequence_no"`
	LoanID     int64 `json:"loan_id"`
}

func (lcr *LoanMarkAsPaidRequest) Build(ctx *gin.Context) error {
	err := ctx.ShouldBindJSON(lcr)

	if err != nil {
		return err
	}

	return nil
}
