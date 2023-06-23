package loans

import (
	"context"
)

// Service .
type Service struct {
	Core ILoanCore
}

// NewLoanService returns a service.
func NewLoanService(core ILoanCore) *Service {
	return &Service{
		Core: core,
	}
}

func (s *Service) Create(ctx context.Context, input interface{}) (interface{}, error) {
	response, err := s.Core.Create(ctx, input)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) FetchLoans(ctx context.Context, input interface{}) (interface{}, error) {
	response, err := s.Core.FetchLoans(ctx, input)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) ApproveLoan(ctx context.Context, input interface{}) (interface{}, error) {
	response, err := s.Core.ApproveLoan(ctx, input)

	if err != nil {
		return nil, err
	}

	return response, nil
}
