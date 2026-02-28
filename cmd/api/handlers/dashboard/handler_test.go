// Package dashboard_test contains tests for the dashboard handler.
package dashboard_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/cmd/api/handlers/dashboard"
	appDashboard "github.com/financial-manager/api/internal/application/dashboard"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		uc         *fakeUseCase
		wantStatus int
	}{
		{
			name: "returns dashboard with all data",
			uc: &fakeUseCase{out: appDashboard.Output{
				GlobalBalance: 2500.00,
				MonthlySummary: appDashboard.MonthlySummary{
					TotalIncome:  3000.00,
					TotalExpense: 500.00,
					NetBalance:   2500.00,
				},
				ExpensesByCategory: []appDashboard.ExpenseByCategory{
					{CategoryID: "cat-1", CategoryName: "Alimentación", Total: 300.00, Percentage: 60.0},
					{CategoryID: "cat-2", CategoryName: "Transporte", Total: 200.00, Percentage: 40.0},
				},
				RecentTransactions: []appDashboard.RecentTransaction{
					{ID: "tx-1", Type: "income", Amount: 1000.00, Date: "2026-02-28", Description: "Salary", CategoryName: "Income"},
					{ID: "tx-2", Type: "expense", Amount: 50.00, Date: "2026-02-27", Description: "Groceries", CategoryName: "Alimentación"},
				},
			}},
			wantStatus: http.StatusOK,
		},
		{
			name: "returns empty dashboard when no data",
			uc: &fakeUseCase{out: appDashboard.Output{
				GlobalBalance:      0,
				MonthlySummary:     appDashboard.MonthlySummary{},
				ExpensesByCategory: []appDashboard.ExpenseByCategory{},
				RecentTransactions: []appDashboard.RecentTransaction{},
			}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "use case error returns 500",
			uc:         &fakeUseCase{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := dashboard.New(tc.uc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/dashboard", nil)
			rec := httptest.NewRecorder()

			h.Handle(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}
