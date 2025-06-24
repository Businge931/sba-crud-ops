package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateOddsUseCase_Execute(t *testing.T) {
	// Common test data
	testTime := time.Now()
	validRequest := domain.CreateOddsRequest{
		League:          "English Premier League",
		HomeTeam:        "Arsenal",
		AwayTeam:        "Chelsea",
		HomeTeamWinOdds: 2.1,
		AwayTeamWinOdds: 3.2,
		DrawOdds:        3.5,
		GameDate:        testTime,
	}

	tests := []updateOddsTestCase{
		{
			name: "successfully updates odds",
			deps: updateOddsTestDeps{
				mockUpdater: new(mockOddsUpdater),
			},
			args: updateOddsTestArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(t *testing.T, deps *updateOddsTestDeps) {
				deps.setupLogger()

				deps.mockUpdater.On("UpdateOdds", mock.Anything, validRequest).
					Return(nil).
					Once()

				// Clear any existing entries from the hook
				if deps.loggerHook != nil {
					deps.loggerHook.Reset()
				}
			},
			after: func(t *testing.T, deps *updateOddsTestDeps) {
				deps.mockUpdater.AssertExpectations(t)

				// Verify logs
				require.Len(t, deps.loggerHook.Entries, 2, "Expected 2 log entries")

				// Verify start log
				startLog := deps.loggerHook.Entries[0]
				assert.Equal(t, "Starting operation", startLog.Message)
				assert.Equal(t, logrus.InfoLevel, startLog.Level)
				assert.Equal(t, "English Premier League", startLog.Data["league"])
				assert.Equal(t, "Arsenal", startLog.Data["homeTeam"])
				assert.Equal(t, "Chelsea", startLog.Data["awayTeam"])
				assert.NotEmpty(t, startLog.Data["gameDate"])
				assert.Equal(t, "update_odds", startLog.Data["operation"])

				// Verify success log
				successLog := deps.loggerHook.Entries[1]
				assert.Equal(t, "Successfully completed operation", successLog.Message)
				assert.Equal(t, logrus.InfoLevel, successLog.Level)
				assert.NotEmpty(t, successLog.Data["duration"])
				assert.Equal(t, "update_odds", successLog.Data["operation"])

				deps.teardownLogger()
			},
		},
		{
			name: "returns error when updater fails",
			deps: updateOddsTestDeps{
				mockUpdater: new(mockOddsUpdater),
			},
			args: updateOddsTestArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(t *testing.T, deps *updateOddsTestDeps) {
				deps.setupLogger()

				deps.mockUpdater.On("UpdateOdds", mock.Anything, validRequest).
					Return(assert.AnError).
					Once()

				// Clear any existing entries from the hook
				if deps.loggerHook != nil {
					deps.loggerHook.Reset()
				}
			},
			after: func(t *testing.T, deps *updateOddsTestDeps) {
				deps.mockUpdater.AssertExpectations(t)

				// Verify logs
				require.Len(t, deps.loggerHook.Entries, 2, "Expected 2 log entries")

				// Verify start log
				startLog := deps.loggerHook.Entries[0]
				assert.Equal(t, "Starting operation", startLog.Message)
				assert.Equal(t, logrus.InfoLevel, startLog.Level)

				// Verify error log
				errorLog := deps.loggerHook.Entries[1]
				assert.Equal(t, "Operation failed", errorLog.Message)
				assert.Equal(t, logrus.ErrorLevel, errorLog.Level)
				assert.Equal(t, assert.AnError.Error(), errorLog.Data["error"])
				assert.Equal(t, "update_odds", errorLog.Data["operation"])
				assert.NotEmpty(t, errorLog.Data["duration"])

				deps.teardownLogger()
			},
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test
			if tt.before != nil {
				tt.before(t, &tt.deps)
			}

			// Create use case with mock dependencies
			useCase := &UpdateOddsUseCase{
				oddsUpdater: tt.deps.mockUpdater,
				logger:      NewBaseLogger(tt.deps.logger),
			}

			// Execute
			err := useCase.Execute(tt.args.ctx, tt.args.request)

			// Verify results
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// Verify post-conditions
			if tt.after != nil {
				tt.after(t, &tt.deps)
			}
		})
	}
}

// mockOddsUpdater is a mock implementation of the OddsUpdater interface
type mockOddsUpdater struct {
	mock.Mock
}

func (m *mockOddsUpdater) UpdateOdds(ctx context.Context, req domain.CreateOddsRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

type updateOddsTestDeps struct {
	mockUpdater *mockOddsUpdater
	loggerHook  *test.Hook
	logger      *logrus.Logger
}

func (d *updateOddsTestDeps) setupLogger() {
	d.logger = logrus.New()
	d.logger.SetLevel(logrus.DebugLevel)
	d.logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})

	hook := test.NewLocal(d.logger)
	d.loggerHook = hook
}

func (d *updateOddsTestDeps) teardownLogger() {
	if d.loggerHook != nil {
		d.loggerHook.Reset()
	}
	d.logger = nil
}

type updateOddsTestArgs struct {
	ctx     context.Context
	request domain.CreateOddsRequest
}

type updateOddsTestCase struct {
	name          string
	deps          updateOddsTestDeps
	args          updateOddsTestArgs
	before        func(*testing.T, *updateOddsTestDeps)
	after         func(*testing.T, *updateOddsTestDeps)
	expectedError error
}
