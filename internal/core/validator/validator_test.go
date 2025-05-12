package validator

import (
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLeagueRegistry is a mock implementation of the LeagueRegistry interface
type MockLeagueRegistry struct {
	mock.Mock
}

func (m *MockLeagueRegistry) IsSupported(league string) bool {
	args := m.Called(league)
	return args.Bool(0)
}

func (m *MockLeagueRegistry) SupportedLeagues() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockLeagueRegistry) GetSupportedLeagues() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockLeagueRegistry) RegisterLeague(league string) error {
	args := m.Called(league)
	return args.Error(0)
}

func (m *MockLeagueRegistry) UnregisterLeague(league string) error {
	args := m.Called(league)
	return args.Error(0)
}

func TestValidateCreateRequest(t *testing.T) {
	// Create a mock league registry
	mockRegistry := new(MockLeagueRegistry)
	mockRegistry.On("IsSupported", "English Premier League").Return(true)
	mockRegistry.On("IsSupported", "Invalid League").Return(false)

	validator := NewDefaultOddsValidator(mockRegistry)

	// Setup test cases
	tests := []struct {
		name    string
		request domain.CreateOddsRequest
		wantErr error
	}{
		{
			name: "Valid request",
			request: domain.CreateOddsRequest{
				League:          "English Premier League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        time.Now().Add(24 * time.Hour), // Tomorrow
			},
			wantErr: nil,
		},
		{
			name: "Invalid league",
			request: domain.CreateOddsRequest{
				League:          "Invalid League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        time.Now().Add(24 * time.Hour), // Tomorrow
			},
			wantErr: domain.ErrInvalidLeague,
		},
		{
			name: "Empty home team",
			request: domain.CreateOddsRequest{
				League:          "English Premier League",
				HomeTeam:        "",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        time.Now().Add(24 * time.Hour),
			},
			wantErr: domain.ErrEmptyTeamName,
		},
		{
			name: "Empty away team",
			request: domain.CreateOddsRequest{
				League:          "English Premier League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        time.Now().Add(24 * time.Hour),
			},
			wantErr: domain.ErrEmptyTeamName,
		},
		{
			name: "Invalid home team odds",
			request: domain.CreateOddsRequest{
				League:          "English Premier League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 0.5, // Invalid odds value
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        time.Now().Add(24 * time.Hour),
			},
			wantErr: domain.ErrInvalidOddsValue,
		},
		{
			name: "Invalid away team odds",
			request: domain.CreateOddsRequest{
				League:          "English Premier League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 0.9, // Invalid odds value
				DrawOdds:        3.0,
				GameDate:        time.Now().Add(24 * time.Hour),
			},
			wantErr: domain.ErrInvalidOddsValue,
		},
		{
			name: "Invalid draw odds",
			request: domain.CreateOddsRequest{
				League:          "English Premier League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.1,
				DrawOdds:        0.8, // Invalid odds value
				GameDate:        time.Now().Add(24 * time.Hour),
			},
			wantErr: domain.ErrInvalidOddsValue,
		},
		{
			name: "Empty game date",
			request: domain.CreateOddsRequest{
				League:          "English Premier League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        time.Time{}, // Zero time
			},
			wantErr: domain.ErrEmptyGameDate,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreateRequest(tt.request)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateReadRequest(t *testing.T) {
	// Create a mock league registry
	mockRegistry := new(MockLeagueRegistry)
	mockRegistry.On("IsSupported", "English Premier League").Return(true)
	mockRegistry.On("IsSupported", "Invalid League").Return(false)

	validator := NewDefaultOddsValidator(mockRegistry)

	// Setup test cases
	tests := []struct {
		name    string
		request domain.ReadOddsRequest
		wantErr error
	}{
		{
			name: "Valid request",
			request: domain.ReadOddsRequest{
				League: "English Premier League",
				Date:   time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "Invalid league",
			request: domain.ReadOddsRequest{
				League: "Invalid League",
				Date:   time.Now(),
			},
			wantErr: domain.ErrInvalidLeague,
		},
		{
			name: "Empty date",
			request: domain.ReadOddsRequest{
				League: "English Premier League",
				Date:   time.Time{}, // Zero time
			},
			wantErr: domain.ErrEmptyDate,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateReadRequest(tt.request)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDeleteRequest(t *testing.T) {
	// Create a mock league registry
	mockRegistry := new(MockLeagueRegistry)
	mockRegistry.On("IsSupported", "English Premier League").Return(true)
	mockRegistry.On("IsSupported", "Invalid League").Return(false)

	validator := NewDefaultOddsValidator(mockRegistry)

	// Setup test cases
	tests := []struct {
		name    string
		request domain.DeleteOddsRequest
		wantErr error
	}{
		{
			name: "Valid request",
			request: domain.DeleteOddsRequest{
				League:   "English Premier League",
				HomeTeam: "Manchester United",
				AwayTeam: "Liverpool",
				GameDate: time.Now(),
			},
			wantErr: nil,
		},
		{
			name: "Invalid league",
			request: domain.DeleteOddsRequest{
				League:   "Invalid League",
				HomeTeam: "Manchester United",
				AwayTeam: "Liverpool",
				GameDate: time.Now(),
			},
			wantErr: domain.ErrInvalidLeague,
		},
		{
			name: "Empty home team",
			request: domain.DeleteOddsRequest{
				League:   "English Premier League",
				HomeTeam: "",
				AwayTeam: "Liverpool",
				GameDate: time.Now(),
			},
			wantErr: domain.ErrEmptyTeamName,
		},
		{
			name: "Empty away team",
			request: domain.DeleteOddsRequest{
				League:   "English Premier League",
				HomeTeam: "Manchester United",
				AwayTeam: "",
				GameDate: time.Now(),
			},
			wantErr: domain.ErrEmptyTeamName,
		},
		{
			name: "Empty game date",
			request: domain.DeleteOddsRequest{
				League:   "English Premier League",
				HomeTeam: "Manchester United",
				AwayTeam: "Liverpool",
				GameDate: time.Time{}, // Zero time
			},
			wantErr: domain.ErrEmptyGameDate,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateDeleteRequest(tt.request)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateUpdateRequest(t *testing.T) {
	// Create a mock league registry
	mockRegistry := new(MockLeagueRegistry)
	mockRegistry.On("IsSupported", "English Premier League").Return(true)
	mockRegistry.On("IsSupported", "Invalid League").Return(false)

	validator := NewDefaultOddsValidator(mockRegistry)

	// Setup test cases
	tests := []struct {
		name    string
		request domain.CreateOddsRequest
		wantErr error
	}{
		{
			name: "Valid update request",
			request: domain.CreateOddsRequest{
				League:          "English Premier League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        time.Now().Add(24 * time.Hour), // Tomorrow
			},
			wantErr: nil,
		},
		{
			name: "Invalid league for update",
			request: domain.CreateOddsRequest{
				League:          "Invalid League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        time.Now().Add(24 * time.Hour),
			},
			wantErr: domain.ErrInvalidLeague,
		},
		{
			name: "Updated invalid odds",
			request: domain.CreateOddsRequest{
				League:          "English Premier League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 0.9, // Invalid odds value
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        time.Now().Add(24 * time.Hour),
			},
			wantErr: domain.ErrInvalidOddsValue,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateUpdateRequest(tt.request)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
