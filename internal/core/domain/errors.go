package domain

import (
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	ErrInvalidLeague          = errors.New("invalid league")
	ErrEmptyTeams             = errors.New("team names cannot be empty")
	ErrInvalidOdds            = errors.New("odds must be positive numbers")
	ErrInvalidStartDate       = errors.New("game date must be in the future")
	ErrSameTeams              = errors.New("home and away teams cannot be the same")
	ErrInvalidOddsProbability = errors.New("sum of odds probabilities must be valid")
	ErrInvalidDateFormat      = errors.New("invalid date format")
	ErrOddsNotFound           = errors.New("odds not found")
)

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	WriteJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	WriteJSONError(w, http.StatusBadRequest, err.Error())
}
func NotFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err)

	WriteJSONError(w, http.StatusNotFound, "not found")
}
