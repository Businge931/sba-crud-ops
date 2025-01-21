package ports

import (
	"context"

	"github.com/Businge931/api-gateway/internal/core/domain"
)


type OddsService interface {
	CreateOdds(ctx context.Context, request domain.CreateOddsRequest) error
	ReadOdds(ctx context.Context, request domain.ReadOddsRequest) ([]domain.Odds, error)
	UpdateOdds(ctx context.Context, request domain.CreateOddsRequest) error
	DeleteOdds(ctx context.Context, request domain.DeleteOddsRequest) error
}