package service

import (
	"context"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
	"github.com/stretchr/testify/mock"
)

type testDependencies struct {
	repo      *MockOddsRepository
	validator *MockOddsValidator
	service   ports.OddsService
}

// MockOddsRepository is a mock implementation of the OddsRepository interface
type MockOddsRepository struct {
	mock.Mock
}

func (m *MockOddsRepository) Create(ctx context.Context, odds *domain.Odds) error {
	args := m.Called(ctx, odds)
	return args.Error(0)
}

func (m *MockOddsRepository) Read(ctx context.Context, league string, date time.Time) ([]domain.Odds, error) {
	args := m.Called(ctx, league, date)
	return args.Get(0).([]domain.Odds), args.Error(1)
}

func (m *MockOddsRepository) Update(ctx context.Context, odds *domain.Odds) error {
	args := m.Called(ctx, odds)
	return args.Error(0)
}

func (m *MockOddsRepository) Delete(ctx context.Context, league, homeTeam, awayTeam string, gameDate time.Time) error {
	args := m.Called(ctx, league, homeTeam, awayTeam, gameDate)
	return args.Error(0)
}

// MockOddsValidator is a mock implementation of the OddsValidator interface
type MockOddsValidator struct {
	mock.Mock
}

func (m *MockOddsValidator) ValidateCreateRequest(request domain.CreateOddsRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockOddsValidator) ValidateReadRequest(request domain.ReadOddsRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockOddsValidator) ValidateUpdateRequest(request domain.CreateOddsRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockOddsValidator) ValidateDeleteRequest(request domain.DeleteOddsRequest) error {
	args := m.Called(request)
	return args.Error(0)
}
