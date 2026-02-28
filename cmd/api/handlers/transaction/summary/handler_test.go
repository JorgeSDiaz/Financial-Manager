package summary_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/cmd/api/handlers/transaction/response"
	"github.com/financial-manager/api/cmd/api/handlers/transaction/summary"
	appsummary "github.com/financial-manager/api/internal/application/transaction/summary"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		uc         *fakeUseCase
		wantStatus int
		wantBody   response.Summary
	}{
		{
			name: "returns summary with calculations",
			uc: &fakeUseCase{out: appsummary.Summary{
				TotalIncome:  1000.0,
				TotalExpense: 500.0,
				Balance:      500.0,
			}},
			wantStatus: http.StatusOK,
			wantBody: response.Summary{
				TotalIncome:  1000.0,
				TotalExpense: 500.0,
				Balance:      500.0,
			},
		},
		{
			name: "returns zero summary when no transactions",
			uc: &fakeUseCase{out: appsummary.Summary{
				TotalIncome:  0,
				TotalExpense: 0,
				Balance:      0,
			}},
			wantStatus: http.StatusOK,
			wantBody: response.Summary{
				TotalIncome:  0,
				TotalExpense: 0,
				Balance:      0,
			},
		},
		{
			name:       "use case error returns 500",
			uc:         &fakeUseCase{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
			wantBody:   response.Summary{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := summary.New(tc.uc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/summary", nil)
			rec := httptest.NewRecorder()

			h.Handle(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}
