package validator

import (
	"errors"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"
)

// assertError checks if the error contains the expected error message
func assertError(t *testing.T, err error, wantErr error) {
	if wantErr != nil {
		require.Error(t, err)
		require.Contains(t, err.Error(), wantErr.Error())
	} else {
		require.NoError(t, err)
	}
}

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
			wantErr: errors.New("league: unsupported league: Invalid League."),
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
			wantErr: errors.New("home_team: home team is required."),
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
			wantErr: errors.New("away_team: away team is required."),
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
			wantErr: errors.New("game_date: game date is required."),
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCreateRequest(tt.request)
			assertError(t, err, tt.wantErr)
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
			wantErr: errors.New("league: unsupported league: Invalid League."),
		},
		{
			name: "Empty date",
			request: domain.ReadOddsRequest{
				League: "English Premier League",
				Date:   time.Time{}, // Zero time
			},
			wantErr: errors.New("date: date is required."),
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateReadRequest(tt.request)
			assertError(t, err, tt.wantErr)
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
			wantErr: errors.New("league: unsupported league: Invalid League."),
		},
		{
			name: "Empty home team",
			request: domain.DeleteOddsRequest{
				League:   "English Premier League",
				HomeTeam: "",
				AwayTeam: "Liverpool",
				GameDate: time.Now(),
			},
			wantErr: errors.New("home_team: home team is required."),
		},
		{
			name: "Empty away team",
			request: domain.DeleteOddsRequest{
				League:   "English Premier League",
				HomeTeam: "Manchester United",
				AwayTeam: "",
				GameDate: time.Now(),
			},
			wantErr: errors.New("away_team: away team is required."),
		},
		{
			name: "Empty game date",
			request: domain.DeleteOddsRequest{
				League:   "English Premier League",
				HomeTeam: "Manchester United",
				AwayTeam: "Liverpool",
				GameDate: time.Time{}, // Zero time
			},
			wantErr: errors.New("game_date: game date is required."),
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateDeleteRequest(tt.request)
			assertError(t, err, tt.wantErr)
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
			wantErr: errors.New("league: unsupported league: Invalid League."),
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
			wantErr: errors.New("home_team_win_odds: odds must be greater than 1.0."),
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateUpdateRequest(tt.request)
			assertError(t, err, tt.wantErr)
		})
	}
}
