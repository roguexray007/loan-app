package loans

import (
	"context"
)

// Service .
type Service struct {
	Core IUserCore
}

// NewUserService returns a service.
func NewUserService(core IUserCore) *Service {
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
