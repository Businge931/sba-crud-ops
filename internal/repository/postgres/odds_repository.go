package postgres

import (
	"context"
	"time"

	"github.com/Businge931/api-gateway/internal/core/domain"
	"github.com/Businge931/api-gateway/internal/core/ports"
	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type oddsRepository struct {
	pool *pgxpool.Pool
}

func NewOddsRepository(pool *pgxpool.Pool) ports.OddsRepository {
	return &oddsRepository{
		pool: pool,
	}
}

func (r *oddsRepository) Create(ctx context.Context, odds *domain.Odds) error {
	query := `
		INSERT INTO odds (league, home_team, away_team, home_team_win_odds, away_team_win_odds, draw_odds, game_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := r.pool.QueryRow(
		ctx,
		query,
		odds.League,
		odds.HomeTeam,
		odds.AwayTeam,
		odds.HomeTeamWinOdds,
		odds.AwayTeamWinOdds,
		odds.DrawOdds,
		odds.GameDate,
		odds.CreatedAt,
		odds.UpdatedAt,
	).Scan(&odds.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *oddsRepository) Read(ctx context.Context, league string, date time.Time) ([]domain.Odds, error) {
	query := `
		SELECT id, league, home_team, away_team, home_team_win_odds, away_team_win_odds, draw_odds, game_date, created_at, updated_at
		FROM odds
		WHERE league = $1 AND DATE(game_date) = DATE($2)`

	rows, err := r.pool.Query(ctx, query, league, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var odds []domain.Odds
	for rows.Next() {
		var odd domain.Odds
		err := rows.Scan(
			&odd.ID,
			&odd.League,
			&odd.HomeTeam,
			&odd.AwayTeam,
			&odd.HomeTeamWinOdds,
			&odd.AwayTeamWinOdds,
			&odd.DrawOdds,
			&odd.GameDate,
			&odd.CreatedAt,
			&odd.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		odds = append(odds, odd)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return odds, nil
}

func (r *oddsRepository) Update(ctx context.Context, odds *domain.Odds) error {
	query := `
		UPDATE odds
		SET home_team_win_odds = $1, away_team_win_odds = $2, draw_odds = $3, updated_at = $4
		WHERE league = $5 AND home_team = $6 AND away_team = $7 AND DATE(game_date) = DATE($8)`

	commandTag, err := r.pool.Exec(
		ctx,
		query,
		odds.HomeTeamWinOdds,
		odds.AwayTeamWinOdds,
		odds.DrawOdds,
		odds.UpdatedAt,
		odds.League,
		odds.HomeTeam,
		odds.AwayTeam,
		odds.GameDate,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrOddsNotFound
	}

	return nil
}

func (r *oddsRepository) Delete(ctx context.Context, league, homeTeam, awayTeam string, gameDate time.Time) error {
	query := `
		DELETE FROM odds
		WHERE league = $1 AND home_team = $2 AND away_team = $3 AND DATE(game_date) = DATE($4)`

	commandTag, err := r.pool.Exec(ctx, query, league, homeTeam, awayTeam, gameDate)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return domain.ErrOddsNotFound
	}

	return nil
}
