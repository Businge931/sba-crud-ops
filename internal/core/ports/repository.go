package ports

import (
	"context"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
)

// OddsReader represents the read operations for the odds data
type OddsReader interface {
	// Read retrieves odds for a specific league and date
	Read(ctx context.Context, league string, date time.Time) ([]domain.Odds, error)
}

// OddsWriter represents the write operations for the odds data
type OddsWriter interface {
	// Create adds new odds to the data store
	Create(ctx context.Context, odds *domain.Odds) error

	// Update modifies existing odds in the data store
	Update(ctx context.Context, odds *domain.Odds) error

	// Delete removes odds from the data store
	Delete(ctx context.Context, league, homeTeam, awayTeam string, gameDate time.Time) error
}

// OddsRepository combines the read and write interfaces for odds data
type OddsRepository interface {
	OddsReader
	OddsWriter
}
