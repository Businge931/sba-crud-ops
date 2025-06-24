package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateOddsUseCase_Execute(t *testing.T) {
	// Save original logger state
	oldOut := logrus.StandardLogger().Out
	oldFormatter := logrus.StandardLogger().Formatter
	oldLevel := logrus.StandardLogger().Level

	// Restore original logger when test is done
	defer func() {
		logrus.SetOutput(oldOut)
		logrus.SetFormatter(oldFormatter)
		logrus.SetLevel(oldLevel)
	}()

	testTime := time.Now()

	tests := []TestCase{
		{
			name: "should successfully create odds",
			deps: TestDependencies{
				mockCreator: &mockOddsCreator{},
			},
			args: TestArgs{
				ctx: context.Background(),
				request: domain.CreateOddsRequest{
					League:   "Premier League",
					HomeTeam: "Arsenal",
					AwayTeam: "Chelsea",
					GameDate: testTime,
				},
			},
			before: func() *test.Hook {
				hook := test.NewGlobal()
				logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
				logrus.SetLevel(logrus.DebugLevel)
				return hook
			},
			after: func() {
				// Cleanup handled in defer
			},
			expectedError: "",
			validateLogs: func(t *testing.T, entries []*logrus.Entry) {
				assert.GreaterOrEqual(t, len(entries), 1, "expected at least 1 log entry")
				startEntry := entries[0]
				assert.Equal(t, "Starting operation", startEntry.Message)
				assert.Equal(t, "create_odds", startEntry.Data["operation"])
				assert.Equal(t, "Premier League", startEntry.Data["league"])
				assert.Equal(t, "Arsenal", startEntry.Data["homeTeam"])
				assert.Equal(t, "Chelsea", startEntry.Data["awayTeam"])

				if len(entries) > 1 {
					successEntry := entries[1]
					assert.Equal(t, "Successfully completed operation", successEntry.Message)
				}
			},
		},
		{
			name: "should handle creation error",
			deps: TestDependencies{
				mockCreator: &mockOddsCreator{},
			},
			args: TestArgs{
				ctx: context.Background(),
				request: domain.CreateOddsRequest{
					League:   "Premier League",
					HomeTeam: "Arsenal",
					AwayTeam: "Chelsea",
					GameDate: testTime,
				},
			},
			before: func() *test.Hook {
				hook := test.NewGlobal()
				logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
				logrus.SetLevel(logrus.DebugLevel)
				return hook
			},
			after: func() {
				// Cleanup handled in defer
			},
			expectedError: "database error",
			validateLogs: func(t *testing.T, entries []*logrus.Entry) {
				assert.GreaterOrEqual(t, len(entries), 1, "expected at least 1 log entry")
				startEntry := entries[0]
				assert.Equal(t, "Starting operation", startEntry.Message)
				assert.Equal(t, "create_odds", startEntry.Data["operation"])

				if len(entries) > 1 {
					errorEntry := entries[1]
					assert.Equal(t, "Operation failed", errorEntry.Message)
					assert.Equal(t, "database error", errorEntry.Data["error"])
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test dependencies
			mockCreator := tc.deps.mockCreator
			if mockCreator == nil {
				mockCreator = &mockOddsCreator{}
			}

			// Setup mock expectations
			if tc.expectedError != "" {
				mockCreator.On("CreateOdds", mock.Anything, mock.Anything).
					Return(errors.New(tc.expectedError))
			} else {
				mockCreator.On("CreateOdds", mock.Anything, mock.Anything).
					Return(nil)
			}

			// Run before hook if provided
			var hook *test.Hook
			if tc.before != nil {
				hook = tc.before()
			}

			// Create use case with mock dependencies
			uc := NewCreateOddsUseCase(mockCreator)

			// Execute with test args
			err := uc.Execute(tc.args.ctx, tc.args.request)

			// Verify errors match expected
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}

			// Verify mock expectations
			mockCreator.AssertExpectations(t)

			// Run after hook if provided
			if tc.after != nil {
				tc.after()
			}

			// Validate logs if validation function is provided
			if tc.validateLogs != nil && hook != nil {
				tc.validateLogs(t, hook.AllEntries())
			}

			// Reset the hook for the next test
			if hook != nil {
				hook.Reset()
			}
		})
	}
}

type TestDependencies struct {
	mockCreator *mockOddsCreator
}

type TestArgs struct {
	ctx     context.Context
	request domain.CreateOddsRequest
}

type TestCase struct {
	name          string
	deps          TestDependencies
	args          TestArgs
	before        func() *test.Hook
	after         func()
	expectedError string
	validateLogs  func(*testing.T, []*logrus.Entry)
}

// mockOddsCreator is a mock implementation of the OddsCreator interface
type mockOddsCreator struct {
	mock.Mock
}

func (m *mockOddsCreator) CreateOdds(ctx context.Context, req domain.CreateOddsRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
