package delete_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	accountdelete "github.com/financial-manager/api/cmd/api/handlers/account/delete"
	domainaccount "github.com/financial-manager/api/internal/domain/account"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		id         string
		uc         *fakeUseCase
		wantStatus int
	}{
		{
			name:       "existing account returns 204",
			id:         "acc-1",
			uc:         &fakeUseCase{},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "nonexistent account returns 404",
			id:         "missing",
			uc:         &fakeUseCase{err: domainshared.ErrNotFound},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "account with transactions returns 409",
			id:         "acc-2",
			uc:         &fakeUseCase{err: domainaccount.ErrAccountHasTransactions},
			wantStatus: http.StatusConflict,
		},
		{
			name:       "repository error returns 500",
			id:         "acc-1",
			uc:         &fakeUseCase{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := accountdelete.New(tc.uc)

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/accounts/"+tc.id, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rec := httptest.NewRecorder()

			h.Handle(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}
