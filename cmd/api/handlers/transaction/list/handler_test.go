package list_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/cmd/api/handlers/transaction/list"
	"github.com/financial-manager/api/cmd/api/handlers/transaction/response"
	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

func TestHandler_HandleIncomes(t *testing.T) {
	t.Parallel()

	income1 := buildDomainTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeIncome)
	income2 := buildDomainTransaction("tx-2", "acc-002", domaintransaction.TransactionTypeIncome)

	tests := []struct {
		name       string
		uc         *fakeIncomeLister
		wantStatus int
		wantBody   map[string]any
	}{
		{
			name:       "returns income transactions",
			uc:         &fakeIncomeLister{out: []domaintransaction.Transaction{income1, income2}},
			wantStatus: http.StatusOK,
			wantBody: map[string]any{
				"transactions": []response.Transaction{
					response.ToTransaction(income1),
					response.ToTransaction(income2),
				},
			},
		},
		{
			name:       "empty list returns empty array",
			uc:         &fakeIncomeLister{out: nil},
			wantStatus: http.StatusOK,
			wantBody: map[string]any{
				"transactions": []response.Transaction{},
			},
		},
		{
			name:       "use case error returns 500",
			uc:         &fakeIncomeLister{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
			wantBody:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := list.New(tc.uc, &fakeExpenseLister{})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/incomes", nil)
			rec := httptest.NewRecorder()

			h.HandleIncomes(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}

func TestHandler_HandleExpenses(t *testing.T) {
	t.Parallel()

	expense1 := buildDomainTransaction("tx-1", "acc-001", domaintransaction.TransactionTypeExpense)
	expense2 := buildDomainTransaction("tx-2", "acc-002", domaintransaction.TransactionTypeExpense)

	tests := []struct {
		name       string
		uc         *fakeExpenseLister
		wantStatus int
		wantBody   map[string]any
	}{
		{
			name:       "returns expense transactions",
			uc:         &fakeExpenseLister{out: []domaintransaction.Transaction{expense1, expense2}},
			wantStatus: http.StatusOK,
			wantBody: map[string]any{
				"transactions": []response.Transaction{
					response.ToTransaction(expense1),
					response.ToTransaction(expense2),
				},
			},
		},
		{
			name:       "empty list returns empty array",
			uc:         &fakeExpenseLister{out: nil},
			wantStatus: http.StatusOK,
			wantBody: map[string]any{
				"transactions": []response.Transaction{},
			},
		},
		{
			name:       "use case error returns 500",
			uc:         &fakeExpenseLister{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
			wantBody:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := list.New(&fakeIncomeLister{}, tc.uc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/expenses", nil)
			rec := httptest.NewRecorder()

			h.HandleExpenses(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}
