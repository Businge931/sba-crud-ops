package service

import (
	"context"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/adoptors/secondary/validator"
	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
)

type oddsService struct {
	repo      ports.OddsRepository
	validator ports.OddsValidator
}

func NewOddsService(repo ports.OddsRepository, v ports.OddsValidator) ports.OddsService {
	// If no validator is provided, create a default one
	if v == nil {
		v = validator.NewDefaultOddsValidator(nil)
	}

	return &oddsService{
		repo:      repo,
		validator: v,
	}
}

func (s *oddsService) CreateOdds(ctx context.Context, request domain.CreateOddsRequest) error {
	if err := s.validator.ValidateCreateRequest(request); err != nil {
		return err
	}

	odds := &domain.Odds{
		League:          request.League,
		HomeTeam:        request.HomeTeam,
		AwayTeam:        request.AwayTeam,
		HomeTeamWinOdds: request.HomeTeamWinOdds,
		AwayTeamWinOdds: request.AwayTeamWinOdds,
		DrawOdds:        request.DrawOdds,
		GameDate:        request.GameDate,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return s.repo.Create(ctx, odds)
}

// ReadOdds implements OddsRetriever interface
func (s *oddsService) ReadOdds(ctx context.Context, request domain.ReadOddsRequest) ([]domain.Odds, error) {
	if err := s.validator.ValidateReadRequest(request); err != nil {
		return nil, err
	}

	return s.repo.Read(ctx, request.League, request.Date)
}

// UpdateOdds implements OddsUpdater interface
func (s *oddsService) UpdateOdds(ctx context.Context, request domain.CreateOddsRequest) error {
	if err := s.validator.ValidateUpdateRequest(request); err != nil {
		return err
	}

	odds := &domain.Odds{
		League:          request.League,
		HomeTeam:        request.HomeTeam,
		AwayTeam:        request.AwayTeam,
		HomeTeamWinOdds: request.HomeTeamWinOdds,
		AwayTeamWinOdds: request.AwayTeamWinOdds,
		DrawOdds:        request.DrawOdds,
		GameDate:        request.GameDate,
		UpdatedAt:       time.Now(),
	}

	return s.repo.Update(ctx, odds)
}

// DeleteOdds implements OddsDeleter interface
func (s *oddsService) DeleteOdds(ctx context.Context, request domain.DeleteOddsRequest) error {
	if err := s.validator.ValidateDeleteRequest(request); err != nil {
		return err
	}

	return s.repo.Delete(ctx, request.League, request.HomeTeam, request.AwayTeam, request.GameDate)
}
