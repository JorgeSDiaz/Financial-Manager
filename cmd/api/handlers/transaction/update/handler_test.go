package update_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/financial-manager/api/cmd/api/handlers/transaction/response"
	"github.com/financial-manager/api/cmd/api/handlers/transaction/update"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	tx := buildDomainTransaction("tx-1", "acc-001")
	txResp := response.ToTransaction(tx)

	tests := []struct {
		name       string
		id         string
		body       any
		uc         *fakeUseCase
		wantStatus int
		wantBody   any
	}{
		{
			name: "valid update returns 200 with updated transaction",
			id:   "tx-1",
			body: map[string]any{
				"description": "Updated description",
				"amount":      100.0,
			},
			uc:         &fakeUseCase{out: tx},
			wantStatus: http.StatusOK,
			wantBody:   txResp,
		},
		{
			name:       "missing id returns 400",
			id:         "",
			body:       map[string]any{"description": "Updated"},
			uc:         &fakeUseCase{},
			wantStatus: http.StatusBadRequest,
			wantBody:   response.Error{Error: "transaction id is required"},
		},
		{
			name:       "invalid JSON body returns 400",
			id:         "tx-1",
			body:       "not-json",
			uc:         &fakeUseCase{},
			wantStatus: http.StatusBadRequest,
			wantBody:   response.Error{Error: "invalid request body"},
		},
		{
			name:       "nonexistent transaction returns 404",
			id:         "missing",
			body:       map[string]any{"description": "Updated"},
			uc:         &fakeUseCase{err: domainshared.ErrNotFound},
			wantStatus: http.StatusNotFound,
			wantBody:   response.Error{Error: "transaction not found"},
		},
		{
			name:       "validation error returns 400",
			id:         "tx-1",
			body:       map[string]any{"description": ""},
			uc:         &fakeUseCase{err: errors.New("invalid date format")},
			wantStatus: http.StatusBadRequest,
			wantBody:   response.Error{Error: "invalid date format"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := update.New(tc.uc)

			bodyBytes, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/transactions/"+tc.id, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Set up chi router context with URL parameter
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rec := httptest.NewRecorder()

			h.Handle(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Equal(t, tc.wantBody, decodeAs(t, rec, tc.wantBody))
		})
	}
}

// decodeAs decodes the recorder body into a new value of the same type as want.
func decodeAs(t *testing.T, rec *httptest.ResponseRecorder, want any) any {
	t.Helper()
	ptr := reflect.New(reflect.TypeOf(want))
	require.NoError(t, json.NewDecoder(rec.Body).Decode(ptr.Interface()))
	return ptr.Elem().Interface()
}
