package handlers

// import (
// 	"net/http"

// 	"github.com/Businge931/sba-crud-ops/internal/core/domain"
// 	"github.com/Businge931/sba-crud-ops/internal/core/ports"
// )

// type OddsHandler struct {
// 	oddsService ports.OddsService
// }

// func NewOddsHandler(oddsService ports.OddsService) *OddsHandler {
// 	return &OddsHandler{
// 		oddsService: oddsService,
// 	}
// }

// type ErrorResponse struct {
// 	ErrorDetails string `json:"details"`
// }

// func (h *OddsHandler) CreateOddsHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var req domain.CreateOddsRequest
// 	if err := domain.ReadJSON(w, r, &req); err != nil {
// 		domain.BadRequestResponse(w, r, err)
// 		return
// 	}

// 	if err := h.oddsService.CreateOdds(r.Context(), req); err != nil {
// 		switch err {
// 		case domain.ErrInvalidLeague, domain.ErrEmptyTeams, domain.ErrInvalidOdds,
// 			domain.ErrInvalidStartDate, domain.ErrSameTeams, domain.ErrInvalidOddsProbability:
// 			domain.BadRequestResponse(w, r, err)
// 		default:
// 			domain.InternalServerError(w, r, err)
// 		}
// 		return
// 	}

// 	domain.WriteJSON(w, http.StatusCreated, domain.JSONResponse{
// 		Error:   false,
// 		Message: "odds created successfully",
// 	})
// }

// func (h *OddsHandler) ReadOddsHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var req domain.ReadOddsRequest
// 	if err := domain.ReadJSON(w, r, &req); err != nil {
// 		domain.BadRequestResponse(w, r, err)
// 		return
// 	}

// 	odds, err := h.oddsService.ReadOdds(r.Context(), req)
// 	if err != nil {
// 		switch err {
// 		case domain.ErrInvalidLeague:
// 			domain.BadRequestResponse(w, r, err)
// 		case domain.ErrOddsNotFound:
// 			domain.WriteJSONError(w, http.StatusNotFound, err.Error())
// 		default:
// 			domain.InternalServerError(w, r, err)
// 		}
// 		return
// 	}

// 	domain.WriteJSON(w, http.StatusOK, odds)
// }

// func (h *OddsHandler) UpdateOddsHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPut {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var req domain.CreateOddsRequest
// 	if err := domain.ReadJSON(w, r, &req); err != nil {
// 		domain.BadRequestResponse(w, r, err)
// 		return
// 	}

// 	if err := h.oddsService.UpdateOdds(r.Context(), req); err != nil {
// 		switch err {
// 		case domain.ErrInvalidLeague, domain.ErrEmptyTeams, domain.ErrInvalidOdds,
// 			domain.ErrInvalidStartDate, domain.ErrSameTeams, domain.ErrInvalidOddsProbability:
// 			domain.BadRequestResponse(w, r, err)
// 		case domain.ErrOddsNotFound:
// 			domain.WriteJSONError(w, http.StatusNotFound, err.Error())
// 		default:
// 			domain.InternalServerError(w, r, err)
// 		}
// 		return
// 	}

// 	domain.WriteJSON(w, http.StatusOK, domain.JSONResponse{
// 		Error:   false,
// 		Message: "odds updated successfully",
// 	})
// }

// func (h *OddsHandler) DeleteOddsHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodDelete {
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var req domain.DeleteOddsRequest
// 	if err := domain.ReadJSON(w, r, &req); err != nil {
// 		domain.BadRequestResponse(w, r, err)
// 		return
// 	}

// 	if err := h.oddsService.DeleteOdds(r.Context(), req); err != nil {
// 		switch err {
// 		case domain.ErrInvalidLeague:
// 			domain.BadRequestResponse(w, r, err)
// 		case domain.ErrOddsNotFound:
// 			domain.WriteJSONError(w, http.StatusNotFound, err.Error())
// 		default:
// 			domain.InternalServerError(w, r, err)
// 		}
// 		return
// 	}

// 	domain.WriteJSON(w, http.StatusOK, domain.JSONResponse{
// 		Error:   false,
// 		Message: "odds deleted successfully",
// 	})
// }
