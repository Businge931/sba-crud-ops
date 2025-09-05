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

func TestDeleteOddsUseCase_Execute(t *testing.T) {
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

	tests := []DeleteTestCase{
		{
			name: "should successfully delete odds",
			deps: DeleteTestDependencies{
				mockDeleter: &mockOddsDeleter{},
			},
			args: DeleteTestArgs{
				ctx: context.Background(),
				request: domain.DeleteOddsRequest{
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
				assert.Equal(t, "delete_odds", startEntry.Data["operation"])
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
			name: "should handle deletion error",
			deps: DeleteTestDependencies{
				mockDeleter: &mockOddsDeleter{},
			},
			args: DeleteTestArgs{
				ctx: context.Background(),
				request: domain.DeleteOddsRequest{
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
				assert.Equal(t, "delete_odds", startEntry.Data["operation"])

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
			mockDeleter := tc.deps.mockDeleter
			if mockDeleter == nil {
				mockDeleter = &mockOddsDeleter{}
			}

			// Setup mock expectations
			if tc.expectedError != "" {
				mockDeleter.On("DeleteOdds", mock.Anything, mock.Anything).
					Return(errors.New(tc.expectedError))
			} else {
				mockDeleter.On("DeleteOdds", mock.Anything, mock.Anything).
					Return(nil)
			}

			// Run before hook if provided
			var hook *test.Hook
			if tc.before != nil {
				hook = tc.before()
			}

			// Create use case with mock dependencies
			uc := NewDeleteOddsUseCase(mockDeleter)

			// Execute with test args
			err := uc.Execute(tc.args.ctx, tc.args.request)

			// Verify errors match expected
			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}

			// Verify mock expectations
			mockDeleter.AssertExpectations(t)

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

type DeleteTestDependencies struct {
	mockDeleter *mockOddsDeleter
}

type DeleteTestArgs struct {
	ctx     context.Context
	request domain.DeleteOddsRequest
}

type DeleteTestCase struct {
	name          string
	deps          DeleteTestDependencies
	args          DeleteTestArgs
	before        func() *test.Hook
	after         func()
	expectedError string
	validateLogs  func(*testing.T, []*logrus.Entry)
}

// mockOddsDeleter is a mock implementation of the OddsDeleter interface
type mockOddsDeleter struct {
	mock.Mock
}

func (m *mockOddsDeleter) DeleteOdds(ctx context.Context, req domain.DeleteOddsRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
