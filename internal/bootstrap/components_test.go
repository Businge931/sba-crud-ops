package bootstrap

import (
	"context"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testArgs struct {
	config *Config
}

type testDeps struct {
	pool *pgxpool.Pool
}

type testExpected struct {
	shouldSkip bool
	skipReason string
	shouldErr  bool
}

type testCase struct {
	name   string
	before func(t *testing.T, deps *testDeps, expected *testExpected) *testArgs
	after  func(t *testing.T, deps *testDeps, args *testArgs, expected *testExpected)
}

func TestSetupComponents(t *testing.T) {
	tests := []testCase{
		{
			name: "successful component setup",
			before: func(t *testing.T, deps *testDeps, expected *testExpected) *testArgs {
				// Setup test dependencies
				pool, err := pgxpool.New(context.Background(), "postgresql://test:test@localhost:5432/testdb?sslmode=disable")
				if err != nil {
					expected.shouldSkip = true
					expected.skipReason = "Skipping test as test database is not available"
					return nil
				}
				deps.pool = pool

				// Setup expected results
				expected.shouldErr = false

				// Return test arguments
				return &testArgs{
					config: &Config{
						SupportedLeagues: []string{"Test League"},
					},
				}
			},
			after: func(t *testing.T, deps *testDeps, args *testArgs, expected *testExpected) {
				// Cleanup
				if deps.pool != nil {
					deps.pool.Close()
				}

				// Skip if test was marked to be skipped
				if expected.shouldSkip {
					t.Skip(expected.skipReason)
				}

				// Execute the function under test
				components := SetupComponents(args.config, deps.pool, nil)

				// Verify the results
				require.NotNil(t, components, "Components should not be nil")
				assert.NotNil(t, components.OddsService, "OddsService should not be nil")
				assert.NotNil(t, components.OddsValidator, "OddsValidator should not be nil")

				// Test league registry was properly initialized
				if len(args.config.SupportedLeagues) > 0 {
					testRequest := domain.CreateOddsRequest{
						League:          args.config.SupportedLeagues[0],
						HomeTeam:        "Test Home",
						AwayTeam:        "Test Away",
						HomeTeamWinOdds: 2.0,
						AwayTeamWinOdds: 3.0,
						DrawOdds:        3.5,
						GameDate:        time.Now().Add(24 * time.Hour),
					}
					err := components.OddsValidator.ValidateCreateRequest(testRequest)
					assert.NoError(t, err, "Expected valid request to pass validation")
				}
			},
		},
		{
			name: "nil database pool",
			before: func(t *testing.T, deps *testDeps, expected *testExpected) *testArgs {
				// before test dependencies
				deps.pool = nil

				// before expected results
				expected.shouldErr = true

				// Return test arguments
				return &testArgs{
					config: &Config{
						SupportedLeagues: []string{"Test League"},
					},
				}
			},
			after: func(t *testing.T, deps *testDeps, args *testArgs, expected *testExpected) {
				// Execute the function under test
				components := SetupComponents(args.config, deps.pool, nil)

				// Verify the results
				assert.NotNil(t, components, "Expected non-nil components even with nil pool")
			},
		},
		{
			name: "invalid database connection",
			before: func(t *testing.T, deps *testDeps, expected *testExpected) *testArgs {
				// before test dependencies with an invalid connection
				pool, err := pgxpool.New(context.Background(), "postgresql://test:test@localhost:5432/invalid-db?sslmode=disable")
				if err != nil {
					expected.shouldSkip = true
					expected.skipReason = "Skipping test as test database is not available"
					return nil
				}
				deps.pool = pool

				// before expected results
				expected.shouldErr = false

				// Return test arguments
				return &testArgs{
					config: &Config{
						SupportedLeagues: []string{"Test League"},
					},
				}
			},
			after: func(t *testing.T, deps *testDeps, args *testArgs, expected *testExpected) {
				// Skip if test was marked to be skipped
				if expected.shouldSkip {
					t.Skip(expected.skipReason)
				}

				// Execute the function under test
				components := SetupComponents(args.config, deps.pool, nil)

				// Verify the results
				require.NotNil(t, components, "Components should not be nil")

				// Cleanup
				if deps.pool != nil {
					deps.pool.Close()
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize test dependencies and expected results
			deps := &testDeps{}
			expected := &testExpected{}

			// Run before to get test arguments
			args := tt.before(t, deps, expected)

			// Skip if before indicated we should skip
			if expected.shouldSkip {
				t.Skip(expected.skipReason)
			}

			// Register after for cleanup
			if tt.after != nil {
				t.Cleanup(func() {
					tt.after(t, deps, args, expected)
				})
			}

			// Execute the function under test
			components := SetupComponents(args.config, deps.pool, nil)

			// Store components in expected for after verification if needed
			if tt.after == nil {
				// If no after, run basic assertions here
				if expected.shouldErr {
					assert.NotNil(t, components, "Expected non-nil components even with expected error")
				} else {
					require.NotNil(t, components, "Components should not be nil")
				}
			}
		})
	}
}
