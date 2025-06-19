package postgres

import (
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"

	"github.com/stretchr/testify/require"
)

func TestCreateOdds(t *testing.T) {
	testCases := []struct {
		name    string
		before  func(t *testing.T, deps *testDependencies) domain.Odds
		after   func(t *testing.T, deps *testDependencies, odds domain.Odds)
		wantErr bool
	}{
		{
			name: "Successfully create new odds",
			before: func(t *testing.T, deps *testDependencies) domain.Odds {
				return createTestOdds("Premier League", "Chelsea", "Arsenal", 1.5, 2.5, 3.0)
			},
			after: func(t *testing.T, deps *testDependencies, odds domain.Odds) {
				found, err := deps.repo.Read(deps.ctx, odds.League, odds.GameDate)
				require.NoError(t, err)
				require.NotEmpty(t, found)
				require.Equal(t, "Chelsea", found[0].HomeTeam)
				require.Equal(t, "Arsenal", found[0].AwayTeam)
			},
			wantErr: false,
		},
		{
			name: "Fail to create duplicate odds",
			before: func(t *testing.T, deps *testDependencies) domain.Odds {
				// First create the initial odds
				initialOdds := createTestOdds("La Liga", "Barcelona", "Real Madrid", 1.8, 2.0, 3.5)
				err := deps.repo.Create(deps.ctx, &initialOdds)
				require.NoError(t, err, "Failed to create initial odds for duplicate test")
				// Return the same odds to attempt duplicate creation
				return initialOdds
			},
			after: func(t *testing.T, deps *testDependencies, odds domain.Odds) {
				// Verify only one record exists (no duplicates were created)
				found, err := deps.repo.Read(deps.ctx, odds.League, odds.GameDate)
				require.NoError(t, err)
				require.Len(t, found, 1, "Expected only one record, duplicate might have been created")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deps, cleanup := setupTestDependencies(t)
			defer cleanup()

			odds := tc.before(t, &deps)

			err := deps.repo.Create(deps.ctx, &odds)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Create() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && tc.after != nil {
				tc.after(t, &deps, odds)
			}
		})
	}
}

func TestReadOdds(t *testing.T) {
	testCases := []struct {
		name    string
		before  func(t *testing.T, deps *testDependencies) (league string, gameDate time.Time)
		after   func(t *testing.T, found []domain.Odds)
		wantErr bool
	}{
		{
			name: "Successfully read existing odds",
			before: func(t *testing.T, deps *testDependencies) (string, time.Time) {
				odds := createTestOdds("Bundesliga", "Bayern Munich", "Dortmund", 1.6, 2.2, 3.8)
				err := deps.repo.Create(deps.ctx, &odds)
				require.NoError(t, err)
				return odds.League, odds.GameDate
			},
			after: func(t *testing.T, found []domain.Odds) {
				require.Len(t, found, 1)
				require.Equal(t, "Bayern Munich", found[0].HomeTeam)
				require.Equal(t, "Dortmund", found[0].AwayTeam)
			},
			wantErr: false,
		},
		{
			name: "Return empty slice for non-existent league/date",
			before: func(t *testing.T, deps *testDependencies) (string, time.Time) {
				// Create some test data but query for different league/date
				odds := createTestOdds("Ligue 1", "PSG", "Marseille", 1.7, 2.1, 3.9)
				err := deps.repo.Create(deps.ctx, &odds)
				require.NoError(t, err)
				
				// Return non-existent league and date
				return "NonExistentLeague", time.Now().Add(24 * time.Hour)
			},
			after: func(t *testing.T, found []domain.Odds) {
				// Should return empty slice, not nil and no error
				require.NotNil(t, found)
				require.Empty(t, found)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deps, cleanup := setupTestDependencies(t)
			defer cleanup()

			league, gameDate := tc.before(t, &deps)
			found, err := deps.repo.Read(deps.ctx, league, gameDate)

			if (err != nil) != tc.wantErr {
				t.Fatalf("Read() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && tc.after != nil {
				tc.after(t, found)
			}
		})
	}
}

func TestUpdateOdds(t *testing.T) {
	testCases := []struct {
		name    string
		before  func(t *testing.T, deps *testDependencies) domain.Odds
		update  func(odds *domain.Odds)
		verify  func(t *testing.T, deps *testDependencies, odds domain.Odds)
		wantErr bool
	}{
		{
			name: "Successfully update existing odds",
			before: func(t *testing.T, deps *testDependencies) domain.Odds {
				odds := createTestOdds("Bundesliga", "Bayern", "Dortmund", 1.6, 2.2, 3.8)
				require.NoError(t, deps.repo.Create(deps.ctx, &odds))
				return odds
			},
			update: func(odds *domain.Odds) {
				odds.HomeTeamWinOdds = 1.8
				odds.AwayTeamWinOdds = 2.0
				odds.DrawOdds = 4.0
			},
			verify: func(t *testing.T, deps *testDependencies, odds domain.Odds) {
				updated, err := deps.repo.Read(deps.ctx, odds.League, odds.GameDate)
				require.NoError(t, err)
				require.Len(t, updated, 1)
				require.Equal(t, 1.8, updated[0].HomeTeamWinOdds)
				require.Equal(t, 2.0, updated[0].AwayTeamWinOdds)
				require.Equal(t, 4.0, updated[0].DrawOdds)
			},
			wantErr: false,
		},
		{
			name: "Successfully update existing odds",
			before: func(t *testing.T, deps *testDependencies) domain.Odds {
				odds := createTestOdds("Serie A", "Juventus", "AC Milan", 2.1, 1.9, 3.4)
				err := deps.repo.Create(deps.ctx, &odds)
				require.NoError(t, err)
				return odds
			},
			update: func(odds *domain.Odds) {
				odds.HomeTeamWinOdds = 2.3
				odds.AwayTeamWinOdds = 1.8
			},
			verify: func(t *testing.T, deps *testDependencies, odds domain.Odds) {
				found, err := deps.repo.Read(deps.ctx, odds.League, odds.GameDate)
				require.NoError(t, err)
				require.Len(t, found, 1)
				require.Equal(t, 2.3, found[0].HomeTeamWinOdds)
				require.Equal(t, 1.8, found[0].AwayTeamWinOdds)
			},
			wantErr: false,
		},
		{
			name: "Fail to update non-existent odds",
			before: func(t *testing.T, deps *testDependencies) domain.Odds {
				// Create a test record but return a different one to attempt update
				testOdds := createTestOdds("Eredivisie", "Ajax", "PSV", 1.8, 2.0, 3.5)
				err := deps.repo.Create(deps.ctx, &testOdds)
				require.NoError(t, err)
				
				// Return a different odds object that doesn't exist in DB
				nonExistentOdds := createTestOdds("Eredivisie", "Feyenoord", "AZ Alkmaar", 2.1, 1.9, 3.4)
				return nonExistentOdds
			},
			update: func(odds *domain.Odds) {
				// Try to update the non-existent record
				odds.HomeTeamWinOdds = 2.0
				odds.AwayTeamWinOdds = 1.95
			},
			verify: func(t *testing.T, deps *testDependencies, odds domain.Odds) {
				// Verify no new record was created
				found, err := deps.repo.Read(deps.ctx, odds.League, odds.GameDate)
				require.NoError(t, err)
				
				// Should only find the original test record
				require.Len(t, found, 1)
				require.Equal(t, "Ajax", found[0].HomeTeam)
				require.Equal(t, "PSV", found[0].AwayTeam)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deps, cleanup := setupTestDependencies(t)
			defer cleanup()

			odds := tc.before(t, &deps)
			tc.update(&odds)

			err := deps.repo.Update(deps.ctx, &odds)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Update() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && tc.verify != nil {
				tc.verify(t, &deps, odds)
			}
		})
	}
}

func TestDeleteOdds(t *testing.T) {
	testCases := []struct {
		name    string
		before  func(t *testing.T, deps *testDependencies) domain.Odds
		after   func(t *testing.T, deps *testDependencies, odds domain.Odds)
		wantErr bool
	}{
		{
			name: "Successfully delete existing odds",
			before: func(t *testing.T, deps *testDependencies) domain.Odds {
				odds := createTestOdds("MLS", "LA Galaxy", "LAFC", 2.5, 2.5, 3.2)
				err := deps.repo.Create(deps.ctx, &odds)
				require.NoError(t, err)
				return odds
			},
			after: func(t *testing.T, deps *testDependencies, odds domain.Odds) {
				found, err := deps.repo.Read(deps.ctx, odds.League, odds.GameDate)
				require.NoError(t, err)
				require.Empty(t, found, "Expected no records after deletion")
			},
			wantErr: false,
		},
		{
			name: "No error when deleting non-existent odds",
			before: func(t *testing.T, deps *testDependencies) domain.Odds {
				// Create a test record but return a different one to attempt deletion
				testOdds := createTestOdds("J-League", "Kawasaki Frontale", "Yokohama F. Marinos", 2.1, 1.9, 3.4)
				err := deps.repo.Create(deps.ctx, &testOdds)
				require.NoError(t, err)
				
				// Return a different odds object that doesn't exist in DB
				nonExistentOdds := createTestOdds("J-League", "Urawa Red Diamonds", "Kashima Antlers", 2.0, 2.0, 3.5)
				return nonExistentOdds
			},
			after: func(t *testing.T, deps *testDependencies, odds domain.Odds) {
				// The original record should still exist
				found, err := deps.repo.Read(deps.ctx, "J-League", odds.GameDate)
				require.NoError(t, err)
				
				// Should find the original test record
				require.Len(t, found, 1)
				require.Equal(t, "Kawasaki Frontale", found[0].HomeTeam)
				require.Equal(t, "Yokohama F. Marinos", found[0].AwayTeam)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deps, cleanup := setupTestDependencies(t)
			defer cleanup()

			odds := tc.before(t, &deps)

			err := deps.repo.Delete(deps.ctx, odds.League, odds.HomeTeam, odds.AwayTeam, odds.GameDate)
			if (err != nil) != tc.wantErr {
				t.Fatalf("Delete() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && tc.after != nil {
				tc.after(t, &deps, odds)
			} else if tc.wantErr && tc.after != nil {
				tc.after(t, &deps, odds)
			}
		})
	}
}
