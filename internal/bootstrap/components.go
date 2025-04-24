package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Businge931/sba-crud-ops/internal/core/ports"
	"github.com/Businge931/sba-crud-ops/internal/core/validator"
	"github.com/Businge931/sba-crud-ops/internal/repository/postgres"
	"github.com/Businge931/sba-crud-ops/internal/service"
)

// ApplicationComponents holds all the application components and services
type ApplicationComponents struct {
	// Services
	OddsService ports.OddsService

	// Other components needed by the server
	OddsValidator validator.OddsValidator
}

// SetupComponents initializes and wires all application components
func SetupComponents(cfg *Config, pool *pgxpool.Pool) *ApplicationComponents {
	// Initialize league registry
	leagueRegistry := NewLeagueRegistry(cfg)

	// Initialize validator
	oddsValidator := validator.NewDefaultOddsValidator(leagueRegistry)

	// Initialize repository
	oddsRepo := postgres.NewOddsRepository(pool)

	// Initialize service
	oddsService := service.NewOddsService(oddsRepo, oddsValidator)

	return &ApplicationComponents{
		OddsService:   oddsService,
		OddsValidator: oddsValidator,
	}
}
