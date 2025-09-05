package config

import (
	"sync"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
)

// LeagueRegistry manages the available leagues in the system
type LeagueRegistry interface {
	// IsSupported checks if a league is supported
	IsSupported(league string) bool

	// GetSupportedLeagues returns all supported leagues
	GetSupportedLeagues() []string

	// RegisterLeague adds a new supported league
	RegisterLeague(league string) error

	// UnregisterLeague removes a league from supported leagues
	UnregisterLeague(league string) error
}

// DefaultLeagueRegistry is the default implementation of LeagueRegistry
type DefaultLeagueRegistry struct {
	supportedLeagues map[string]struct{}
	mu               sync.RWMutex
}

// NewLeagueRegistry creates a new league registry with initial leagues
func NewLeagueRegistry(initialLeagues []string) *DefaultLeagueRegistry {
	registry := &DefaultLeagueRegistry{
		supportedLeagues: make(map[string]struct{}),
	}

	// Add default if nothing provided
	if len(initialLeagues) == 0 {
		registry.supportedLeagues["English Premier League"] = struct{}{}
	} else {
		for _, league := range initialLeagues {
			registry.supportedLeagues[league] = struct{}{}
		}
	}

	return registry
}

// IsSupported checks if a league is supported
func (r *DefaultLeagueRegistry) IsSupported(league string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.supportedLeagues[league]
	return exists
}

// GetSupportedLeagues returns all supported leagues
func (r *DefaultLeagueRegistry) GetSupportedLeagues() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	leagues := make([]string, 0, len(r.supportedLeagues))
	for league := range r.supportedLeagues {
		leagues = append(leagues, league)
	}

	return leagues
}

// RegisterLeague adds a new supported league
func (r *DefaultLeagueRegistry) RegisterLeague(league string) error {
	if league == "" {
		return domain.ErrEmptyLeagueName
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.supportedLeagues[league] = struct{}{}
	return nil
}

// UnregisterLeague removes a league from supported leagues
func (r *DefaultLeagueRegistry) UnregisterLeague(league string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.supportedLeagues[league]; !exists {
		return domain.ErrLeagueNotFound
	}

	delete(r.supportedLeagues, league)
	return nil
}
