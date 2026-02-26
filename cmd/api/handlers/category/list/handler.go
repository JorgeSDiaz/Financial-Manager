// Package list handles GET /api/v1/categories.
package list

import (
	"context"
	"net/http"

	"github.com/financial-manager/api/cmd/api/handlers/category/response"
	appList "github.com/financial-manager/api/internal/application/category/list"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

type useCase interface {
	Execute(ctx context.Context, in appList.Input) ([]domaincategory.Category, error)
}

// Handler handles GET /api/v1/categories.
type Handler struct {
	uc useCase
}

// New creates a Handler with its required use case dependency.
func New(uc useCase) *Handler {
	return &Handler{uc: uc}
}

// Handle processes GET /api/v1/categories and returns all active categories.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	categoryType := r.URL.Query().Get("type")

	var input appList.Input
	if categoryType != "" {
		input.Type = &categoryType
	}

	categories, err := h.uc.Execute(r.Context(), input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := make([]response.CategoryResponse, len(categories))
	for i, c := range categories {
		resp[i] = response.ToCategory(c)
	}

	response.WriteJSON(w, http.StatusOK, resp)
}
