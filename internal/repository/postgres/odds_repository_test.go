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

	// Get database connection
	pool, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer pool.Close()

	// Create a clean test table
	err = createTestTable(pool)
	require.NoError(t, err, "Failed to create test table")

	// Create repository
	repo := NewOddsRepository(pool)

	// Setup test context
	ctx := context.Background()

	// Test data
	league := "English Premier League"
	homeTeam := "Manchester United"
	awayTeam := "Liverpool"
	gameDate := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour) // Tomorrow at midnight
	
	// Create test odds
	odds := &domain.Odds{
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

	// Test Create
	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, odds)
		assert.NoError(t, err, "Failed to create odds")
	})

	// Test Read 
	t.Run("Read", func(t *testing.T) {
		result, err := repo.Read(ctx, league, gameDate)
		assert.NoError(t, err, "Failed to read odds")
		assert.NotEmpty(t, result, "No odds returned")
		assert.Equal(t, homeTeam, result[0].HomeTeam, "Home team doesn't match")
		assert.Equal(t, awayTeam, result[0].AwayTeam, "Away team doesn't match")
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		// Modify odds
		odds.HomeTeamWinOdds = 3.0
		odds.AwayTeamWinOdds = 2.5
		odds.DrawOdds = 3.5

		err := repo.Update(ctx, odds)
		assert.NoError(t, err, "Failed to update odds")

		// Read back to verify
		result, err := repo.Read(ctx, league, gameDate)
		assert.NoError(t, err, "Failed to read updated odds")
		assert.NotEmpty(t, result, "No odds returned after update")
		assert.Equal(t, 3.0, result[0].HomeTeamWinOdds, "Home team odds not updated")
		assert.Equal(t, 2.5, result[0].AwayTeamWinOdds, "Away team odds not updated")
		assert.Equal(t, 3.5, result[0].DrawOdds, "Draw odds not updated")
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, league, homeTeam, awayTeam, gameDate)
		assert.NoError(t, err, "Failed to delete odds")

		// Verify it's gone
		result, err := repo.Read(ctx, league, gameDate)
		assert.NoError(t, err, "Error when reading after delete")
		assert.Empty(t, result, "Odds still present after delete")
	})
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
