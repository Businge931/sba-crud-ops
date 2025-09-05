package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"

	postgresRepo "github.com/Businge931/sba-crud-ops/internal/adoptors/secondary/postgres"
	"github.com/Businge931/sba-crud-ops/internal/adoptors/secondary/validator"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
	"github.com/Businge931/sba-crud-ops/internal/service"
)

type ApplicationComponents struct {
	// Services
	OddsService ports.OddsService

	// Other components needed by the server
	OddsValidator ports.OddsValidator
}

// SetupComponents initializes and wires all application components
func SetupComponents(cfg *Config, pool *pgxpool.Pool, gormDB *gorm.DB) *ApplicationComponents {
	// Initialize league registry
	leagueRegistry := NewLeagueRegistry(cfg)

	// Initialize validator
	oddsValidator := validator.NewDefaultOddsValidator(leagueRegistry)

	// Always use GORM implementation
	oddsRepo := postgresRepo.NewOddsRepository(gormDB)

	// Initialize service
	oddsService := service.NewOddsService(oddsRepo, oddsValidator)

	return &ApplicationComponents{
		OddsService:   oddsService,
		OddsValidator: oddsValidator,
	}
}
