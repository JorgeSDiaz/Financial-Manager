package create_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/financial-manager/api/cmd/api/handlers/category/create"
	"github.com/financial-manager/api/cmd/api/handlers/category/response"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	category := buildDomainCategory("cat-1", "Food")
	categoryResp := response.ToCategory(category)

	tests := []struct {
		name       string
		body       any
		uc         *fakeUseCase
		wantStatus int
		wantBody   any
	}{
		{
			name: "valid body returns 201 with created category",
			body: map[string]any{
				"name": "Food", "type": "expense",
				"color": "#FF5733", "icon": "restaurant",
			},
			uc:         &fakeUseCase{out: category},
			wantStatus: http.StatusCreated,
			wantBody:   categoryResp,
		},
		{
			name:       "invalid JSON body returns 400",
			body:       "not-json",
			uc:         &fakeUseCase{},
			wantStatus: http.StatusBadRequest,
			wantBody:   response.Error{Error: "invalid request body"},
		},
		{
			name:       "use case validation error returns 400",
			body:       map[string]any{"name": "", "type": "expense", "color": "red", "icon": "icon"},
			uc:         &fakeUseCase{err: errors.New("category name is required")},
			wantStatus: http.StatusBadRequest,
			wantBody:   response.Error{Error: "category name is required"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := create.New(tc.uc)

			bodyBytes, _ := json.Marshal(tc.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/categories", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
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
