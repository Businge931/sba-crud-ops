package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Businge931/api-gateway/internal/core/domain"
	"github.com/Businge931/api-gateway/internal/core/ports"
)

type OddsHandler struct {
	oddsService ports.OddsService
}

func NewOddsHandler(oddsService ports.OddsService) *OddsHandler {
	return &OddsHandler{
		oddsService: oddsService,
	}
}

type ErrorResponse struct {
	Details string `json:"details"`
}

func (h *OddsHandler) CreateOddsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req domain.CreateOddsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Details: "invalid request format"})
		return
	}

	if err := h.oddsService.CreateOdds(r.Context(), req); err != nil {
		switch err {
		case domain.ErrInvalidLeague, domain.ErrEmptyTeams, domain.ErrInvalidOdds,
			domain.ErrInvalidStartDate, domain.ErrSameTeams, domain.ErrInvalidOddsProbability:
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{Details: err.Error()})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Details: "internal server error"})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OddsHandler) ReadOddsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req domain.ReadOddsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Details: "invalid request format"})
		return
	}

	odds, err := h.oddsService.ReadOdds(r.Context(), req)
	if err != nil {
		switch err {
		case domain.ErrInvalidLeague, domain.ErrInvalidDateFormat:
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{Details: err.Error()})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Details: "internal server error"})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(odds)
}

func (h *OddsHandler) UpdateOddsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req domain.CreateOddsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Details: "invalid request format"})
		return
	}

	if err := h.oddsService.UpdateOdds(r.Context(), req); err != nil {
		switch err {
		case domain.ErrInvalidLeague, domain.ErrEmptyTeams, domain.ErrInvalidOdds,
			domain.ErrInvalidStartDate, domain.ErrSameTeams, domain.ErrInvalidOddsProbability:
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{Details: err.Error()})
		case domain.ErrOddsNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Details: err.Error()})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Details: "internal server error"})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OddsHandler) DeleteOddsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req domain.DeleteOddsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Details: "invalid request format"})
		return
	}

	if err := h.oddsService.DeleteOdds(r.Context(), req); err != nil {
		switch err {
		case domain.ErrInvalidLeague:
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{Details: err.Error()})
		case domain.ErrOddsNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Details: err.Error()})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Details: "internal server error"})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
