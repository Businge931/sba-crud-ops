package usecase

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
)

// CreateOddsUseCase orchestrates the creation of odds
type CreateOddsUseCase struct {
	oddsCreator ports.OddsCreator
	logger      BaseLogger
}

// NewCreateOddsUseCase creates a new instance of CreateOddsUseCase
func NewCreateOddsUseCase(oddsCreator ports.OddsCreator) *CreateOddsUseCase {
	return &CreateOddsUseCase{
		oddsCreator: oddsCreator,
		logger:      BaseLogger{},
	}
}

// Execute processes the create odds request
func (uc *CreateOddsUseCase) Execute(ctx context.Context, request domain.CreateOddsRequest) error {
	const operation = "create_odds"

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
	err := uc.oddsCreator.CreateOdds(ctx, request)
	if err != nil {
		uc.logger.LogError(operation, startTime, err)
		return err
	}

	// Log success
	uc.logger.LogSuccess(operation, startTime, nil)

	return nil
}
