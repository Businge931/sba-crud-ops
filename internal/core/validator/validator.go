package validator

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/config"
	"github.com/Businge931/sba-crud-ops/internal/core/domain"
)

// OddsValidator defines the interface for validating odds-related requests
type OddsValidator interface {
	ValidateCreateRequest(request domain.CreateOddsRequest) error
	ValidateReadRequest(request domain.ReadOddsRequest) error
	ValidateUpdateRequest(request domain.CreateOddsRequest) error
	ValidateDeleteRequest(request domain.DeleteOddsRequest) error
}

// DefaultOddsValidator implements the OddsValidator interface with default validation logic
type DefaultOddsValidator struct {
	leagueRegistry config.LeagueRegistry
}

// NewDefaultOddsValidator creates a new instance of DefaultOddsValidator
func NewDefaultOddsValidator(leagueRegistry config.LeagueRegistry) *DefaultOddsValidator {
	// If no registry is provided, create a default one
	if leagueRegistry == nil {
		leagueRegistry = config.NewLeagueRegistry(nil) // Uses default of English Premier League
	}

	return &DefaultOddsValidator{
		leagueRegistry: leagueRegistry,
	}
}

// ValidateCreateRequest validates a create odds request
func (validator *DefaultOddsValidator) ValidateCreateRequest(request domain.CreateOddsRequest) error {
	// Validate league
	if err := validator.validateLeague(request.League); err != nil {
		return err
	}

	// Validate team names
	if err := validator.validateTeam(request.HomeTeam); err != nil {
		return fmt.Errorf("invalid home team: %w", err)
	}

	if err := validator.validateTeam(request.AwayTeam); err != nil {
		return fmt.Errorf("invalid away team: %w", err)
	}

	// Validate odds values
	if err := validator.validateOddsValue(request.HomeTeamWinOdds); err != nil {
		return fmt.Errorf("invalid home team win odds: %w", err)
	}

	if err := validator.validateOddsValue(request.AwayTeamWinOdds); err != nil {
		return fmt.Errorf("invalid away team win odds: %w", err)
	}

	if err := validator.validateOddsValue(request.DrawOdds); err != nil {
		return fmt.Errorf("invalid draw odds: %w", err)
	}

	// Validate game date
	if err := validator.validateGameDate(request.GameDate); err != nil {
		return err
	}

	return nil
}

// ValidateReadRequest validates a read odds request
func (validator *DefaultOddsValidator) ValidateReadRequest(request domain.ReadOddsRequest) error {
	// Validate league
	if err := validator.validateLeague(request.League); err != nil {
		return err
	}

	// Validate date (optional extra validation can be added)
	if request.Date.IsZero() {
		return domain.ErrEmptyDate
	}

	return nil
}

// ValidateUpdateRequest validates an update odds request
func (validator *DefaultOddsValidator) ValidateUpdateRequest(request domain.CreateOddsRequest) error {
	// Reuse create validation since they're similar
	return validator.ValidateCreateRequest(request)
}

// ValidateDeleteRequest validates a delete odds request
func (validator *DefaultOddsValidator) ValidateDeleteRequest(request domain.DeleteOddsRequest) error {
	// Validate league
	if err := validator.validateLeague(request.League); err != nil {
		return err
	}

	// Validate team names
	if err := validator.validateTeam(request.HomeTeam); err != nil {
		return fmt.Errorf("invalid home team: %w", err)
	}

	if err := validator.validateTeam(request.AwayTeam); err != nil {
		return fmt.Errorf("invalid away team: %w", err)
	}

	// Validate game date
	if request.GameDate.IsZero() {
		return domain.ErrEmptyGameDate
	}

	return nil
}

// Helper validation methods
func (validator *DefaultOddsValidator) validateLeague(league string) error {
	if league == "" {
		return domain.ErrEmptyLeagueName
	}

	if !validator.leagueRegistry.IsSupported(league) {
		return domain.ErrInvalidLeague
	}

	return nil
}

func (validator *DefaultOddsValidator) validateTeam(team string) error {
	if team == "" {
		return domain.ErrEmptyTeamName
	}

	// Ensure team name is properly formatted (alphanumeric with spaces)
	matched, _ := regexp.MatchString(`^[A-Za-z0-9\s\-]+$`, team)
	if !matched {
		return domain.ErrInvalidTeamName
	}

	return nil
}

func (validator *DefaultOddsValidator) validateOddsValue(odds float64) error {
	if odds <= 1.0 {
		return domain.ErrInvalidOddsValue
	}

	return nil
}

func (validator *DefaultOddsValidator) validateGameDate(date time.Time) error {
	if date.IsZero() {
		return domain.ErrEmptyGameDate
	}

	// Optional: Add more date validations if needed
	// For example, ensure date is not in the past
	// if date.Before(time.Now()) {
	//     return domain.ErrInvalidStartDate
	// }

	return nil
}
