package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func TestCreateOdds(t *testing.T) {
	// Create mock service
	mockService := new(MockOddsService)
	
	// Create server with mock service
	server := NewOddsServer(mockService)
	
	// Setup test context
	ctx := context.Background()
	
	// Define test time
	testTime := time.Now().Add(24 * time.Hour)
	testDateStr := testTime.Format(time.RFC3339)
	
	// Create valid request
	validRequest := &proto.CreateOddsRequest{
		League:          "English Premier League",
		HomeTeam:        "Manchester United",
		AwayTeam:        "Liverpool",
		HomeTeamWinOdds: 2.5,
		AwayTeamWinOdds: 2.1,
		DrawOdds:        3.0,
		GameDate:        testDateStr,
	}
	
	// Test cases
	tests := []struct {
		name           string
		request        *proto.CreateOddsRequest
		setupMock      func()
		expectedResult *proto.CreateOddsResponse
		expectedError  error
	}{
		{
			name:    "Successful creation",
			request: validRequest,
			setupMock: func() {
				// Expect the service to be called with appropriate domain request
				mockService.On("CreateOdds", ctx, mock.AnythingOfType("domain.CreateOddsRequest")).Return(nil)
			},
			expectedResult: &proto.CreateOddsResponse{OddsId: "created-successfully"},
			expectedError:  nil,
		},
		{
			name:    "Service error",
			request: validRequest,
			setupMock: func() {
				// Simulate service error
				mockService.On("CreateOdds", ctx, mock.AnythingOfType("domain.CreateOddsRequest")).Return(domain.ErrInvalidLeague)
			},
			expectedResult: nil,
			expectedError:  status.Errorf(codes.Internal, "failed to create odds: %v", domain.ErrInvalidLeague),
		},
	}
	
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockService = new(MockOddsService)
			server = NewOddsServer(mockService)
			
			// Setup mock
			tt.setupMock()
			
			// Call the method
			result, err := server.CreateOdds(ctx, tt.request)
			
			// Check error
			if tt.expectedError != nil {
				assert.Error(t, err)
				// Check error code if it's a gRPC status error
				if statusErr, ok := status.FromError(err); ok {
					expectedStatusErr, _ := status.FromError(tt.expectedError)
					assert.Equal(t, expectedStatusErr.Code(), statusErr.Code())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.OddsId, result.OddsId)
			}
			
			// Verify all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetOdds(t *testing.T) {
	// Create mock service
	mockService := new(MockOddsService)
	
	// Create server with mock service
	server := NewOddsServer(mockService)
	
	// Setup test context
	ctx := context.Background()
	
	// Define test time
	testTime := time.Now()
	
	// Create sample domain odds
	sampleOdds := []domain.Odds{
		{
			ID:              1,
			League:          "English Premier League",
			HomeTeam:        "Manchester United",
			AwayTeam:        "Liverpool",
			HomeTeamWinOdds: 2.5,
			AwayTeamWinOdds: 2.1,
			DrawOdds:        3.0,
			GameDate:        testTime,
			CreatedAt:       testTime,
			UpdatedAt:       testTime,
		},
	}
	
	// Create valid request
	validRequest := &proto.GetOddsRequest{
		OddsId: "1", // In this implementation, we don't really use this ID
	}
	
	// Test cases
	tests := []struct {
		name          string
		request       *proto.GetOddsRequest
		setupMock     func()
		expectedError error
	}{
		{
			name:    "Successful retrieval",
			request: validRequest,
			setupMock: func() {
				// Expect the service to be called and return sample odds
				mockService.On("ReadOdds", ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return(sampleOdds, nil)
			},
			expectedError: nil,
		},
		{
			name:    "Service error",
			request: validRequest,
			setupMock: func() {
				// Simulate service error
				mockService.On("ReadOdds", ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return([]domain.Odds{}, domain.ErrOddsNotFound)
			},
			expectedError: status.Errorf(codes.Internal, "failed to get odds: %v", domain.ErrOddsNotFound),
		},
		{
			name:    "Empty results",
			request: validRequest,
			setupMock: func() {
				// Return empty results
				mockService.On("ReadOdds", ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return([]domain.Odds{}, nil)
			},
			expectedError: status.Errorf(codes.NotFound, "odds not found"),
		},
	}
	
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockService = new(MockOddsService)
			server = NewOddsServer(mockService)
			
			// Setup mock
			tt.setupMock()
			
			// Call the method
			result, err := server.GetOdds(ctx, tt.request)
			
			// Check error
			if tt.expectedError != nil {
				assert.Error(t, err)
				// Check error code if it's a gRPC status error
				if statusErr, ok := status.FromError(err); ok {
					expectedStatusErr, _ := status.FromError(tt.expectedError)
					assert.Equal(t, expectedStatusErr.Code(), statusErr.Code())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, sampleOdds[0].HomeTeam, result.HomeTeam)
			}
			
			// Verify all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateOdds(t *testing.T) {
	// Create mock service
	mockService := new(MockOddsService)
	
	// Create server with mock service
	server := NewOddsServer(mockService)
	
	// Setup test context
	ctx := context.Background()
	
	// Define test time
	testTime := time.Now()
	testDateStr := testTime.Format(time.RFC3339)
	
	// Create valid request
	validRequest := &proto.UpdateOddsRequest{
		OddsId:          "1", // Not used in current implementation
		League:          "English Premier League",
		HomeTeam:        "Manchester United",
		AwayTeam:        "Liverpool",
		HomeTeamWinOdds: 2.5,
		AwayTeamWinOdds: 2.1,
		DrawOdds:        3.0,
		GameDate:        testDateStr,
	}
	
	// Test cases
	tests := []struct {
		name           string
		request        *proto.UpdateOddsRequest
		setupMock      func()
		expectedResult *proto.UpdateOddsResponse
		expectedError  error
	}{
		{
			name:    "Successful update",
			request: validRequest,
			setupMock: func() {
				// Expect the service to be called with appropriate domain request
				mockService.On("UpdateOdds", ctx, mock.AnythingOfType("domain.CreateOddsRequest")).Return(nil)
			},
			expectedResult: &proto.UpdateOddsResponse{Success: true},
			expectedError:  nil,
		},
		{
			name:    "Service error",
			request: validRequest,
			setupMock: func() {
				// Simulate service error
				mockService.On("UpdateOdds", ctx, mock.AnythingOfType("domain.CreateOddsRequest")).Return(domain.ErrOddsNotFound)
			},
			expectedResult: nil,
			expectedError:  status.Errorf(codes.Internal, "failed to update odds: %v", domain.ErrOddsNotFound),
		},
	}
	
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockService = new(MockOddsService)
			server = NewOddsServer(mockService)
			
			// Setup mock
			tt.setupMock()
			
			// Call the method
			result, err := server.UpdateOdds(ctx, tt.request)
			
			// Check error
			if tt.expectedError != nil {
				assert.Error(t, err)
				// Check error code if it's a gRPC status error
				if statusErr, ok := status.FromError(err); ok {
					expectedStatusErr, _ := status.FromError(tt.expectedError)
					assert.Equal(t, expectedStatusErr.Code(), statusErr.Code())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Success, result.Success)
			}
			
			// Verify all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestDeleteOdds(t *testing.T) {
	// Create mock service
	mockService := new(MockOddsService)
	
	// Create server with mock service
	server := NewOddsServer(mockService)
	
	// Setup test context
	ctx := context.Background()
	
	// Define test time
	testTime := time.Now()
	
	// Create sample domain odds for the mock to return
	sampleOdds := []domain.Odds{
		{
			ID:              1,
			League:          "English Premier League",
			HomeTeam:        "Manchester United",
			AwayTeam:        "Liverpool",
			HomeTeamWinOdds: 2.5,
			AwayTeamWinOdds: 2.1,
			DrawOdds:        3.0,
			GameDate:        testTime,
			CreatedAt:       testTime,
			UpdatedAt:       testTime,
		},
	}
	
	// Create valid request
	validRequest := &proto.DeleteOddsRequest{
		OddsId: "1",
	}
	
	// Test cases
	tests := []struct {
		name           string
		request        *proto.DeleteOddsRequest
		setupMock      func()
		expectedResult *proto.DeleteOddsResponse
		expectedError  error
	}{
		{
			name:    "Successful deletion",
			request: validRequest,
			setupMock: func() {
				// First we fetch the odds by ReadOdds
				mockService.On("ReadOdds", ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return(sampleOdds, nil)
				// Then we delete them
				mockService.On("DeleteOdds", ctx, mock.AnythingOfType("domain.DeleteOddsRequest")).Return(nil)
			},
			expectedResult: &proto.DeleteOddsResponse{Success: true},
			expectedError:  nil,
		},
		{
			name:    "Not found for deletion",
			request: validRequest,
			setupMock: func() {
				// Return empty results from ReadOdds
				mockService.On("ReadOdds", ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return([]domain.Odds{}, nil)
				// DeleteOdds should not be called
			},
			expectedResult: nil,
			expectedError:  status.Errorf(codes.NotFound, "odds not found for deletion"),
		},
		{
			name:    "Error fetching for deletion",
			request: validRequest,
			setupMock: func() {
				// Simulate error when trying to read
				mockService.On("ReadOdds", ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return([]domain.Odds{}, domain.ErrInvalidLeague)
				// DeleteOdds should not be called
			},
			expectedResult: nil,
			expectedError:  status.Errorf(codes.Internal, "failed to prepare for delete: %v", domain.ErrInvalidLeague),
		},
		{
			name:    "Error during deletion",
			request: validRequest,
			setupMock: func() {
				// Successfully fetch odds
				mockService.On("ReadOdds", ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return(sampleOdds, nil)
				// But fail to delete
				mockService.On("DeleteOdds", ctx, mock.AnythingOfType("domain.DeleteOddsRequest")).Return(domain.ErrInvalidConfiguration)
			},
			expectedResult: nil,
			expectedError:  status.Errorf(codes.Internal, "failed to delete odds: %v", domain.ErrInvalidConfiguration),
		},
	}
	
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock
			mockService = new(MockOddsService)
			server = NewOddsServer(mockService)
			
			// Setup mock
			tt.setupMock()
			
			// Call the method
			result, err := server.DeleteOdds(ctx, tt.request)
			
			// Check error
			if tt.expectedError != nil {
				assert.Error(t, err)
				// Check error code if it's a gRPC status error
				if statusErr, ok := status.FromError(err); ok {
					expectedStatusErr, _ := status.FromError(tt.expectedError)
					assert.Equal(t, expectedStatusErr.Code(), statusErr.Code())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Success, result.Success)
			}
			
			// Verify all expectations were met
			mockService.AssertExpectations(t)
		})
	}
}
