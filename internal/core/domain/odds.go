package domain

import "time"

type Odds struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	League          string    `gorm:"type:varchar(100);not null;index" json:"league"`
	HomeTeam        string    `gorm:"type:varchar(100);not null;index" json:"home_team"`
	AwayTeam        string    `gorm:"type:varchar(100);not null;index" json:"away_team"`
	HomeTeamWinOdds float64   `gorm:"type:decimal(10,2);not null" json:"home_team_win_odds"`
	AwayTeamWinOdds float64   `gorm:"type:decimal(10,2);not null" json:"away_team_win_odds"`
	DrawOdds        float64   `gorm:"type:decimal(10,2);not null" json:"draw_odds"`
	GameDate        time.Time `gorm:"type:date;not null;index" json:"game_date"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
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
