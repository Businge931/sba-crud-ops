package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
)

type testDependencies struct {
	ctx       context.Context
	container testcontainers.Container
	db        *gorm.DB
	repo      ports.OddsRepository
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

func createTestOdds(league, homeTeam, awayTeam string, homeOdds, awayOdds, drawOdds float64) domain.Odds {
	now := time.Now()
	return domain.Odds{
		League:          league,
		HomeTeam:        homeTeam,
		AwayTeam:        awayTeam,
		HomeTeamWinOdds: homeOdds,
		AwayTeamWinOdds: awayOdds,
		DrawOdds:        drawOdds,
		GameDate:        now.Truncate(24 * time.Hour),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func setupTestDependencies(t *testing.T) (testDependencies, func()) {
	ctx := context.Background()
	container, db, cleanup := setupTestContainerGorm(t, ctx)
	repo := NewOddsRepositoryGorm(db)
	return testDependencies{
		ctx:       ctx,
		container: container,
		db:        db,
		repo:      repo,
	}, cleanup
}
