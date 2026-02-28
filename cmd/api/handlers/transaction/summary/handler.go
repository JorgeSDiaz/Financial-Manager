// Package summary handles GET /api/v1/transactions/summary.
package summary

import (
	"context"
	"net/http"

	"github.com/financial-manager/api/cmd/api/handlers/transaction/response"
	appSummary "github.com/financial-manager/api/internal/application/transaction/summary"
)

type useCase interface {
	Execute(ctx context.Context, in appSummary.Input) (appSummary.Summary, error)
}

// Handler handles GET /api/v1/transactions/summary.
type Handler struct {
	uc useCase
}

// New creates a Handler with its required use case dependency.
func New(uc useCase) *Handler {
	return &Handler{uc: uc}
}

// Handle processes GET /api/v1/transactions/summary.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	input := appSummary.Input{
		AccountID: r.URL.Query().Get("account_id"),
		StartDate: r.URL.Query().Get("start_date"),
		EndDate:   r.URL.Query().Get("end_date"),
	}

	sum, err := h.uc.Execute(r.Context(), input)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, response.Summary{
		TotalIncome:  sum.TotalIncome,
		TotalExpense: sum.TotalExpense,
		Balance:      sum.Balance,
	})
}
