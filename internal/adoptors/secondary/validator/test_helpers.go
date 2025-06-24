package validator

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
