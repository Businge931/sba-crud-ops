package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestOddsRepository is an integration test that requires a real database connection
// These tests should be run with a dedicated test database, not the production database
func TestOddsRepository(t *testing.T) {
	type dependencies struct {
		ctx       context.Context
		container testcontainers.Container
		db        *gorm.DB
		repo      ports.OddsRepository
	}
	type args struct {
		odds domain.Odds
	}

	testCases := []struct {
		name    string
		args    args
		before  func(t *testing.T, deps *dependencies, args *args)
		after   func(t *testing.T, deps *dependencies, args *args)
		want    any
		wantErr bool
	}{
		{
			name: "Create and Read Odds",
			args: args{
				odds: domain.Odds{
					League:          "Premier League",
					HomeTeam:        "Chelsea",
					AwayTeam:        "Arsenal",
					HomeTeamWinOdds: 1.5,
					AwayTeamWinOdds: 2.5,
					DrawOdds:        3.0,
					GameDate:        time.Now().Truncate(24 * time.Hour),
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
			before:  nil,
			after:   nil,
			want:    "Chelsea vs Arsenal",
			wantErr: false,
		},
		{
			name: "Update Odds",
			args: args{
				odds: domain.Odds{
					League:          "Premier League",
					HomeTeam:        "Chelsea",
					AwayTeam:        "Arsenal",
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 2.1,
					DrawOdds:        3.0,
					GameDate:        time.Now().Truncate(24 * time.Hour),
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
			before: func(t *testing.T, deps *dependencies, args *args) {
				require.NoError(t, deps.repo.Create(deps.ctx, &args.odds), "Failed to setup test odds")
			},
			after:   nil,
			want:    3.0, // new HomeTeamWinOdds
			wantErr: false,
		},
		{
			name: "Delete Odds",
			args: args{
				odds: domain.Odds{
					League:          "Premier League",
					HomeTeam:        "Chelsea",
					AwayTeam:        "Arsenal",
					HomeTeamWinOdds: 1.5,
					AwayTeamWinOdds: 2.5,
					DrawOdds:        3.0,
					GameDate:        time.Now().Truncate(24 * time.Hour),
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
			before: func(t *testing.T, deps *dependencies, args *args) {
				require.NoError(t, deps.repo.Create(deps.ctx, &args.odds), "Failed to setup test odds")
			},
			after:   nil,
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test case
			ctx := context.Background()
			container, db, cleanup := setupTestContainerGorm(t, ctx)
			repo := NewOddsRepositoryGorm(db)
			deps := dependencies{
				ctx:       ctx,
				container: container,
				db:        db,
				repo:      repo,
			}
			args := tt.args
			if tt.before != nil {
				tt.before(t, &deps, &args)
			}
			var err error
			switch tt.name {
			case "Create and Read Odds":
				err = repo.Create(ctx, &args.odds)
				if (err != nil) != tt.wantErr {
					t.Errorf("Create error = %v, wantErr %v", err, tt.wantErr)
				}
				found, err := repo.Read(ctx, args.odds.League, args.odds.GameDate)
				if err != nil {
					t.Errorf("failed to read odds: %v", err)
				}
				if len(found) == 0 {
					t.Errorf("no odds found")
				}
				got := fmt.Sprintf("%s vs %s", found[0].HomeTeam, found[0].AwayTeam)
				if got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
			case "Update Odds":
				found, err := repo.Read(ctx, args.odds.League, args.odds.GameDate)
				if err != nil || len(found) == 0 {
					t.Fatalf("failed to find odds for update: %v", err)
				}
				odds := found[0]
				odds.HomeTeamWinOdds = 3.0
				odds.AwayTeamWinOdds = 2.5
				odds.DrawOdds = 3.5
				err = repo.Update(ctx, &odds)
				if (err != nil) != tt.wantErr {
					t.Errorf("Update error = %v, wantErr %v", err, tt.wantErr)
				}
				updated, err := repo.Read(ctx, odds.League, odds.GameDate)
				if err != nil || len(updated) == 0 {
					t.Errorf("failed to read updated odds: %v", err)
				}
				if updated[0].HomeTeamWinOdds != tt.want {
					t.Errorf("got HomeTeamWinOdds=%v, want %v", updated[0].HomeTeamWinOdds, tt.want)
				}
			case "Delete Odds":
				found, err := repo.Read(ctx, args.odds.League, args.odds.GameDate)
				if err != nil || len(found) == 0 {
					t.Fatalf("failed to find odds for delete: %v", err)
				}
				odds := found[0]
				err = repo.Delete(ctx, odds.League, odds.HomeTeam, odds.AwayTeam, odds.GameDate)
				if (err != nil) != tt.wantErr {
					t.Errorf("Delete error = %v, wantErr %v", err, tt.wantErr)
				}
				check, err := repo.Read(ctx, odds.League, odds.GameDate)
				if err == nil && len(check) > 0 {
					t.Errorf("record was not deleted")
				}
				if tt.want != true {
					t.Errorf("expected record to be deleted")
				}
			}
			if tt.after != nil {
				tt.after(t, &deps, &args)
			}
			cleanup()
		})
	}
}

// Helper function to setup a test GORM database connection using testcontainers
func setupTestContainerGorm(t *testing.T, ctx context.Context) (testcontainers.Container, *gorm.DB, func()) {
	t.Helper()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "testuser",
				"POSTGRES_PASSWORD": "testpass",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}
	pgPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		cleanup()
		t.Fatalf("failed to get mapped port: %v", err)
	}
	dsn := fmt.Sprintf("host=localhost port=%d user=testuser password=testpass dbname=testdb sslmode=disable", pgPort.Int())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		cleanup()
		t.Fatalf("failed to connect to database: %v", err)
	}
	if err := db.AutoMigrate(&domain.Odds{}); err != nil {
		cleanup()
		t.Fatalf("failed to migrate database: %v", err)
	}
	return container, db, cleanup
}
