package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Businge931/sba-crud-ops/internal/bootstrap"
	"github.com/Businge931/sba-crud-ops/internal/db"
)

func main() {
	// Load configuration
	cfg := bootstrap.LoadConfig()

	// Initialize database based on configuration
	var pool *pgxpool.Pool
	var gormDB *gorm.DB
	var err error

	if cfg.UseGORM {
		// Initialize GORM database
		gormDB, err = db.SetupGormDB(cfg)
		if err != nil {
			log.Panicf("Failed to connect to GORM database: %v", err)
		}

		// Run migrations if needed
		err = db.RunMigrations(gormDB)
		if err != nil {
			log.Panicf("Failed to run database migrations: %v", err)
		}

		// For backward compatibility, i'll still create a pgx pool
		// but it won't be used if GORM is enabled
		pool, err = pgxpool.New(context.Background(), cfg.DBAddr)
		if err != nil {
			log.Panicf("Failed to create pgx pool: %v", err)
		}
	} else {
		// Use the original pgx pool
		pool, err = pgxpool.New(context.Background(), cfg.DBAddr)
		if err != nil {
			log.Panicf("Failed to create pgx pool: %v", err)
		}
	}
	defer pool.Close()

	// Setup application components
	components := bootstrap.SetupComponents(cfg, pool, gormDB)

	// Setup server
	server := bootstrap.SetupServer(cfg, components)

	// Try to register with API Gateway in the background
	server.RegisterWithGateway()

	// Start the server
	if err := server.Start(); err != nil {
		log.Panicf("Failed to serve: %v", err)
	}
}
