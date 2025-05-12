package service

import (
	"context"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestCreateOdds(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockOddsRepository)
	mockValidator := new(MockOddsValidator)

	// Create the service with mocks
	service := NewOddsService(mockRepo, mockValidator)

	// Test context
	ctx := context.Background()

	// Setup test data
	validRequest := domain.CreateOddsRequest{
		League:          "English Premier League",
		HomeTeam:        "Manchester United",
		AwayTeam:        "Liverpool",
		HomeTeamWinOdds: 2.5,
		AwayTeamWinOdds: 2.1,
		DrawOdds:        3.0,
		GameDate:        time.Now().Add(24 * time.Hour),
	}

	tests := []struct {
		name           string
		request        domain.CreateOddsRequest
		setupMocks     func()
		expectedResult error
	}{
		{
			name:    "Successful creation",
			request: validRequest,
			setupMocks: func() {
				// Validator validates successfully
				mockValidator.On("ValidateCreateRequest", validRequest).Return(nil)

				// Repository creates successfully
				// Note: We're not checking exact Odds struct as it has time.Now() fields which are hard to predict
				mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Odds")).Return(nil)
			},
			expectedResult: nil,
		},
		{
			name:    "Validation error",
			request: validRequest,
			setupMocks: func() {
				// Validator returns error
				mockValidator.On("ValidateCreateRequest", validRequest).Return(domain.ErrInvalidLeague)
				
				// Repository should not be called
			},
			expectedResult: domain.ErrInvalidLeague,
		},
		{
			name:    "Repository error",
			request: validRequest,
			setupMocks: func() {
				// Validator validates successfully
				mockValidator.On("ValidateCreateRequest", validRequest).Return(nil)

				// Repository returns error
				mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Odds")).Return(domain.ErrInvalidConfiguration)
			},
			expectedResult: domain.ErrInvalidConfiguration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock expectations and calls
			mockRepo = new(MockOddsRepository)
			mockValidator = new(MockOddsValidator)
			service = NewOddsService(mockRepo, mockValidator)

			// Setup mock expectations
			tt.setupMocks()

			// Call the method being tested
			result := service.CreateOdds(ctx, tt.request)

			// Assert result
			assert.Equal(t, tt.expectedResult, result)

			// Assert all expectations were met (all expected methods were called)
			mockRepo.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
		})
	}
}

func TestReadOdds(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockOddsRepository)
	mockValidator := new(MockOddsValidator)

	// Create the service with mocks
	service := NewOddsService(mockRepo, mockValidator)

	// Test context
	ctx := context.Background()

	// Setup test data
	validRequest := domain.ReadOddsRequest{
		League: "English Premier League",
		Date:   time.Now(),
	}

	expectedOdds := []domain.Odds{
		{
			ID:              1,
			League:          "English Premier League",
			HomeTeam:        "Manchester United",
			AwayTeam:        "Liverpool",
			HomeTeamWinOdds: 2.5,
			AwayTeamWinOdds: 2.1,
			DrawOdds:        3.0,
			GameDate:        time.Now().Add(24 * time.Hour),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	tests := []struct {
		name           string
		request        domain.ReadOddsRequest
		setupMocks     func()
		expectedResult []domain.Odds
		expectedError  error
	}{
		{
			name:    "Successful read",
			request: validRequest,
			setupMocks: func() {
				// Validator validates successfully
				mockValidator.On("ValidateReadRequest", validRequest).Return(nil)

				// Repository reads successfully
				mockRepo.On("Read", ctx, validRequest.League, validRequest.Date).Return(expectedOdds, nil)
			},
			expectedResult: expectedOdds,
			expectedError:  nil,
		},
		{
			name:    "Validation error",
			request: validRequest,
			setupMocks: func() {
				// Validator returns error
				mockValidator.On("ValidateReadRequest", validRequest).Return(domain.ErrInvalidLeague)
				
				// Repository should not be called
			},
			expectedResult: nil,
			expectedError:  domain.ErrInvalidLeague,
		},
		{
			name:    "Repository error",
			request: validRequest,
			setupMocks: func() {
				// Validator validates successfully
				mockValidator.On("ValidateReadRequest", validRequest).Return(nil)

				// Repository returns error
				mockRepo.On("Read", ctx, validRequest.League, validRequest.Date).Return([]domain.Odds{}, domain.ErrOddsNotFound)
			},
			expectedResult: nil,
			expectedError:  domain.ErrOddsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo = new(MockOddsRepository)
			mockValidator = new(MockOddsValidator)
			service = NewOddsService(mockRepo, mockValidator)

			// Setup mock expectations
			tt.setupMocks()

			// Call the method being tested
			result, err := service.ReadOdds(ctx, tt.request)

			// Assert result
			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError == nil {
				assert.Equal(t, tt.expectedResult, result)
			}

			// Assert all expectations were met (all expected methods were called)
			mockRepo.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
		})
	}
}

func TestUpdateOdds(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockOddsRepository)
	mockValidator := new(MockOddsValidator)

	// Create the service with mocks
	service := NewOddsService(mockRepo, mockValidator)

	// Test context
	ctx := context.Background()

	// Setup test data
	validRequest := domain.CreateOddsRequest{
		League:          "English Premier League",
		HomeTeam:        "Manchester United",
		AwayTeam:        "Liverpool",
		HomeTeamWinOdds: 2.5,
		AwayTeamWinOdds: 2.1,
		DrawOdds:        3.0,
		GameDate:        time.Now().Add(24 * time.Hour),
	}

	tests := []struct {
		name           string
		request        domain.CreateOddsRequest
		setupMocks     func()
		expectedResult error
	}{
		{
			name:    "Successful update",
			request: validRequest,
			setupMocks: func() {
				// Validator validates successfully
				mockValidator.On("ValidateUpdateRequest", validRequest).Return(nil)

				// Repository updates successfully
				mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Odds")).Return(nil)
			},
			expectedResult: nil,
		},
		{
			name:    "Validation error",
			request: validRequest,
			setupMocks: func() {
				// Validator returns error
				mockValidator.On("ValidateUpdateRequest", validRequest).Return(domain.ErrInvalidLeague)
				
				// Repository should not be called
			},
			expectedResult: domain.ErrInvalidLeague,
		},
		{
			name:    "Repository error",
			request: validRequest,
			setupMocks: func() {
				// Validator validates successfully
				mockValidator.On("ValidateUpdateRequest", validRequest).Return(nil)

				// Repository returns error
				mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Odds")).Return(domain.ErrOddsNotFound)
			},
			expectedResult: domain.ErrOddsNotFound,
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock expectations and calls
			mockRepo = new(MockOddsRepository)
			mockValidator = new(MockOddsValidator)
			service = NewOddsService(mockRepo, mockValidator)

			// Setup mock expectations
			tt.setupMocks()

			// Call the method being tested
			result := service.UpdateOdds(ctx, tt.request)

			// Assert result
			assert.Equal(t, tt.expectedResult, result)

			// Assert all expectations were met (all expected methods were called)
			mockRepo.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
		})
	}
}

func TestDeleteOdds(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockOddsRepository)
	mockValidator := new(MockOddsValidator)

	// Create the service with mocks
	service := NewOddsService(mockRepo, mockValidator)

	// Test context
	ctx := context.Background()

	// Setup test data
	gameDate := time.Now().Add(24 * time.Hour)
	validRequest := domain.DeleteOddsRequest{
		League:   "English Premier League",
		HomeTeam: "Manchester United",
		AwayTeam: "Liverpool",
		GameDate: gameDate,
	}

	tests := []struct {
		name           string
		request        domain.DeleteOddsRequest
		setupMocks     func()
		expectedResult error
	}{
		{
			name:    "Successful delete",
			request: validRequest,
			setupMocks: func() {
				// Validator validates successfully
				mockValidator.On("ValidateDeleteRequest", validRequest).Return(nil)

				// Repository deletes successfully
				mockRepo.On("Delete", ctx, validRequest.League, validRequest.HomeTeam, validRequest.AwayTeam, validRequest.GameDate).Return(nil)
			},
			expectedResult: nil,
		},
		{
			name:    "Validation error",
			request: validRequest,
			setupMocks: func() {
				// Validator returns error
				mockValidator.On("ValidateDeleteRequest", validRequest).Return(domain.ErrInvalidLeague)
				
				// Repository should not be called
			},
			expectedResult: domain.ErrInvalidLeague,
		},
		{
			name:    "Repository error",
			request: validRequest,
			setupMocks: func() {
				// Validator validates successfully
				mockValidator.On("ValidateDeleteRequest", validRequest).Return(nil)

				// Repository returns error
				mockRepo.On("Delete", ctx, validRequest.League, validRequest.HomeTeam, validRequest.AwayTeam, validRequest.GameDate).Return(domain.ErrOddsNotFound)
			},
			expectedResult: domain.ErrOddsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo = new(MockOddsRepository)
			mockValidator = new(MockOddsValidator)
			service = NewOddsService(mockRepo, mockValidator)

			// Setup mock expectations
			tt.setupMocks()

			// Call the method being tested
			result := service.DeleteOdds(ctx, tt.request)

			// Assert result
			assert.Equal(t, tt.expectedResult, result)

			// Assert all expectations were met (all expected methods were called)
			mockRepo.AssertExpectations(t)
			mockValidator.AssertExpectations(t)
		})
	}
}
