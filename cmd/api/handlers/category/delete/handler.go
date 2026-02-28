// Package delete handles DELETE /api/v1/categories/{id}.
package delete

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/financial-manager/api/cmd/api/handlers/category/response"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

type useCase interface {
	Execute(ctx context.Context, id string) error
}

// Handler handles DELETE /api/v1/categories/{id}.
type Handler struct {
	uc useCase
}

// New creates a Handler with its required use case dependency.
func New(uc useCase) *Handler {
	return &Handler{uc: uc}
}

// Handle processes DELETE /api/v1/categories/{id} and returns 204 on success.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.uc.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainshared.ErrNotFound) {
			response.WriteError(w, http.StatusNotFound, "category not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
