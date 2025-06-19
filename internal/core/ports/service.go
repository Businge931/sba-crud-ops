package ports

import (
	"context"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
)

type (
	OddsCreator interface {
		CreateOdds(ctx context.Context, request domain.CreateOddsRequest) error
	}

	OddsRetriever interface {
		ReadOdds(ctx context.Context, request domain.ReadOddsRequest) ([]domain.Odds, error)
	}

	OddsUpdater interface {
		UpdateOdds(ctx context.Context, request domain.CreateOddsRequest) error
	}

	OddsDeleter interface {
		DeleteOdds(ctx context.Context, request domain.DeleteOddsRequest) error
	}

	OddsService interface {
		OddsCreator
		OddsRetriever
		OddsUpdater
		OddsDeleter
	}

	OddsValidator interface {
		ValidateCreateRequest(request domain.CreateOddsRequest) error
		ValidateReadRequest(request domain.ReadOddsRequest) error
		ValidateUpdateRequest(request domain.CreateOddsRequest) error
		ValidateDeleteRequest(request domain.DeleteOddsRequest) error
	}
)
