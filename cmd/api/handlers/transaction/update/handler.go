// Package update handles PUT /api/v1/transactions/{id}.
package update

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/financial-manager/api/cmd/api/handlers/transaction/response"
	appUpdate "github.com/financial-manager/api/internal/application/transaction/update"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

type useCase interface {
	Execute(ctx context.Context, in appUpdate.Input) (domaintransaction.Transaction, error)
}

// Handler handles PUT /api/v1/transactions/{id}.
type Handler struct {
	uc useCase
}

// New creates a Handler with its required use case dependency.
func New(uc useCase) *Handler {
	return &Handler{uc: uc}
}

type updateRequest struct {
	CategoryID  string  `json:"category_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

// Handle processes PUT /api/v1/transactions/{id}.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.WriteError(w, http.StatusBadRequest, "transaction id is required")
		return
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tx, err := h.uc.Execute(r.Context(), appUpdate.Input{
		ID:          id,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
	})
	if err != nil {
		if errors.Is(err, domainshared.ErrNotFound) {
			response.WriteError(w, http.StatusNotFound, "transaction not found")
			return
		}
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, response.ToTransaction(tx))
}
