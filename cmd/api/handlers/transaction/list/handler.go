// Package list handles GET /api/v1/transactions/incomes and GET /api/v1/transactions/expenses.
package list

import (
	"context"
	"net/http"

	"github.com/financial-manager/api/cmd/api/handlers/transaction/response"
	expenselist "github.com/financial-manager/api/internal/application/transaction/expense/list"
	incomelist "github.com/financial-manager/api/internal/application/transaction/income/list"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

type incomeLister interface {
	Execute(ctx context.Context, in incomelist.Input) ([]domaintransaction.Transaction, error)
}

type expenseLister interface {
	Execute(ctx context.Context, in expenselist.Input) ([]domaintransaction.Transaction, error)
}

// Handler handles GET /api/v1/transactions/incomes and GET /api/v1/transactions/expenses.
type Handler struct {
	incomeUC  incomeLister
	expenseUC expenseLister
}

// New creates a Handler with its required use case dependencies.
func New(incomeUC incomeLister, expenseUC expenseLister) *Handler {
	return &Handler{incomeUC: incomeUC, expenseUC: expenseUC}
}

// HandleIncomes processes GET /api/v1/transactions/incomes.
func (h *Handler) HandleIncomes(w http.ResponseWriter, r *http.Request) {
	in := incomelist.Input{
		AccountID:  r.URL.Query().Get("account_id"),
		CategoryID: r.URL.Query().Get("category_id"),
		StartDate:  r.URL.Query().Get("start_date"),
		EndDate:    r.URL.Query().Get("end_date"),
	}

	txs, err := h.incomeUC.Execute(r.Context(), in)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]response.Transaction, 0, len(txs))
	for _, tx := range txs {
		resp = append(resp, response.ToTransaction(tx))
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"transactions": resp,
	})
}

// HandleExpenses processes GET /api/v1/transactions/expenses.
func (h *Handler) HandleExpenses(w http.ResponseWriter, r *http.Request) {
	in := expenselist.Input{
		AccountID:  r.URL.Query().Get("account_id"),
		CategoryID: r.URL.Query().Get("category_id"),
		StartDate:  r.URL.Query().Get("start_date"),
		EndDate:    r.URL.Query().Get("end_date"),
	}

	txs, err := h.expenseUC.Execute(r.Context(), in)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := make([]response.Transaction, 0, len(txs))
	for _, tx := range txs {
		resp = append(resp, response.ToTransaction(tx))
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"transactions": resp,
	})
}
