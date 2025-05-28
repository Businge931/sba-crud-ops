package config

import (
	"testing"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewLeagueRegistry(t *testing.T) {
	tests := []struct {
		name            string
		initialLeagues  []string
		expectedLeagues []string
	}{
		{
			name:            "no initial leagues",
			initialLeagues:  nil,
			expectedLeagues: []string{"English Premier League"},
		},
		{
			name:            "with initial leagues",
			initialLeagues:  []string{"La Liga", "Bundesliga"},
			expectedLeagues: []string{"La Liga", "Bundesliga"},
		},
		{
			name:            "empty initial leagues",
			initialLeagues:  []string{},
			expectedLeagues: []string{"English Premier League"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			r := NewLeagueRegistry(tt.initialLeagues)

			// Execute
			leagues := r.GetSupportedLeagues()

			// Verify
			assert.ElementsMatch(t, tt.expectedLeagues, leagues)
		})
	}
}

func TestIsSupported(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *DefaultLeagueRegistry
		league         string
		expectedResult bool
	}{
		{
			name: "league exists",
			setup: func() *DefaultLeagueRegistry {
				r := NewLeagueRegistry([]string{"La Liga"})
				return r
			},
			league:         "La Liga",
			expectedResult: true,
		},
		{
			name: "league does not exist",
			setup: func() *DefaultLeagueRegistry {
				return NewLeagueRegistry([]string{"La Liga"})
			},
			league:         "Premier League",
			expectedResult: false,
		},
		{
			name: "empty league name",
			setup: func() *DefaultLeagueRegistry {
				return NewLeagueRegistry(nil)
			},
			league:         "",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			r := tt.setup()

			// Execute
			result := r.IsSupported(tt.league)

			// Verify
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestRegisterLeague(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *DefaultLeagueRegistry
		league        string
		expectedError error
		expectedLen   int
	}{
		{
			name: "register new league",
			setup: func() *DefaultLeagueRegistry {
				return NewLeagueRegistry([]string{"La Liga"})
			},
			league:        "Premier League",
			expectedError: nil,
			expectedLen:   2,
		},
		{
			name: "register existing league",
			setup: func() *DefaultLeagueRegistry {
				return NewLeagueRegistry([]string{"La Liga"})
			},
			league:        "La Liga",
			expectedError: nil,
			expectedLen:   1, // Should not add duplicate
		},
		{
			name: "register empty league name",
			setup: func() *DefaultLeagueRegistry {
				return NewLeagueRegistry(nil)
			},
			league:        "",
			expectedError: domain.ErrEmptyLeagueName,
			expectedLen:   1, // Default league should still be there
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			r := tt.setup()

			// Execute
			err := r.RegisterLeague(tt.league)
			leagues := r.GetSupportedLeagues()

			// Verify
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Len(t, leagues, tt.expectedLen)
		})
	}
}

func TestUnregisterLeague(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *DefaultLeagueRegistry
		league        string
		expectedError error
		expectedLen   int
	}{
		{
			name: "unregister existing league",
			setup: func() *DefaultLeagueRegistry {
				return NewLeagueRegistry([]string{"La Liga"})
			},
			league:        "La Liga",
			expectedError: nil,
			expectedLen:   0,
		},
		{
			name: "unregister non-existent league",
			setup: func() *DefaultLeagueRegistry {
				return NewLeagueRegistry([]string{"La Liga"})
			},
			league:        "Premier League",
			expectedError: domain.ErrLeagueNotFound,
			expectedLen:   1,
		},
		{
			name: "unregister empty league name",
			setup: func() *DefaultLeagueRegistry {
				return NewLeagueRegistry([]string{"La Liga"})
			},
			league:        "",
			expectedError: domain.ErrLeagueNotFound,
			expectedLen:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			r := tt.setup()

			// Execute
			err := r.UnregisterLeague(tt.league)
			leagues := r.GetSupportedLeagues()

			// Verify
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Len(t, leagues, tt.expectedLen)
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	r := NewLeagueRegistry([]string{"La Liga"})

	// Number of concurrent operations
	const numOps = 100
	done := make(chan bool)

	// Start multiple goroutines to simulate concurrent access
	for i := 0; i < numOps; i++ {
		go func(n int) {
			// Alternate between register and unregister
			if n%2 == 0 {
				_ = r.RegisterLeague("League" + string(rune('A'+n)))
			} else {
				_ = r.UnregisterLeague("La Liga")
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numOps; i++ {
		<-done
	}

	// Verify the registry is in a valid state
	leagues := r.GetSupportedLeagues()
	assert.True(t, len(leagues) > 0, "Registry should not be empty after concurrent access")
}
