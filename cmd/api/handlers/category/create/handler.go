// Package create handles POST /api/v1/categories.
package create

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/financial-manager/api/cmd/api/handlers/category/response"
	appCreate "github.com/financial-manager/api/internal/application/category/create"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

type useCase interface {
	Execute(ctx context.Context, in appCreate.Input) (domaincategory.Category, error)
}

// Handler handles POST /api/v1/categories.
type Handler struct {
	uc useCase
}

// New creates a Handler with its required use case dependency.
func New(uc useCase) *Handler {
	return &Handler{uc: uc}
}

type createRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

// Handle processes POST /api/v1/categories and returns 201 with the created category.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var req createRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cat, err := h.uc.Execute(r.Context(), appCreate.Input{
		Name:  req.Name,
		Type:  req.Type,
		Color: req.Color,
		Icon:  req.Icon,
	})
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.ToCategory(cat))
}
