// Package create handles POST /api/v1/transactions/expenses.
package create

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/financial-manager/api/cmd/api/handlers/transaction/response"
	appCreate "github.com/financial-manager/api/internal/application/transaction/expense/create"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

type useCase interface {
	Execute(ctx context.Context, in appCreate.Input) (domaintransaction.Transaction, error)
}

// Handler handles POST /api/v1/transactions/expenses.
type Handler struct {
	uc useCase
}

// New creates a Handler with its required use case dependency.
func New(uc useCase) *Handler {
	return &Handler{uc: uc}
}

type createRequest struct {
	AccountID   string  `json:"account_id"`
	CategoryID  string  `json:"category_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

// Handle processes POST /api/v1/transactions/expenses and returns 201 with the created transaction.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var req createRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tx, err := h.uc.Execute(r.Context(), appCreate.Input{
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
	})
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, response.ToTransaction(tx))
}
