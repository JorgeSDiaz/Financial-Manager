// Package balance handles GET /api/v1/accounts/{id}/balance.
package balance

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/financial-manager/api/cmd/api/handlers/account/response"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

type useCase interface {
	Execute(ctx context.Context, id string) (domainaccount.Account, error)
}

// Handler handles GET /api/v1/accounts/{id}/balance.
type Handler struct {
	uc useCase
}

// New creates a Handler with its required use case dependency.
func New(uc useCase) *Handler {
	return &Handler{uc: uc}
}

type balanceResponse struct {
	AccountID      string  `json:"account_id"`
	CurrentBalance float64 `json:"current_balance"`
}

// Handle processes GET /api/v1/accounts/{id}/balance and returns 200 with the current balance.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	acc, err := h.uc.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainshared.ErrNotFound) {
			response.WriteError(w, http.StatusNotFound, "account not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, balanceResponse{
		AccountID:      acc.ID,
		CurrentBalance: acc.CurrentBalance,
	})
}
