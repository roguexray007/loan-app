package enum

type BuilderType int

const (
	// Requests
	LoanCreateRequest BuilderType = 1 + iota
	UserCreateRequest
	LoanPaymentCreateRequest
	LoanPaymentRequest
	LoanFetchMultipleRequest
	LoanApproveRequest
	LoanPayRequest

	// Responses
)
