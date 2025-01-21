package service

import (
	"context"
	"time"

	"github.com/Businge931/api-gateway/internal/core/domain"
	"github.com/Businge931/api-gateway/internal/core/ports"
)

type oddsService struct {
	repo ports.OddsRepository
}

func NewOddsService(repo ports.OddsRepository) ports.OddsService {
	return &oddsService{
		repo: repo,
	}
}

func (s *oddsService) CreateOdds(ctx context.Context, request domain.CreateOddsRequest) error {
	if err := ValidateOddsRequest(request); err != nil {
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

func (s *oddsService) ReadOdds(ctx context.Context, request domain.ReadOddsRequest) ([]domain.Odds, error) {
	if request.League != "English Premier League" {
		return nil, domain.ErrInvalidLeague
	}

	return s.repo.Read(ctx, request.League, request.Date)
}

func (s *oddsService) UpdateOdds(ctx context.Context, request domain.CreateOddsRequest) error {
	if err := ValidateOddsRequest(request); err != nil {
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

func (s *oddsService) DeleteOdds(ctx context.Context, request domain.DeleteOddsRequest) error {
	if request.League != "English Premier League" {
		return domain.ErrInvalidLeague
	}

	return s.repo.Delete(ctx, request.League, request.HomeTeam, request.AwayTeam, request.GameDate)
}
