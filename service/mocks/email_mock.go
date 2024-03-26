package mocks

import (
	"context"
	"webapp/types"
)

type MockEmailService struct{}

func NewMockEmailService() *MockEmailService {
	return &MockEmailService{}
}

func (s *MockEmailService) SendVerificationEmailToQueue(ctx context.Context, user types.User) error {
	// Do nothing since it is a mock method
	return nil
}
