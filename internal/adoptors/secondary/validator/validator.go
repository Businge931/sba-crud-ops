package validator

import (
	"regexp"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	registry "github.com/Businge931/sba-crud-ops/internal/core/registry"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"
)

// DefaultOddsValidator implements the OddsValidator interface with default validation logic
type DefaultOddsValidator struct {
	leagueRegistry registry.LeagueRegistry
}

// NewDefaultOddsValidator creates a new instance of DefaultOddsValidator
func NewDefaultOddsValidator(leagueRegistry registry.LeagueRegistry) *DefaultOddsValidator {
	// If no registry is provided, create a default one
	if leagueRegistry == nil {
		leagueRegistry = registry.NewLeagueRegistry(nil) // Uses default of English Premier League
	}

	return &DefaultOddsValidator{
		leagueRegistry: leagueRegistry,
	}
}

// validateLeagueRegistry validates if the league registry is properly initialized
func (validator *DefaultOddsValidator) validateLeagueRegistry() error {
	if validator.leagueRegistry == nil {
		return errors.New("league registry is not initialized")
	}
	return nil
}

// ValidateCreateRequest validates a create odds request using ozzo-validation
func (validator *DefaultOddsValidator) ValidateCreateRequest(request domain.CreateOddsRequest) error {
	err := validation.ValidateStruct(&request,
		validation.Field(&request.League,
			validation.Required.Error("league is required"),
			validation.By(validator.validateLeague)),
		validation.Field(&request.HomeTeam,
			validation.Required.Error("home team is required"),
			validation.Match(regexp.MustCompile(`^[A-Za-z0-9\s\-]+$`)).Error("invalid team name format")),
		validation.Field(&request.AwayTeam,
			validation.Required.Error("away team is required"),
			validation.Match(regexp.MustCompile(`^[A-Za-z0-9\s\-]+$`)).Error("invalid team name format")),
		validation.Field(&request.HomeTeamWinOdds,
			validation.Required.Error("home team win odds are required"),
			validation.Min(1.0).Exclusive().Error("odds must be greater than 1.0")),
		validation.Field(&request.AwayTeamWinOdds,
			validation.Required.Error("away team win odds are required"),
			validation.Min(1.0).Exclusive().Error("odds must be greater than 1.0")),
		validation.Field(&request.DrawOdds,
			validation.Required.Error("draw odds are required"),
			validation.Min(1.0).Exclusive().Error("odds must be greater than 1.0")),
		validation.Field(&request.GameDate,
			validation.Required.Error("game date is required"),
			validation.Min(time.Now().Add(-24*time.Hour)).Error("game date cannot be in the past")),
	)

	if err != nil {
		return errors.Wrap(err, "validation failed")
	}

	return nil
}

// ValidateReadRequest validates a read odds request using ozzo-validation
func (validator *DefaultOddsValidator) ValidateReadRequest(request domain.ReadOddsRequest) error {
	err := validation.ValidateStruct(&request,
		validation.Field(&request.League,
			validation.Required.Error("league is required"),
			validation.By(validator.validateLeague)),
		validation.Field(&request.Date,
			validation.Required.Error("date is required")),
	)

	if err != nil {
		return errors.Wrap(err, "validation failed")
	}

	return nil
}

// ValidateUpdateRequest validates an update odds request using ozzo-validation
func (validator *DefaultOddsValidator) ValidateUpdateRequest(request domain.CreateOddsRequest) error {
	// Reuse create validation since they have the same structure
	return validator.ValidateCreateRequest(request)
}

// ValidateDeleteRequest validates a delete odds request using ozzo-validation
func (validator *DefaultOddsValidator) ValidateDeleteRequest(request domain.DeleteOddsRequest) error {
	err := validation.ValidateStruct(&request,
		validation.Field(&request.League,
			validation.Required.Error("league is required"),
			validation.By(validator.validateLeague)),
		validation.Field(&request.HomeTeam,
			validation.Required.Error("home team is required"),
			validation.Match(regexp.MustCompile(`^[A-Za-z0-9\s\-]+$`)).Error("invalid team name format")),
		validation.Field(&request.AwayTeam,
			validation.Required.Error("away team is required"),
			validation.Match(regexp.MustCompile(`^[A-Za-z0-9\s\-]+$`)).Error("invalid team name format")),
		validation.Field(&request.GameDate,
			validation.Required.Error("game date is required")),
	)

	if err != nil {
		return errors.Wrap(err, "validation failed")
	}

	return nil
}

// validateLeague is a custom validation function for league field
func (validator *DefaultOddsValidator) validateLeague(value interface{}) error {
	s, ok := value.(string)
	if !ok || s == "" {
		return errors.New("league name is required")
	}

	if err := validator.validateLeagueRegistry(); err != nil {
		return errors.Wrap(err, "failed to validate league registry")
	}

	if !validator.leagueRegistry.IsSupported(s) {
		return errors.Errorf("unsupported league: %s", s)
	}

	return nil
}
