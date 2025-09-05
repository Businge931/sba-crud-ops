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

func TestReadOddsUseCase_Execute(t *testing.T) {
	// Common test data
	now := time.Now()
	validRequest := domain.ReadOddsRequest{
		League: "English Premier League",
		Date:   now,
	}

	sampleOdds := []domain.Odds{
		{
			League:          "English Premier League",
			HomeTeam:        "Arsenal",
			AwayTeam:        "Chelsea",
			HomeTeamWinOdds: 2.1,
			AwayTeamWinOdds: 3.2,
			DrawOdds:        3.5,
			GameDate:        now,
		},
	}

	tests := []readOddsTestCase{
		{
			name: "successfully retrieves odds",
			deps: readOddsTestDeps{
				mockRetriever: new(MockOddsRetriever),
			},
			args: readOddsTestArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(t *testing.T, deps *readOddsTestDeps) {
				deps.setupLogger()

				deps.mockRetriever.On("ReadOdds", mock.Anything, validRequest).
					Return(sampleOdds, nil).
					Once()

				// Clear any existing entries from the hook
				if deps.loggerHook != nil {
					deps.loggerHook.Reset()
				}
			},
			after: func(t *testing.T, deps *readOddsTestDeps) {
				deps.mockRetriever.AssertExpectations(t)

				// Verify logs
				require.Len(t, deps.loggerHook.Entries, 2, "Expected 2 log entries")

				// Verify start log
				startLog := deps.loggerHook.Entries[0]
				assert.Equal(t, "Starting operation", startLog.Message)
				assert.Equal(t, logrus.InfoLevel, startLog.Level)
				assert.Equal(t, "English Premier League", startLog.Data["league"])
				assert.NotEmpty(t, startLog.Data["date"])
				assert.Equal(t, "read_odds", startLog.Data["operation"])

				// Verify success log
				successLog := deps.loggerHook.Entries[1]
				assert.Equal(t, "Successfully completed operation", successLog.Message)
				assert.Equal(t, logrus.InfoLevel, successLog.Level)
				assert.NotEmpty(t, successLog.Data["duration"])
				assert.Equal(t, 1, successLog.Data["odds_count"])
				deps.teardownLogger()
			},
			expectedOdds: sampleOdds,
		},
		{
			name: "returns error when retriever fails",
			deps: readOddsTestDeps{
				mockRetriever: new(MockOddsRetriever),
			},
			args: readOddsTestArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(t *testing.T, deps *readOddsTestDeps) {
				deps.setupLogger()

				deps.mockRetriever.On("ReadOdds", mock.Anything, validRequest).
					Return(([]domain.Odds)(nil), assert.AnError).
					Once()
			},
			after: func(t *testing.T, deps *readOddsTestDeps) {
				deps.mockRetriever.AssertExpectations(t)

				// Verify error log
				require.Len(t, deps.loggerHook.Entries, 2, "Expected 2 log entries")

				// Verify start log
				startLog := deps.loggerHook.Entries[0]
				assert.Equal(t, "Starting operation", startLog.Message)
				assert.Equal(t, logrus.InfoLevel, startLog.Level)
				assert.Equal(t, "English Premier League", startLog.Data["league"])
				assert.NotEmpty(t, startLog.Data["date"])
				assert.Equal(t, "read_odds", startLog.Data["operation"])

				// Verify error log
				errorLog := deps.loggerHook.Entries[1]
				assert.Equal(t, "Operation failed", errorLog.Message)
				assert.Equal(t, logrus.ErrorLevel, errorLog.Level)
				assert.Equal(t, assert.AnError.Error(), errorLog.Data["error"])
				assert.Equal(t, "read_odds", errorLog.Data["operation"])
				assert.NotEmpty(t, errorLog.Data["duration"])
				deps.teardownLogger()
			},
			expectedErr: assert.AnError,
		},
		{
			name: "handles empty result set",
			deps: readOddsTestDeps{
				mockRetriever: new(MockOddsRetriever),
			},
			args: readOddsTestArgs{
				ctx:     context.Background(),
				request: validRequest,
			},
			before: func(t *testing.T, deps *readOddsTestDeps) {
				deps.setupLogger()

				deps.mockRetriever.On("ReadOdds", mock.Anything, validRequest).
					Return([]domain.Odds{}, nil).
					Once()
			},
			after: func(t *testing.T, deps *readOddsTestDeps) {
				deps.mockRetriever.AssertExpectations(t)

				// Verify logs
				require.Len(t, deps.loggerHook.Entries, 2, "Expected 2 log entries")

				// Verify start log
				startLog := deps.loggerHook.Entries[0]
				assert.Equal(t, "Starting operation", startLog.Message)
				assert.Equal(t, logrus.InfoLevel, startLog.Level)
				assert.Equal(t, "English Premier League", startLog.Data["league"])
				assert.NotEmpty(t, startLog.Data["date"])
				assert.Equal(t, "read_odds", startLog.Data["operation"])

				// Verify success log with empty result
				successLog := deps.loggerHook.Entries[1]
				assert.Equal(t, "Successfully completed operation", successLog.Message)
				assert.Equal(t, logrus.InfoLevel, successLog.Level)
				assert.Equal(t, 0, successLog.Data["odds_count"])
				assert.Equal(t, "read_odds", successLog.Data["operation"])
				assert.NotEmpty(t, successLog.Data["duration"])
				deps.teardownLogger()
			},
			expectedOdds: []domain.Odds{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test
			if tt.before != nil {
				tt.before(t, &tt.deps)
			}

			// Create use case with mock dependencies
			useCase := &ReadOddsUseCase{
				oddsRetriever: tt.deps.mockRetriever,
				logger:        NewBaseLogger(tt.deps.logger),
			}

			// Execute
			result, err := useCase.Execute(tt.args.ctx, tt.args.request)

			// Verify results
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOdds, result)
			}

			// Verify post-conditions
			if tt.after != nil {
				tt.after(t, &tt.deps)
			}
		})
	}
}

// MockOddsRetriever is a mock implementation of the OddsRetriever interface
type MockOddsRetriever struct {
	mock.Mock
}

func (m *MockOddsRetriever) ReadOdds(ctx context.Context, request domain.ReadOddsRequest) ([]domain.Odds, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Odds), args.Error(1)
}

type readOddsTestDeps struct {
	mockRetriever *MockOddsRetriever
	loggerHook    *test.Hook
	logger        *logrus.Logger
}

func (d *readOddsTestDeps) setupLogger() {
	// Create a new logger with test hook
	d.logger = logrus.New()
	d.logger.SetLevel(logrus.DebugLevel)
	d.logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})

	// Add hook for testing
	hook := test.NewLocal(d.logger)
	d.loggerHook = hook
}

func (d *readOddsTestDeps) teardownLogger() {
	// Reset the hook after each test
	if d.loggerHook != nil {
		d.loggerHook.Reset()
	}
	// Reset to default logger
	d.logger = nil
}

type readOddsTestArgs struct {
	ctx     context.Context
	request domain.ReadOddsRequest
}

type readOddsTestCase struct {
	name         string
	deps         readOddsTestDeps
	args         readOddsTestArgs
	before       func(*testing.T, *readOddsTestDeps)
	after        func(*testing.T, *readOddsTestDeps)
	expectedErr  error
	expectedOdds []domain.Odds
}
