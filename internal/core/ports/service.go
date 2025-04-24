package ports

import (
	"context"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
)

type OddsCreator interface {
	CreateOdds(ctx context.Context, request domain.CreateOddsRequest) error
}

type OddsRetriever interface {
	ReadOdds(ctx context.Context, request domain.ReadOddsRequest) ([]domain.Odds, error)
}

type OddsUpdater interface {
	UpdateOdds(ctx context.Context, request domain.CreateOddsRequest) error
}

type OddsDeleter interface {
	DeleteOdds(ctx context.Context, request domain.DeleteOddsRequest) error
}

type OddsService interface {
	OddsCreator
	OddsRetriever
	OddsUpdater
	OddsDeleter
}
