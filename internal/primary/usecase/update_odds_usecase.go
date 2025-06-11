package usecase

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
)

// UpdateOddsUseCase orchestrates updating odds information
type UpdateOddsUseCase struct {
	oddsUpdater ports.OddsUpdater
	logger      BaseLogger
}

// NewUpdateOddsUseCase creates a new instance of UpdateOddsUseCase
func NewUpdateOddsUseCase(oddsUpdater ports.OddsUpdater) *UpdateOddsUseCase {
	return &UpdateOddsUseCase{
		oddsUpdater: oddsUpdater,
		logger:      BaseLogger{},
	}
}

// Execute processes the update odds request
func (uc *UpdateOddsUseCase) Execute(ctx context.Context, request domain.CreateOddsRequest) error {
	const operation = "update_odds"

	// Start time for metrics
	startTime := time.Now()

	// Log operation start
	fields := log.Fields{
		"league":   request.League,
		"homeTeam": request.HomeTeam,
		"awayTeam": request.AwayTeam,
		"gameDate": request.GameDate,
	}
	uc.logger.LogStart(operation, fields)

	// Call the domain service
	err := uc.oddsUpdater.UpdateOdds(ctx, request)
	if err != nil {
		uc.logger.LogError(operation, startTime, err)
		return err
	}

	// Log success
	uc.logger.LogSuccess(operation, startTime, nil)

	return nil
}
