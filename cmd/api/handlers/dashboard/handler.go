// Package dashboard handles GET /api/v1/dashboard.
package dashboard

import (
	"context"
	"encoding/json"
	"net/http"

	appDashboard "github.com/financial-manager/api/internal/application/dashboard"
)

// Response represents the dashboard JSON response.
type Response struct {
	GlobalBalance      float64             `json:"global_balance"`
	MonthlySummary     MonthlySummary      `json:"monthly_summary"`
	ExpensesByCategory []ExpenseByCategory `json:"expenses_by_category"`
	RecentTransactions []RecentTransaction `json:"recent_transactions"`
}

// MonthlySummary represents the monthly financial summary.
type MonthlySummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetBalance   float64 `json:"net_balance"`
}

// ExpenseByCategory represents expense breakdown by category.
type ExpenseByCategory struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Total        float64 `json:"total"`
	Percentage   float64 `json:"percentage"`
}

// RecentTransaction represents a recent transaction.
type RecentTransaction struct {
	ID           string  `json:"id"`
	Type         string  `json:"type"`
	Amount       float64 `json:"amount"`
	Date         string  `json:"date"`
	Description  string  `json:"description"`
	CategoryName string  `json:"category_name"`
}

type useCase interface {
	Execute(ctx context.Context) (appDashboard.Output, error)
}

// Handler handles GET /api/v1/dashboard.
type Handler struct {
	uc useCase
}

// New creates a Handler with its required use case dependency.
func New(uc useCase) *Handler {
	return &Handler{uc: uc}
}

// Handle processes GET /api/v1/dashboard.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	out, err := h.uc.Execute(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Transform to response
	expensesByCategory := make([]ExpenseByCategory, len(out.ExpensesByCategory))
	for i, e := range out.ExpensesByCategory {
		expensesByCategory[i] = ExpenseByCategory{
			CategoryID:   e.CategoryID,
			CategoryName: e.CategoryName,
			Total:        e.Total,
			Percentage:   e.Percentage,
		}
	}

	recentTransactions := make([]RecentTransaction, len(out.RecentTransactions))
	for i, t := range out.RecentTransactions {
		recentTransactions[i] = RecentTransaction{
			ID:           t.ID,
			Type:         t.Type,
			Amount:       t.Amount,
			Date:         t.Date,
			Description:  t.Description,
			CategoryName: t.CategoryName,
		}
	}

	resp := Response{
		GlobalBalance: out.GlobalBalance,
		MonthlySummary: MonthlySummary{
			TotalIncome:  out.MonthlySummary.TotalIncome,
			TotalExpense: out.MonthlySummary.TotalExpense,
			NetBalance:   out.MonthlySummary.NetBalance,
		},
		ExpensesByCategory: expensesByCategory,
		RecentTransactions: recentTransactions,
	}

	writeJSON(w, http.StatusOK, resp)
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError writes an error response.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
