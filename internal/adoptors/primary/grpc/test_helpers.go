package grpc

import (
	"context"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

// MockOddsService is a mock implementation of the OddsService interface for testing
type MockOddsService struct {
	mock.Mock
}

func (m *MockOddsService) CreateOdds(ctx context.Context, request domain.CreateOddsRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockOddsService) ReadOdds(ctx context.Context, request domain.ReadOddsRequest) ([]domain.Odds, error) {
	args := m.Called(ctx, request)
	return args.Get(0).([]domain.Odds), args.Error(1)
}

func (m *MockOddsService) UpdateOdds(ctx context.Context, request domain.CreateOddsRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockOddsService) DeleteOdds(ctx context.Context, request domain.DeleteOddsRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}
