// Package list handles GET /api/v1/accounts.
package list

import (
	"context"
	"net/http"

	"github.com/financial-manager/api/cmd/api/handlers/account/response"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

type lister interface {
	Execute(ctx context.Context) ([]domainaccount.Account, error)
}

type balanceGetter interface {
	Execute(ctx context.Context) (float64, error)
}

// Handler handles GET /api/v1/accounts.
type Handler struct {
	lister        lister
	balanceGetter balanceGetter
}

// New creates a Handler with its required use case dependencies.
func New(lister lister, balanceGetter balanceGetter) *Handler {
	return &Handler{lister: lister, balanceGetter: balanceGetter}
}

type listResponse struct {
	Accounts      []response.Account `json:"accounts"`
	GlobalBalance float64            `json:"global_balance"`
}

// Handle processes GET /api/v1/accounts and returns 200 with all accounts and the global balance.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.lister.Execute(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	total, err := h.balanceGetter.Execute(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]response.Account, 0, len(accounts))
	for _, a := range accounts {
		resp = append(resp, response.ToAccount(a))
	}

	response.WriteJSON(w, http.StatusOK, listResponse{
		Accounts:      resp,
		GlobalBalance: total,
	})
}
