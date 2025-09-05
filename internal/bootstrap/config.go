package bootstrap

import (
	"fmt"

	registry "github.com/Businge931/sba-crud-ops/internal/core/registry"
	"github.com/Businge931/sba-crud-ops/internal/env"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	ServicePort string
	ServiceName string
	Version     string
	GatewayAddr string

	// Database configuration
	UseGORM      bool   // Flag to enable/disable GORM
	DBAddr       string // Full connection string
	DBHost       string // Individual connection parameters for GORM
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string

	// Business configuration
	SupportedLeagues []string
}

// LoadConfig loads all application configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		// Server config
		ServicePort: env.GetString("SERVICE_PORT", "50052"),
		ServiceName: env.GetString("SERVICE_NAME", "odds-service"),
		Version:     env.GetString("VERSION", "0.0.1"),
		GatewayAddr: env.GetString("GATEWAY_ADDR", "localhost:8080"),

		// Database config
		UseGORM:      env.GetBool("USE_GORM", true), // Default to using GORM
		DBAddr:       env.GetString("DB_ADDR", "postgresql://admin:adminpassword@localhost:5433/sba_crud_ops?sslmode=disable"),
		DBHost:       env.GetString("DB_HOST", "localhost"),
		DBPort:       env.GetString("DB_PORT", "5433"),
		DBUser:       env.GetString("DB_USER", "admin"),
		DBPassword:   env.GetString("DB_PASSWORD", "adminpassword"),
		DBName:       env.GetString("DB_NAME", "sba_crud_ops"),
		MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	// If DB_ADDR is not set, construct it from individual components
	if cfg.DBAddr == "" {
		cfg.DBAddr = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	}

	// Business config
	cfg.SupportedLeagues = []string{
		"English Premier League",
		"La Liga",
		"Serie A",
		"Bundesliga",
		"Ligue 1",
	}

	return cfg
}

// NewLeagueRegistry creates and configures the league registry
func NewLeagueRegistry(cfg *Config) registry.LeagueRegistry {
	if cfg == nil {
		// Return a default registry with no supported leagues if registry is nil
		return registry.NewLeagueRegistry(nil)
	}
	return registry.NewLeagueRegistry(cfg.SupportedLeagues)
}
