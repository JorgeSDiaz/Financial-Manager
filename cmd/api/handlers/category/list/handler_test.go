package list_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/financial-manager/api/cmd/api/handlers/category/list"
	"github.com/financial-manager/api/cmd/api/handlers/category/response"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	categories := buildDomainCategories()
	categoriesResp := []response.CategoryResponse{
		response.ToCategory(categories[0]),
		response.ToCategory(categories[1]),
	}

	tests := []struct {
		name       string
		query      string
		uc         *fakeUseCase
		wantStatus int
		wantBody   any
	}{
		{
			name:       "list all categories returns 200 with categories",
			query:      "",
			uc:         &fakeUseCase{out: categories},
			wantStatus: http.StatusOK,
			wantBody:   categoriesResp,
		},
		{
			name:       "list with type filter returns 200 with filtered categories",
			query:      "?type=expense",
			uc:         &fakeUseCase{out: categories[:1]},
			wantStatus: http.StatusOK,
			wantBody:   []response.CategoryResponse{categoriesResp[0]},
		},
		{
			name:       "use case error returns 400",
			query:      "",
			uc:         buildFailingUseCase(errors.New("db error")),
			wantStatus: http.StatusBadRequest,
			wantBody:   response.Error{Error: "db error"},
		},
		{
			name:       "empty list returns 200 with empty array",
			query:      "",
			uc:         &fakeUseCase{out: []domaincategory.Category{}},
			wantStatus: http.StatusOK,
			wantBody:   []response.CategoryResponse{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := list.New(tc.uc)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/categories"+tc.query, nil)
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
