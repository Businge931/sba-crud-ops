package usecase

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
)

// ReadOddsUseCase orchestrates reading odds information
type ReadOddsUseCase struct {
	oddsRetriever ports.OddsRetriever
	logger        BaseLogger
}

// NewReadOddsUseCase creates a new instance of ReadOddsUseCase
func NewReadOddsUseCase(oddsRetriever ports.OddsRetriever) *ReadOddsUseCase {
	return &ReadOddsUseCase{
		oddsRetriever: oddsRetriever,
		logger:        BaseLogger{},
	}
}

// Execute processes the read odds request
func (uc *ReadOddsUseCase) Execute(ctx context.Context, request domain.ReadOddsRequest) ([]domain.Odds, error) {
	const operation = "read_odds"

	// Start time for metrics
	startTime := time.Now()

	// Log operation start
	fields := log.Fields{
		"league": request.League,
		"date":   request.Date,
	}
	uc.logger.LogStart(operation, fields)

	// Call the domain service
	odds, err := uc.oddsRetriever.ReadOdds(ctx, request)
	if err != nil {
		uc.logger.LogError(operation, startTime, err)
		return nil, err
	}

	// Log success with extra field for odds count
	extraFields := log.Fields{
		"odds_count": len(odds),
	}
	uc.logger.LogSuccess(operation, startTime, extraFields)

	return odds, nil
}
