package service

import (
	"context"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOdds(t *testing.T) {
	validRequest := domain.CreateOddsRequest{
		League:          "English Premier League",
		HomeTeam:        "Manchester United",
		AwayTeam:        "Liverpool",
		HomeTeamWinOdds: 2.5,
		AwayTeamWinOdds: 2.1,
		DrawOdds:        3.0,
		GameDate:        time.Now().Add(24 * time.Hour),
	}

	type testArgs struct {
		ctx     context.Context
		request domain.CreateOddsRequest
	}

	tests := []struct {
		name        string
		args        testArgs
		before      func(*testDependencies)
		after       func(*testing.T, *testDependencies, error)
		expectedErr error
	}{
		{
			name: "Successful creation",
			args: testArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(d *testDependencies) {
				d.validator.On("ValidateCreateRequest", validRequest).Return(nil)
				d.repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Odds")).Return(nil)
			},
			after: func(t *testing.T, d *testDependencies, err error) {
				d.repo.AssertExpectations(t)
				d.validator.AssertExpectations(t)
				assert.NoError(t, err)

			},
			expectedErr: nil,
		},
		{
			name: "Validation error",
			args: testArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(d *testDependencies) {
				d.validator.On("ValidateCreateRequest", validRequest).Return(domain.ErrInvalidLeague)
			},
			after: func(t *testing.T, d *testDependencies, err error) {
				d.validator.AssertExpectations(t)
				d.repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
			},
			expectedErr: domain.ErrInvalidLeague,
		},
		{
			name: "Repository error",
			args: testArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(d *testDependencies) {
				d.validator.On("ValidateCreateRequest", validRequest).Return(nil)
				d.repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Odds")).Return(domain.ErrInvalidConfiguration)
			},
			after: func(t *testing.T, d *testDependencies, err error) {
				d.validator.AssertExpectations(t)
				d.repo.AssertExpectations(t)
			},
			expectedErr: domain.ErrInvalidConfiguration,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize dependencies
			deps := testDependencies{
				repo:      new(MockOddsRepository),
				validator: new(MockOddsValidator),
			}
			deps.service = NewOddsService(deps.repo, deps.validator)

			// Setup test case
			if tt.before != nil {
				tt.before(&deps)
			}

			// Execute
			err := deps.service.CreateOdds(tt.args.ctx, tt.args.request)

			// Cleanup
			if tt.after != nil {
				tt.after(t, &deps, err)
			}
		})
	}
}

func TestReadOdds(t *testing.T) {
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

	type testArgs struct {
		ctx     context.Context
		request domain.ReadOddsRequest
	}

	tests := []struct {
		name        string
		args        testArgs
		before      func(*testDependencies)
		after       func(*testing.T, *testDependencies)
		expected    []domain.Odds
		expectedErr error
	}{
		{
			name: "Successful read",
			args: testArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(d *testDependencies) {
				d.validator.On("ValidateReadRequest", validRequest).Return(nil)
				d.repo.On("Read", mock.Anything, validRequest.League, validRequest.Date).Return(expectedOdds, nil)
			},
			after: func(t *testing.T, d *testDependencies) {
				d.validator.AssertExpectations(t)
				d.repo.AssertExpectations(t)
			},
			expected:    expectedOdds,
			expectedErr: nil,
		},
		{
			name: "Validation error",
			args: testArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(d *testDependencies) {
				d.validator.On("ValidateReadRequest", validRequest).Return(domain.ErrInvalidLeague)
			},
			after: func(t *testing.T, d *testDependencies) {
				d.validator.AssertExpectations(t)
				d.repo.AssertNotCalled(t, "Read", mock.Anything, mock.Anything, mock.Anything)
			},
			expected:    nil,
			expectedErr: domain.ErrInvalidLeague,
		},
		{
			name: "Repository error",
			args: testArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(d *testDependencies) {
				d.validator.On("ValidateReadRequest", validRequest).Return(nil)
				d.repo.On("Read", mock.Anything, validRequest.League, validRequest.Date).Return([]domain.Odds{}, domain.ErrOddsNotFound)
			},
			after: func(t *testing.T, d *testDependencies) {
				d.validator.AssertExpectations(t)
				d.repo.AssertExpectations(t)
			},
			expected:    nil,
			expectedErr: domain.ErrOddsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize dependencies
			deps := testDependencies{
				repo:      new(MockOddsRepository),
				validator: new(MockOddsValidator),
			}
			deps.service = NewOddsService(deps.repo, deps.validator)

			// Setup test case
			if tt.before != nil {
				tt.before(&deps)
			}

			// Execute
			result, err := deps.service.ReadOdds(tt.args.ctx, tt.args.request)

			// Verify results
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			// Cleanup
			if tt.after != nil {
				tt.after(t, &deps)
			}
		})
	}
}

func TestUpdateOdds(t *testing.T) {
	validRequest := domain.CreateOddsRequest{
		League:          "English Premier League",
		HomeTeam:        "Manchester United",
		AwayTeam:        "Liverpool",
		HomeTeamWinOdds: 2.5,
		AwayTeamWinOdds: 2.1,
		DrawOdds:        3.0,
		GameDate:        time.Now().Add(24 * time.Hour),
	}

	expectedOdds := &domain.Odds{
		League:          validRequest.League,
		HomeTeam:        validRequest.HomeTeam,
		AwayTeam:        validRequest.AwayTeam,
		HomeTeamWinOdds: validRequest.HomeTeamWinOdds,
		AwayTeamWinOdds: validRequest.AwayTeamWinOdds,
		DrawOdds:        validRequest.DrawOdds,
		GameDate:        validRequest.GameDate,
	}

	type testArgs struct {
		ctx     context.Context
		request domain.CreateOddsRequest
	}

	tests := []struct {
		name        string
		args        testArgs
		before      func(*testDependencies)
		after       func(*testing.T, *testDependencies)
		expected    *domain.Odds
		expectedErr error
	}{
		{
			name: "Successful update",
			args: testArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(d *testDependencies) {
				d.validator.On("ValidateUpdateRequest", validRequest).Return(nil)
				d.repo.On("Update", mock.Anything, mock.MatchedBy(func(odds *domain.Odds) bool {
					return odds.League == validRequest.League &&
						odds.HomeTeam == validRequest.HomeTeam &&
						odds.AwayTeam == validRequest.AwayTeam &&
						odds.HomeTeamWinOdds == validRequest.HomeTeamWinOdds &&
						odds.AwayTeamWinOdds == validRequest.AwayTeamWinOdds &&
						odds.DrawOdds == validRequest.DrawOdds
				})).Return(nil)
			},
			after: func(t *testing.T, d *testDependencies) {
				d.validator.AssertExpectations(t)
				d.repo.AssertExpectations(t)
			},
			expected:    expectedOdds,
			expectedErr: nil,
		},
		{
			name: "Validation error",
			args: testArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(d *testDependencies) {
				d.validator.On("ValidateUpdateRequest", validRequest).Return(domain.ErrInvalidLeague)
			},
			after: func(t *testing.T, d *testDependencies) {
				d.validator.AssertExpectations(t)
				d.repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			},
			expected:    nil,
			expectedErr: domain.ErrInvalidLeague,
		},
		{
			name: "Repository error",
			args: testArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(d *testDependencies) {
				d.validator.On("ValidateUpdateRequest", validRequest).Return(nil)
				d.repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Odds")).Return(domain.ErrOddsNotFound)
			},
			after: func(t *testing.T, d *testDependencies) {
				d.validator.AssertExpectations(t)
				d.repo.AssertExpectations(t)
			},
			expected:    nil,
			expectedErr: domain.ErrOddsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize dependencies
			deps := testDependencies{
				repo:      new(MockOddsRepository),
				validator: new(MockOddsValidator),
			}
			deps.service = NewOddsService(deps.repo, deps.validator)

			// Setup test case
			if tt.before != nil {
				tt.before(&deps)
			}

			// Execute
			err := deps.service.UpdateOdds(tt.args.ctx, tt.args.request)

			// Verify results
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			// Cleanup
			if tt.after != nil {
				tt.after(t, &deps)
			}
		})
	}
}

func TestDeleteOdds(t *testing.T) {
	gameDate := time.Now().Add(24 * time.Hour)
	validRequest := domain.DeleteOddsRequest{
		League:   "English Premier League",
		HomeTeam: "Manchester United",
		AwayTeam: "Liverpool",
		GameDate: gameDate,
	}

	type testArgs struct {
		ctx     context.Context
		request domain.DeleteOddsRequest
	}

	tests := []struct {
		name        string
		args        testArgs
		before      func(*testDependencies)
		after       func(*testing.T, *testDependencies)
		expectedErr error
	}{
		{
			name: "Successful delete",
			args: testArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(d *testDependencies) {
				d.validator.On("ValidateDeleteRequest", validRequest).Return(nil)
				d.repo.On("Delete", mock.Anything, validRequest.League, validRequest.HomeTeam, validRequest.AwayTeam, validRequest.GameDate).Return(nil)
			},
			after: func(t *testing.T, d *testDependencies) {
				d.validator.AssertExpectations(t)
				d.repo.AssertExpectations(t)
			},
			expectedErr: nil,
		},
		{
			name: "Validation error",
			args: testArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(d *testDependencies) {
				d.validator.On("ValidateDeleteRequest", validRequest).Return(domain.ErrInvalidLeague)
			},
			after: func(t *testing.T, d *testDependencies) {
				d.validator.AssertExpectations(t)
				d.repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			},
			expectedErr: domain.ErrInvalidLeague,
		},
		{
			name: "Repository error",
			args: testArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(d *testDependencies) {
				d.validator.On("ValidateDeleteRequest", validRequest).Return(nil)
				d.repo.On("Delete", mock.Anything, validRequest.League, validRequest.HomeTeam, validRequest.AwayTeam, validRequest.GameDate).Return(domain.ErrOddsNotFound)
			},
			after: func(t *testing.T, d *testDependencies) {
				d.validator.AssertExpectations(t)
				d.repo.AssertExpectations(t)
			},
			expectedErr: domain.ErrOddsNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize dependencies
			deps := testDependencies{
				repo:      new(MockOddsRepository),
				validator: new(MockOddsValidator),
			}
			deps.service = NewOddsService(deps.repo, deps.validator)

			// Setup test case
			if tt.before != nil {
				tt.before(&deps)
			}

			// Execute
			err := deps.service.DeleteOdds(tt.args.ctx, tt.args.request)

			// Verify results
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			// Cleanup
			if tt.after != nil {
				tt.after(t, &deps)
			}
		})
	}
}
