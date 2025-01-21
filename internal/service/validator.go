package service

import (
	"time"

	"github.com/Businge931/api-gateway/internal/core/domain"
)

func ValidateOddsRequest(request domain.CreateOddsRequest) error {
	if request.League != "English Premier League" {
		return domain.ErrInvalidLeague
	}
	if request.HomeTeam == "" || request.AwayTeam == "" {
		return domain.ErrEmptyTeams
	}
	if request.HomeTeamWinOdds <= 0 || request.AwayTeamWinOdds <= 0 || request.DrawOdds <= 0 {
		return domain.ErrInvalidOdds
	}
	if request.GameDate.Before(time.Now()) {
		return domain.ErrInvalidStartDate
	}
	return nil
}
