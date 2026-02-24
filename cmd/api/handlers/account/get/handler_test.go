package get_test

import (
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

	"github.com/financial-manager/api/cmd/api/handlers/account/get"
	"github.com/financial-manager/api/cmd/api/handlers/account/response"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	account := buildDomainAccount("acc-1", "Cash")
	accountResp := response.ToAccount(account)

	tests := []struct {
		name       string
		id         string
		uc         *fakeUseCase
		wantStatus int
		wantBody   any
	}{
		{
			name:       "existing account returns 200",
			id:         "acc-1",
			uc:         &fakeUseCase{out: account},
			wantStatus: http.StatusOK,
			wantBody:   accountResp,
		},
		{
			name:       "nonexistent ID returns 404",
			id:         "missing",
			uc:         &fakeUseCase{err: domainshared.ErrNotFound},
			wantStatus: http.StatusNotFound,
			wantBody:   response.Error{Error: "account not found"},
		},
		{
			name:       "repository error returns 500",
			id:         "acc-1",
			uc:         &fakeUseCase{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
			wantBody:   response.Error{Error: "internal server error"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := get.New(tc.uc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/accounts/"+tc.id, nil)
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
