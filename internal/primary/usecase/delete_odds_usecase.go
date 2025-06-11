package usecase

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
)

// DeleteOddsUseCase orchestrates deleting odds information
type DeleteOddsUseCase struct {
	oddsDeleter ports.OddsDeleter
	logger      BaseLogger
}

// NewDeleteOddsUseCase creates a new instance of DeleteOddsUseCase
func NewDeleteOddsUseCase(oddsDeleter ports.OddsDeleter) *DeleteOddsUseCase {
	return &DeleteOddsUseCase{
		oddsDeleter: oddsDeleter,
		logger:      BaseLogger{},
	}
}

// Execute processes the delete odds request
func (uc *DeleteOddsUseCase) Execute(ctx context.Context, request domain.DeleteOddsRequest) error {
	const operation = "delete_odds"

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
	err := uc.oddsDeleter.DeleteOdds(ctx, request)
	if err != nil {
		uc.logger.LogError(operation, startTime, err)
		return err
	}

	// Log success
	uc.logger.LogSuccess(operation, startTime, nil)

	return nil
}
