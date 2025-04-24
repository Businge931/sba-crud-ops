package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"

	"github.com/Businge931/sba-crud-ops/internal/db"
)

// SetupDatabase initializes and configures the database connection pool
func SetupDatabase(cfg *Config) *pgxpool.Pool {
	pool, err := db.New(
		cfg.DBAddr,
		cfg.MaxOpenConns,
		cfg.MaxIdleConns,
		cfg.MaxIdleTime,
	)
	if err != nil {
		log.Panicf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection pool established")
	return pool
}
