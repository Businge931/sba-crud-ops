package domain

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// ================================================================
// Application Error Definitions
// ================================================================

// Data Not Found Errors
var (
	// ErrOddsNotFound is returned when the requested odds record cannot be found
	ErrOddsNotFound = errors.New("odds not found")
	
	// ErrLeagueNotFound is returned when a specified league doesn't exist
	ErrLeagueNotFound = errors.New("league not found")
)

// League-related Errors
var (
	// ErrInvalidLeague is returned when an unsupported league is specified
	ErrInvalidLeague = errors.New("invalid league")
	
	// ErrEmptyLeagueName is returned when a league name is required but not provided
	ErrEmptyLeagueName = errors.New("league name cannot be empty")
)

// Team-related Errors
var (
	// ErrEmptyTeams is returned when team names are required but not provided
	ErrEmptyTeams = errors.New("team names cannot be empty")
	
	// ErrEmptyTeamName is returned when a single team name is required but not provided
	ErrEmptyTeamName = errors.New("team name cannot be empty")
	
	// ErrInvalidTeamName is returned when a team name contains invalid characters
	ErrInvalidTeamName = errors.New("team name can only contain letters, numbers, spaces, and hyphens")
	
	// ErrSameTeams is returned when the home and away teams specified are the same
	ErrSameTeams = errors.New("home and away teams cannot be the same")
)

// Odds-related Errors
var (
	// ErrInvalidOdds is returned when odds values don't meet basic validity criteria
	ErrInvalidOdds = errors.New("odds must be positive numbers")
	
	// ErrInvalidOddsValue is returned when an odds value is <= 1.0
	ErrInvalidOddsValue = errors.New("odds must be greater than 1.0")
	
	// ErrInvalidOddsProbability is returned when the implied probabilities don't make sense
	ErrInvalidOddsProbability = errors.New("sum of odds probabilities must be valid")
)

// Date-related Errors
var (
	// ErrEmptyDate is returned when a date field is required but not provided
	ErrEmptyDate = errors.New("date cannot be empty")
	
	// ErrEmptyGameDate is returned when a game date is required but not provided
	ErrEmptyGameDate = errors.New("game date cannot be empty")
	
	// ErrInvalidStartDate is returned when a game date is in the past
	ErrInvalidStartDate = errors.New("game date must be in the future")
	
	// ErrInvalidDateFormat is returned when a date string doesn't match the expected format
	ErrInvalidDateFormat = errors.New("invalid date format")
)

// Configuration Errors
var (
	// ErrInvalidConfiguration is returned when application configuration is invalid
	ErrInvalidConfiguration = errors.New("invalid configuration")
)

// ================================================================
// HTTP Response Functions
// ================================================================

// InternalServerError logs the error and sends a 500 Internal Server Error response
// with a generic error message to avoid exposing sensitive information
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"error":  err.Error(),
	}).Error("internal server error")

	WriteJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

// BadRequestResponse logs the error and sends a 400 Bad Request response
// with the specific error message to help clients correct their request
func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path, 
		"error":  err.Error(),
	}).Info("bad request error")

	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

// NotFoundResponse logs the error and sends a 404 Not Found response
func NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"error":  err.Error(),
	}).Info("not found error")

	WriteJSONError(w, http.StatusNotFound, "resource not found")
}
