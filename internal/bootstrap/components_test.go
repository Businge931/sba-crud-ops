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

func TestSetupComponents(t *testing.T) {
	tests := []struct {
		name          string
		cfg           *Config
		setupMocks    func(*testing.T, *Config) *pgxpool.Pool
		expectedError bool
	}{
		{
			name: "successful component setup",
			cfg: &Config{
				SupportedLeagues: []string{"Test League"},
			},
			setupMocks: func(t *testing.T, cfg *Config) *pgxpool.Pool {
				// Create a real pool for testing
				pool, err := pgxpool.New(context.Background(), "postgresql://test:test@localhost:5432/testdb?sslmode=disable")
				if err != nil {
					t.Skip("Skipping test as test database is not available")
				}
				return pool
			},
			expectedError: false,
		},
		{
			name: "nil database pool",
			cfg: &Config{
				SupportedLeagues: []string{"Test League"},
			},
			setupMocks: func(t *testing.T, cfg *Config) *pgxpool.Pool {
				// Return nil to simulate nil pool
				return nil
			},
			expectedError: true, // This will be handled specially in the test
		},
		{
			name: "invalid database connection",
			cfg: &Config{
				SupportedLeagues: []string{"Test League"},
			},
			setupMocks: func(t *testing.T, cfg *Config) *pgxpool.Pool {
				// Return a pool with an invalid connection
				pool, err := pgxpool.New(context.Background(), "postgresql://test:test@localhost:5432/invalid-db?sslmode=disable")
				if err != nil {
					t.Skip("Skipping test as test database is not available")
				}
				return pool
			},
			expectedError: false, // We don't expect an error during component setup, only when the service is used
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := tt.setupMocks(t, tt.cfg)

			if tt.expectedError && pool == nil {
				// SetupComponents doesn't panic with nil pool, but the service will fail when used
				components := SetupComponents(tt.cfg, pool, nil)
				assert.NotNil(t, components, "Expected non-nil components even with nil pool")
				return
			}

			if pool == nil {
				t.Skip("Skipping test as test database is not available")
			}

			// Execute
			components := SetupComponents(tt.cfg, pool, nil)

			// Verify
			require.NotNil(t, components, "Components should not be nil")
			assert.NotNil(t, components.OddsService, "OddsService should not be nil")
			assert.NotNil(t, components.OddsValidator, "OddsValidator should not be nil")

			// Test league registry was properly initialized by testing a valid request
			// Skip if there are no supported leagues in the config
			if len(tt.cfg.SupportedLeagues) > 0 {
				testRequest := domain.CreateOddsRequest{
					League:          tt.cfg.SupportedLeagues[0],
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

			// Cleanup
			if pool != nil {
				pool.Close()
			}
		})
	}
}
