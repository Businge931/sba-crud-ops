package db

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Businge931/sba-crud-ops/internal/bootstrap"
	"github.com/Businge931/sba-crud-ops/internal/core/domain"
)

// SetupGormDB initializes and configures the GORM database connection
func SetupDB(cfg *bootstrap.Config) (*gorm.DB, error) {
	var dsn string
	if cfg.DBAddr != "" {
		dsn = cfg.DBAddr
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DBHost,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBPort,
		)
	}

	// Configure GORM logger
	gormLogger := logger.New(
		log.StandardLogger(),
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Parse MaxIdleTime duration
	maxIdleDuration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE_TIME format: %v", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(maxIdleDuration)

	log.Info("GORM database connection established")
	return db, nil
}

// RunMigrations runs database migrations for all models
func RunMigrations(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil: cannot run migrations on a nil *gorm.DB")
	}
	// Enable UUID extension if not exists
	err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		return fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	// Auto migrate all models
	err = db.AutoMigrate(
		&domain.Odds{},
	)
	if err != nil {
		return fmt.Errorf("failed to run auto migration: %w", err)
	}

	log.Info("Database migrations completed successfully")
	return nil
}
