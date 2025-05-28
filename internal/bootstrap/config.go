package bootstrap

import (
	"github.com/Businge931/sba-crud-ops/internal/core/config"
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
	DBAddr       string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string

	// Business configuration
	SupportedLeagues []string
}

// LoadConfig loads all application configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		// Server config
		ServicePort: env.GetString("SERVICE_PORT", "50052"),
		ServiceName: env.GetString("SERVICE_NAME", "odds-service"),
		Version:     env.GetString("VERSION", "0.0.1"),
		GatewayAddr: env.GetString("GATEWAY_ADDR", "localhost:8080"),

		// Database config
		DBAddr:       env.GetString("DB_ADDR", "postgresql://admin:adminpassword@localhost:5433/sba_crud_ops?sslmode=disable"),
		MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),

		// Business config
		SupportedLeagues: []string{
			"English Premier League",
			"La Liga",
			"Serie A",
			"Bundesliga",
			"Ligue 1",
		},
	}
}

// NewLeagueRegistry creates and configures the league registry
func NewLeagueRegistry(cfg *Config) config.LeagueRegistry {
	if cfg == nil {
		// Return a default registry with no supported leagues if config is nil
		return config.NewLeagueRegistry(nil)
	}
	return config.NewLeagueRegistry(cfg.SupportedLeagues)
}
