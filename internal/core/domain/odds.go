package domain

import "time"

type Odds struct {
	ID              int64     `json:"id"`
	League          string    `json:"league"`
	HomeTeam        string    `json:"home_team"`
	AwayTeam        string    `json:"away_team"`
	HomeTeamWinOdds float64   `json:"home_team_win_odds"`
	AwayTeamWinOdds float64   `json:"away_team_win_odds"`
	DrawOdds        float64   `json:"draw_odds"`
	GameDate        time.Time `json:"game_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateOddsRequest struct {
	League          string    `json:"league"`
	HomeTeam        string    `json:"home_team"`
	AwayTeam        string    `json:"away_team"`
	HomeTeamWinOdds float64   `json:"home_team_win_odds"`
	AwayTeamWinOdds float64   `json:"away_team_win_odds"`
	DrawOdds        float64   `json:"draw_odds"`
	GameDate        time.Time `json:"game_date"`
}

type ReadOddsRequest struct {
	League string    `json:"league"`
	Date   time.Time `json:"date"`
}

type DeleteOddsRequest struct {
	League   string    `json:"league"`
	HomeTeam string    `json:"home_team"`
	AwayTeam string    `json:"away_team"`
	GameDate time.Time `json:"game_date"`
}
