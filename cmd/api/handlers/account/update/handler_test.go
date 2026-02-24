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

	"github.com/financial-manager/api/cmd/api/handlers/account/response"
	"github.com/financial-manager/api/cmd/api/handlers/account/update"
	appUpdate "github.com/financial-manager/api/internal/application/account/update"
	domainshared "github.com/financial-manager/api/internal/domain/shared"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	updatedAccount := buildDomainAccount("acc-1", "Updated")
	updatedResp := response.ToAccount(updatedAccount)

	tests := []struct {
		name       string
		id         string
		body       any
		uc         *fakeUseCase
		wantStatus int
		wantBody   any
	}{
		{
			name: "valid update returns 200 with updated account",
			id:   "acc-1",
			body: map[string]any{"name": "Updated", "color": "#FF0000"},
			uc: &fakeUseCase{
				wantInput: appUpdate.Input{ID: "acc-1", Name: "Updated", Color: "#FF0000"},
				out:       updatedAccount,
			},
			wantStatus: http.StatusOK,
			wantBody:   updatedResp,
		},
		{
			name:       "invalid JSON returns 400",
			id:         "acc-1",
			body:       "not-json",
			uc:         &fakeUseCase{},
			wantStatus: http.StatusBadRequest,
			wantBody:   response.Error{Error: "invalid request body"},
		},
		{
			name:       "nonexistent account returns 404",
			id:         "missing",
			body:       map[string]any{"name": "X"},
			uc:         &fakeUseCase{err: domainshared.ErrNotFound},
			wantStatus: http.StatusNotFound,
			wantBody:   response.Error{Error: "account not found"},
		},
		{
			name:       "validation error returns 400",
			id:         "acc-1",
			body:       map[string]any{"name": ""},
			uc:         &fakeUseCase{err: errors.New("account name is required")},
			wantStatus: http.StatusBadRequest,
			wantBody:   response.Error{Error: "account name is required"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := update.New(tc.uc)

			bodyBytes, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/accounts/"+tc.id, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
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
