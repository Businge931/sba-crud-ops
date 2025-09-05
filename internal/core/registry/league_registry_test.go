package config

import (
	"testing"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewLeagueRegistry(t *testing.T) {
	type dependencies struct{}
	type args struct {
		initialLeagues []string
	}

	testCases := []struct {
		name         string
		dependencies dependencies
		args         args
		before       func(t *testing.T, deps *dependencies, args *args)
		after        func(t *testing.T, deps *dependencies, args *args)
		want         []string
		wantErr      bool
	}{
		{
			name:         "no initial leagues",
			dependencies: dependencies{},
			args:         args{initialLeagues: nil},
			want:         []string{"English Premier League"},
		},
		{
			name:         "with initial leagues",
			dependencies: dependencies{},
			args:         args{initialLeagues: []string{"La Liga", "Bundesliga"}},
			want:         []string{"La Liga", "Bundesliga"},
		},
		{
			name:         "empty initial leagues",
			dependencies: dependencies{},
			args:         args{initialLeagues: []string{}},
			want:         []string{"English Premier League"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.dependencies, &tc.args)
			}

			r := NewLeagueRegistry(tc.args.initialLeagues)
			leagues := r.GetSupportedLeagues()

			assert.ElementsMatch(t, tc.want, leagues)

			if tc.after != nil {
				tc.after(t, &tc.dependencies, &tc.args)
			}
		})
	}
}

func TestIsSupported(t *testing.T) {
	type dependencies struct{}
	type args struct {
		initialLeagues []string
		league         string
	}

	testCases := []struct {
		name         string
		dependencies dependencies
		args         args
		before       func(t *testing.T, deps *dependencies, args *args)
		after        func(t *testing.T, deps *dependencies, args *args)
		want         bool
		wantErr      bool
	}{
		{
			name:         "league exists",
			dependencies: dependencies{},
			args:         args{initialLeagues: []string{"La Liga"}, league: "La Liga"},
			want:         true,
		},
		{
			name:         "league does not exist",
			dependencies: dependencies{},
			args:         args{initialLeagues: []string{"La Liga"}, league: "Premier League"},
			want:         false,
		},
		{
			name:         "empty league name",
			dependencies: dependencies{},
			args:         args{initialLeagues: nil, league: ""},
			want:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.dependencies, &tc.args)
			}
			r := NewLeagueRegistry(tc.args.initialLeagues)
			result := r.IsSupported(tc.args.league)
			assert.Equal(t, tc.want, result)
			if tc.after != nil {
				tc.after(t, &tc.dependencies, &tc.args)
			}
		})
	}
}

func TestRegisterLeague(t *testing.T) {
	type dependencies struct{}
	type args struct {
		initialLeagues []string
		league         string
	}

	testCases := []struct {
		name         string
		dependencies dependencies
		args         args
		before       func(t *testing.T, deps *dependencies, args *args)
		after        func(t *testing.T, deps *dependencies, args *args)
		wantErr      error
		wantLen      int
	}{
		{
			name:         "register new league",
			dependencies: dependencies{},
			args:         args{initialLeagues: []string{"La Liga"}, league: "Premier League"},
			wantErr:      nil,
			wantLen:      2,
		},
		{
			name:         "register existing league",
			dependencies: dependencies{},
			args:         args{initialLeagues: []string{"La Liga"}, league: "La Liga"},
			wantErr:      nil,
			wantLen:      1, // Should not add duplicate
		},
		{
			name:         "register empty league name",
			dependencies: dependencies{},
			args:         args{initialLeagues: nil, league: ""},
			wantErr:      domain.ErrEmptyLeagueName,
			wantLen:     1, // Default league should still be there
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.dependencies, &tc.args)
			}
			r := NewLeagueRegistry(tc.args.initialLeagues)
			err := r.RegisterLeague(tc.args.league)
			leagues := r.GetSupportedLeagues()
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Len(t, leagues, tc.wantLen)
			if tc.after != nil {
				tc.after(t, &tc.dependencies, &tc.args)
			}
		})
	}
}

func TestUnregisterLeague(t *testing.T) {
	type dependencies struct{}
	type args struct {
		initialLeagues []string
		league         string
	}

	testCases := []struct {
		name         string
		dependencies dependencies
		args         args
		before       func(t *testing.T, deps *dependencies, args *args)
		after        func(t *testing.T, deps *dependencies, args *args)
		wantErr      error
		wantLen      int
	}{
		{
			name:         "unregister existing league",
			dependencies: dependencies{},
			args:         args{initialLeagues: []string{"La Liga"}, league: "La Liga"},
			wantErr:      nil,
			wantLen:      0,
		},
		{
			name:         "unregister non-existent league",
			dependencies: dependencies{},
			args:         args{initialLeagues: []string{"La Liga"}, league: "Premier League"},
			wantErr:      domain.ErrLeagueNotFound,
			wantLen:      1,
		},
		{
			name:         "unregister empty league name",
			dependencies: dependencies{},
			args:         args{initialLeagues: []string{"La Liga"}, league: ""},
			wantErr:      domain.ErrLeagueNotFound,
			wantLen:      1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.dependencies, &tc.args)
			}
			r := NewLeagueRegistry(tc.args.initialLeagues)
			err := r.UnregisterLeague(tc.args.league)
			leagues := r.GetSupportedLeagues()
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Len(t, leagues, tc.wantLen)
			if tc.after != nil {
				tc.after(t, &tc.dependencies, &tc.args)
			}
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
