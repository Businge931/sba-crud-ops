package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOddsRepository is an integration test that requires a real database connection
// These tests should be run with a dedicated test database, not the production database
func TestOddsRepository(t *testing.T) {
	// Skip if we're not running integration tests
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=true to run")
	}

	// Setup test database
	pool, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer pool.Close()

	// Create a clean test table
	err = createTestTable(pool)
	require.NoError(t, err, "Failed to create test table")

	// Create repository
	repo := NewOddsRepository(pool)

	// Test context
	ctx := context.Background()

	// Common test data
	league := "English Premier League"
	homeTeam := "Manchester United"
	awayTeam := "Liverpool"
	gameDate := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour) // Tomorrow at midnight

	// Test cases
	tests := []struct {
		name        string
		setup       func(*testing.T, *oddsRepository) *domain.Odds
		action      func(*testing.T, context.Context, *oddsRepository, *domain.Odds) error
		verify      func(*testing.T, context.Context, *oddsRepository, *domain.Odds, error)
		cleanup     func(*testing.T, context.Context, *oddsRepository, *domain.Odds)
		shouldClean bool
	}{
		{
			name: "Create odds",
			setup: func(t *testing.T, _ *oddsRepository) *domain.Odds {
				return &domain.Odds{
					League:          league,
					HomeTeam:        homeTeam,
					AwayTeam:        awayTeam,
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 2.1,
					DrawOdds:        3.0,
					GameDate:        gameDate,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
			},
			action: func(t *testing.T, ctx context.Context, r *oddsRepository, o *domain.Odds) error {
				return r.Create(ctx, o)
			},
			verify: func(t *testing.T, ctx context.Context, r *oddsRepository, o *domain.Odds, err error) {
				assert.NoError(t, err, "Failed to create odds")
				assert.NotZero(t, o.ID, "Expected odds ID to be set after creation")

				// Verify the odds were created
				result, err := r.Read(ctx, o.League, o.GameDate)
				require.NoError(t, err, "Failed to read created odds")
				require.NotEmpty(t, result, "Expected to find created odds")
				assert.Equal(t, o.HomeTeam, result[0].HomeTeam, "Home team doesn't match")
				assert.Equal(t, o.AwayTeam, result[0].AwayTeam, "Away team doesn't match")
			},
			shouldClean: true,
		},
		{
			name: "Read odds",
			setup: func(t *testing.T, r *oddsRepository) *domain.Odds {
				o := &domain.Odds{
					League:          league,
					HomeTeam:        homeTeam,
					AwayTeam:        awayTeam,
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 2.1,
					DrawOdds:        3.0,
					GameDate:        gameDate,
				}
				require.NoError(t, r.Create(ctx, o), "Failed to setup test odds")
				return o
			},
			action: func(t *testing.T, ctx context.Context, r *oddsRepository, o *domain.Odds) error {
				_, err := r.Read(ctx, o.League, o.GameDate)
				return err
			},
			verify: func(t *testing.T, ctx context.Context, r *oddsRepository, o *domain.Odds, err error) {
				assert.NoError(t, err, "Failed to read odds")
				result, err := r.Read(ctx, o.League, o.GameDate)
				require.NoError(t, err, "Failed to read odds in verification")
				assert.NotEmpty(t, result, "Expected to find odds")
				assert.Equal(t, o.HomeTeam, result[0].HomeTeam, "Home team doesn't match")
				assert.Equal(t, o.AwayTeam, result[0].AwayTeam, "Away team doesn't match")
			},
			shouldClean: true,
		},
		{
			name: "Update odds",
			setup: func(t *testing.T, r *oddsRepository) *domain.Odds {
				o := &domain.Odds{
					League:          league,
					HomeTeam:        homeTeam,
					AwayTeam:        awayTeam,
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 2.1,
					DrawOdds:        3.0,
					GameDate:        gameDate,
				}
				require.NoError(t, r.Create(ctx, o), "Failed to setup test odds")
				return o
			},
			action: func(t *testing.T, ctx context.Context, r *oddsRepository, o *domain.Odds) error {
				o.HomeTeamWinOdds = 3.0
				o.AwayTeamWinOdds = 2.5
				o.DrawOdds = 3.5
				return r.Update(ctx, o)
			},
			verify: func(t *testing.T, ctx context.Context, r *oddsRepository, o *domain.Odds, err error) {
				assert.NoError(t, err, "Failed to update odds")
				result, err := r.Read(ctx, o.League, o.GameDate)
				require.NoError(t, err, "Failed to read updated odds")
				require.NotEmpty(t, result, "No odds returned after update")
				assert.Equal(t, 3.0, result[0].HomeTeamWinOdds, "Home team odds not updated")
				assert.Equal(t, 2.5, result[0].AwayTeamWinOdds, "Away team odds not updated")
				assert.Equal(t, 3.5, result[0].DrawOdds, "Draw odds not updated")
			},
			shouldClean: true,
		},
		{
			name: "Delete odds",
			setup: func(t *testing.T, r *oddsRepository) *domain.Odds {
				o := &domain.Odds{
					League:          league,
					HomeTeam:        homeTeam,
					AwayTeam:        awayTeam,
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 2.1,
					DrawOdds:        3.0,
					GameDate:        gameDate,
				}
				require.NoError(t, r.Create(ctx, o), "Failed to setup test odds")
				return o
			},
			action: func(t *testing.T, ctx context.Context, r *oddsRepository, o *domain.Odds) error {
				return r.Delete(ctx, o.League, o.HomeTeam, o.AwayTeam, o.GameDate)
			},
			verify: func(t *testing.T, ctx context.Context, r *oddsRepository, o *domain.Odds, err error) {
				assert.NoError(t, err, "Failed to delete odds")
				result, err := r.Read(ctx, o.League, o.GameDate)
				require.NoError(t, err, "Error when reading after delete")
				assert.Empty(t, result, "Odds still present after delete")
			},
			shouldClean: false, // Already cleaned up by the action
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the test table for each test case
			err := createTestTable(pool)
			require.NoError(t, err, "Failed to reset test table")

			// Setup test case
			o := tt.setup(t, repo.(*oddsRepository))

			// Execute the action
			err = tt.action(t, ctx, repo.(*oddsRepository), o)

			// Verify the results
			if tt.verify != nil {
				tt.verify(t, ctx, repo.(*oddsRepository), o, err)
			}

			// Cleanup if needed
			if tt.shouldClean && tt.cleanup != nil {
				tt.cleanup(t, ctx, repo.(*oddsRepository), o)
			}
		})
	}
}

// Helper function to setup a test database connection
func setupTestDB() (*pgxpool.Pool, error) {
	// Use environment variable for test database connection or a default for local development
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		testDBURL = "postgres://postgres:postgres@localhost:5432/odds_test?sslmode=disable"
	}

	// Connect to database directly rather than using the db package
	return pgxpool.New(context.Background(), testDBURL)
}

// Helper function to create test table
func createTestTable(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Drop the test table if it exists (clean start)
	_, err := pool.Exec(ctx, `
	DROP TABLE IF EXISTS odds;
	`)
	if err != nil {
		return err
	}

	// Create the test table
	_, err = pool.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS odds (
		id SERIAL PRIMARY KEY,
		league VARCHAR(100) NOT NULL,
		home_team VARCHAR(100) NOT NULL,
		away_team VARCHAR(100) NOT NULL,
		home_team_win_odds NUMERIC(6, 2) NOT NULL CHECK (home_team_win_odds > 1.0),
		away_team_win_odds NUMERIC(6, 2) NOT NULL CHECK (away_team_win_odds > 1.0),
		draw_odds NUMERIC(6, 2) NOT NULL CHECK (draw_odds > 1.0),
		game_date TIMESTAMP WITH TIME ZONE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
		UNIQUE (league, home_team, away_team, game_date)
	);
	`)
	return err
}
