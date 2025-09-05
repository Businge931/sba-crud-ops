package ports

import (
	"context"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
)

type (
	OddsReader interface {
		Read(ctx context.Context, league string, date time.Time) ([]domain.Odds, error)
	}

	OddsWriter interface {
		Create(ctx context.Context, odds *domain.Odds) error
		Update(ctx context.Context, odds *domain.Odds) error
		Delete(ctx context.Context, league, homeTeam, awayTeam string, gameDate time.Time) error
	}

	OddsRepository interface {
		OddsReader
		OddsWriter
	}
)
