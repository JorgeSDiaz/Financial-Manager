// Package update handles PUT /api/v1/accounts/{id}.
package update

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/financial-manager/api/cmd/api/handlers/account/response"
	appUpdate "github.com/financial-manager/api/internal/application/account/update"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

type useCase interface {
	Execute(ctx context.Context, in appUpdate.Input) (domainaccount.Account, error)
}

// Handler handles PUT /api/v1/accounts/{id}.
type Handler struct {
	uc useCase
}

// New creates a Handler with its required use case dependency.
func New(uc useCase) *Handler {
	return &Handler{uc: uc}
}

type updateRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

// Handle processes PUT /api/v1/accounts/{id} and returns 200 with the updated account.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	acc, err := h.uc.Execute(r.Context(), appUpdate.Input{
		ID:    id,
		Name:  req.Name,
		Color: req.Color,
		Icon:  req.Icon,
	})
	if err != nil {
		if errors.Is(err, domainshared.ErrNotFound) {
			response.WriteError(w, http.StatusNotFound, "account not found")
			return
		}
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusOK, response.ToAccount(acc))
}
