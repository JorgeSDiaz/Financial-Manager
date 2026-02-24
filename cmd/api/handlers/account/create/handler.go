// Package create handles POST /api/v1/accounts.
package create

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/financial-manager/api/cmd/api/handlers/account/response"
	appCreate "github.com/financial-manager/api/internal/application/account/create"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

type useCase interface {
	Execute(ctx context.Context, in appCreate.Input) (domainaccount.Account, error)
}

// Handler handles POST /api/v1/accounts.
type Handler struct {
	uc useCase
}

// New creates a Handler with its required use case dependency.
func New(uc useCase) *Handler {
	return &Handler{uc: uc}
}

type createRequest struct {
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	InitialBalance float64 `json:"initial_balance"`
	Currency       string  `json:"currency"`
	Color          string  `json:"color"`
	Icon           string  `json:"icon"`
}

// Handle processes POST /api/v1/accounts and returns 201 with the created account.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var req createRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	acc, err := h.uc.Execute(r.Context(), appCreate.Input{
		Name:           req.Name,
		Type:           req.Type,
		InitialBalance: req.InitialBalance,
		Currency:       req.Currency,
		Color:          req.Color,
		Icon:           req.Icon,
	})
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.ToAccount(acc))
}
