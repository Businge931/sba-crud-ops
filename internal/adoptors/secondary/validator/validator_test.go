package validator

import (
	"errors"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
)

type createRequestTestDeps struct {
	mockRegistry *MockLeagueRegistry
}

type createRequestTestArgs struct {
	request domain.CreateOddsRequest
}

type validateCreateRequestTest struct {
	name           string
	deps           createRequestTestDeps
	args           createRequestTestArgs
	before         func(testing.TB, *createRequestTestDeps)
	after          func(testing.TB, *createRequestTestDeps)
	expectedResult any
	expectedErr    error
}

func TestValidateCreateRequest(t *testing.T) {
	// Common test data
	validRequest := domain.CreateOddsRequest{
		League:          "English Premier League",
		HomeTeam:        "Manchester United",
		AwayTeam:        "Liverpool",
		HomeTeamWinOdds: 2.5,
		AwayTeamWinOdds: 2.1,
		DrawOdds:        3.0,
		GameDate:        time.Now().Add(24 * time.Hour), // Tomorrow
	}

	tests := []validateCreateRequestTest{
		{
			name: "Valid request",
			deps: createRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: createRequestTestArgs{
				request: validRequest,
			},
			before: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    nil,
		},
		{
			name: "Invalid league",
			deps: createRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: createRequestTestArgs{
				request: func() domain.CreateOddsRequest {
					r := validRequest
					r.League = "Invalid League"
					return r
				}(),
			},
			before: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "Invalid League").Return(false)
			},
			after: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    errors.New("league: unsupported league: Invalid League."),
		},
		{
			name: "Empty home team",
			deps: createRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: createRequestTestArgs{
				request: func() domain.CreateOddsRequest {
					r := validRequest
					r.HomeTeam = ""
					return r
				}(),
			},
			before: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    errors.New("home_team: home team is required."),
		},
		{
			name: "Empty away team",
			deps: createRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: createRequestTestArgs{
				request: func() domain.CreateOddsRequest {
					r := validRequest
					r.AwayTeam = ""
					return r
				}(),
			},
			before: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    errors.New("away_team: away team is required."),
		},
		{
			name: "Invalid home team odds",
			deps: createRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: createRequestTestArgs{
				request: func() domain.CreateOddsRequest {
					r := validRequest
					r.HomeTeamWinOdds = 0.5
					return r
				}(),
			},
			before: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    domain.ErrInvalidOddsValue,
		},
		{
			name: "Invalid away team odds",
			deps: createRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: createRequestTestArgs{
				request: func() domain.CreateOddsRequest {
					r := validRequest
					r.AwayTeamWinOdds = 0.9
					return r
				}(),
			},
			before: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    domain.ErrInvalidOddsValue,
		},
		{
			name: "Invalid draw odds",
			deps: createRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: createRequestTestArgs{
				request: func() domain.CreateOddsRequest {
					r := validRequest
					r.DrawOdds = 0.8
					return r
				}(),
			},
			before: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    domain.ErrInvalidOddsValue,
		},
		{
			name: "Empty game date",
			deps: createRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: createRequestTestArgs{
				request: func() domain.CreateOddsRequest {
					r := validRequest
					r.GameDate = time.Time{}
					return r
				}(),
			},
			before: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *createRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    errors.New("game_date: game date is required."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test dependencies
			if tt.before != nil {
				tt.before(t, &tt.deps)
			}

			// Create validator with mock dependencies
			validator := NewDefaultOddsValidator(tt.deps.mockRegistry)

			// Execute test
			err := validator.ValidateCreateRequest(tt.args.request)

			// Verify results
			assertError(t, err, tt.expectedErr)

			// Run after function if provided
			if tt.after != nil {
				tt.after(t, &tt.deps)
			}
		})
	}
}

type readRequestTestDeps struct {
	mockRegistry *MockLeagueRegistry
}

type readRequestTestArgs struct {
	request domain.ReadOddsRequest
}

type validateReadRequestTest struct {
	name           string
	deps           readRequestTestDeps
	args           readRequestTestArgs
	before         func(testing.TB, *readRequestTestDeps)
	after          func(testing.TB, *readRequestTestDeps)
	expectedResult any
	expectedErr    error
}

func TestValidateReadRequest(t *testing.T) {
	// Common test data
	now := time.Now()
	validRequest := domain.ReadOddsRequest{
		League: "English Premier League",
		Date:   now,
	}

	tests := []validateReadRequestTest{
		{
			name: "Valid request",
			deps: readRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: readRequestTestArgs{
				request: validRequest,
			},
			before: func(t testing.TB, deps *readRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *readRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    nil,
		},
		{
			name: "Invalid league",
			deps: readRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: readRequestTestArgs{
				request: func() domain.ReadOddsRequest {
					r := validRequest
					r.League = "Invalid League"
					return r
				}(),
			},
			before: func(t testing.TB, deps *readRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "Invalid League").Return(false)
			},
			after: func(t testing.TB, deps *readRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    errors.New("league: unsupported league: Invalid League."),
		},
		{
			name: "Empty date",
			deps: readRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: readRequestTestArgs{
				request: func() domain.ReadOddsRequest {
					r := validRequest
					r.Date = time.Time{}
					return r
				}(),
			},
			before: func(t testing.TB, deps *readRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *readRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    errors.New("date: date is required."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test dependencies
			if tt.before != nil {
				tt.before(t, &tt.deps)
			}

			// Create validator with mock dependencies
			validator := NewDefaultOddsValidator(tt.deps.mockRegistry)

			// Execute test
			err := validator.ValidateReadRequest(tt.args.request)

			// Verify results
			assertError(t, err, tt.expectedErr)

			// Run after function if provided
			if tt.after != nil {
				tt.after(t, &tt.deps)
			}
		})
	}
}

type deleteRequestTestDeps struct {
	mockRegistry *MockLeagueRegistry
}

type deleteRequestTestArgs struct {
	request domain.DeleteOddsRequest
}

type validateDeleteRequestTest struct {
	name           string
	deps           deleteRequestTestDeps
	args           deleteRequestTestArgs
	before         func(testing.TB, *deleteRequestTestDeps)
	after          func(testing.TB, *deleteRequestTestDeps)
	expectedResult any
	expectedErr    error
}

func TestValidateDeleteRequest(t *testing.T) {
	// Common test data
	now := time.Now()
	validRequest := domain.DeleteOddsRequest{
		League:   "English Premier League",
		HomeTeam: "Manchester United",
		AwayTeam: "Liverpool",
		GameDate: now,
	}

	tests := []validateDeleteRequestTest{
		{
			name: "Valid request",
			deps: deleteRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: deleteRequestTestArgs{
				request: validRequest,
			},
			before: func(t testing.TB, deps *deleteRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *deleteRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    nil,
		},
		{
			name: "Invalid league",
			deps: deleteRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: deleteRequestTestArgs{
				request: func() domain.DeleteOddsRequest {
					r := validRequest
					r.League = "Invalid League"
					return r
				}(),
			},
			before: func(t testing.TB, deps *deleteRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "Invalid League").Return(false)
			},
			after: func(t testing.TB, deps *deleteRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    errors.New("league: unsupported league: Invalid League."),
		},
		{
			name: "Empty home team",
			deps: deleteRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: deleteRequestTestArgs{
				request: func() domain.DeleteOddsRequest {
					r := validRequest
					r.HomeTeam = ""
					return r
				}(),
			},
			before: func(t testing.TB, deps *deleteRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *deleteRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    errors.New("home_team: home team is required."),
		},
		{
			name: "Empty away team",
			deps: deleteRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: deleteRequestTestArgs{
				request: func() domain.DeleteOddsRequest {
					r := validRequest
					r.AwayTeam = ""
					return r
				}(),
			},
			before: func(t testing.TB, deps *deleteRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *deleteRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    errors.New("away_team: away team is required."),
		},
		{
			name: "Empty game date",
			deps: deleteRequestTestDeps{
				mockRegistry: new(MockLeagueRegistry),
			},
			args: deleteRequestTestArgs{
				request: func() domain.DeleteOddsRequest {
					r := validRequest
					r.GameDate = time.Time{}
					return r
				}(),
			},
			before: func(t testing.TB, deps *deleteRequestTestDeps) {
				deps.mockRegistry.On("IsSupported", "English Premier League").Return(true)
			},
			after: func(t testing.TB, deps *deleteRequestTestDeps) {
				deps.mockRegistry.AssertExpectations(t)
			},
			expectedResult: nil,
			expectedErr:    errors.New("game_date: game date is required."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test dependencies
			if tt.before != nil {
				tt.before(t, &tt.deps)
			}

			// Create validator with mock dependencies
			validator := NewDefaultOddsValidator(tt.deps.mockRegistry)

			// Execute test
			err := validator.ValidateDeleteRequest(tt.args.request)

			// Verify results
			assertError(t, err, tt.expectedErr)

			// Run after function if provided
			if tt.after != nil {
				tt.after(t, &tt.deps)
			}
		})
	}
}
