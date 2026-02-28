package delete_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/cmd/api/handlers/transaction/delete"
	"github.com/financial-manager/api/cmd/api/handlers/transaction/response"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		id         string
		uc         *fakeUseCase
		wantStatus int
		wantBody   any
	}{
		{
			name:       "existing transaction returns 204",
			id:         "tx-1",
			uc:         &fakeUseCase{err: nil},
			wantStatus: http.StatusNoContent,
			wantBody:   nil,
		},
		{
			name:       "missing id returns 400",
			id:         "",
			uc:         &fakeUseCase{},
			wantStatus: http.StatusBadRequest,
			wantBody:   response.Error{Error: "transaction id is required"},
		},
		{
			name:       "nonexistent transaction returns 404",
			id:         "missing",
			uc:         &fakeUseCase{err: domainshared.ErrNotFound},
			wantStatus: http.StatusNotFound,
			wantBody:   response.Error{Error: "transaction not found"},
		},
		{
			name:       "repository error returns 500",
			id:         "any",
			uc:         &fakeUseCase{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
			wantBody:   response.Error{Error: "internal server error"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := delete.New(tc.uc)

			req := httptest.NewRequest(http.MethodDelete, "/api/v1/transactions/"+tc.id, nil)

			// Set up chi router context with URL parameter
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rec := httptest.NewRecorder()

			h.Handle(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantBody != nil {
				var body response.Error
				json.NewDecoder(rec.Body).Decode(&body)
				assert.Equal(t, tc.wantBody, body)
			}
		})
	}
}
