package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
)

type oddsRepositoryGorm struct {
	db *gorm.DB
}

// NewOddsRepositoryGorm creates a new instance of oddsRepositoryGorm
func NewOddsRepositoryGorm(db *gorm.DB) ports.OddsRepository {
	return &oddsRepositoryGorm{
		db: db,
	}
}

func (r *oddsRepositoryGorm) Create(ctx context.Context, odds *domain.Odds) error {
	result := r.db.WithContext(ctx).Create(odds)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *oddsRepositoryGorm) Read(ctx context.Context, league string, date time.Time) ([]domain.Odds, error) {
	var odds []domain.Odds

	// Format the date to compare only the date part
	dateStr := date.Format("2006-01-02")

	// Query with date comparison
	result := r.db.WithContext(ctx).
		Where("league = ? AND DATE(game_date) = ?", league, dateStr).
		Find(&odds)

	if result.Error != nil {
		return nil, result.Error
	}

	return odds, nil
}

func (r *oddsRepositoryGorm) Update(ctx context.Context, odds *domain.Odds) error {
	// First, find the existing record to get the ID
	existingOdds := &domain.Odds{}
	err := r.db.WithContext(ctx).
		Where("league = ? AND home_team = ? AND away_team = ? AND DATE(game_date) = ?",
			odds.League, odds.HomeTeam, odds.AwayTeam, odds.GameDate.Format("2006-01-02")).
		First(existingOdds).Error

	if err != nil {
		return err
	}

	// Update the odds with the existing ID and timestamps
	odds.ID = existingOdds.ID
	odds.CreatedAt = existingOdds.CreatedAt
	odds.UpdatedAt = time.Now()

	// Save the updated record
	result := r.db.WithContext(ctx).Save(odds)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *oddsRepositoryGorm) Delete(ctx context.Context, league, homeTeam, awayTeam string, gameDate time.Time) error {
	// Format the date to compare only the date part
	dateStr := gameDate.Format("2006-01-02")

	result := r.db.WithContext(ctx).
		Where("league = ? AND home_team = ? AND away_team = ? AND DATE(game_date) = ?",
			league, homeTeam, awayTeam, dateStr).
		Delete(&domain.Odds{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
