package ports

import (
	"context"
	"time"

	"github.com/Businge931/api-gateway/internal/core/domain"
)

// type OddsRepository interface {
// 	Create(ctx context.Context, odds *domain.Odds) error
// 	GetByLeagueAndDate(ctx context.Context, league string, date time.Time) ([]domain.Odds, error)
// 	Update(ctx context.Context, odds *domain.Odds) error
// 	Delete(ctx context.Context, league, homeTeam, awayTeam string, gameDate time.Time) error
// }

type OddsRepository interface {
	Create(ctx context.Context, odds *domain.Odds) error
	Read(ctx context.Context, league string, date time.Time) ([]domain.Odds, error)
	Update(ctx context.Context, odds *domain.Odds) error
	Delete(ctx context.Context, league, homeTeam, awayTeam string, gameDate time.Time) error
}
