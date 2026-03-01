// Package export_test contains tests for the export handlers.
package export_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/cmd/api/handlers/export"
)

func TestHandler_HandleCSV(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		csvUC      *fakeCSVUseCase
		wantStatus int
		wantHeader string
	}{
		{
			name:       "exports CSV successfully",
			csvUC:      &fakeCSVUseCase{csv: "date,type,amount\n2026-02-28,income,100.00\n"},
			wantStatus: http.StatusOK,
			wantHeader: "text/csv",
		},
		{
			name:       "use case error returns 500",
			csvUC:      &fakeCSVUseCase{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
			wantHeader: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := export.New(tc.csvUC, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/export/csv", nil)
			rec := httptest.NewRecorder()

			h.HandleCSV(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantHeader != "" {
				assert.Equal(t, tc.wantHeader, rec.Header().Get("Content-Type"))
				assert.Contains(t, rec.Header().Get("Content-Disposition"), "attachment")
			}
		})
	}
}

func TestHandler_HandleJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		jsonUC     *fakeJSONUseCase
		wantStatus int
		wantHeader string
	}{
		{
			name:       "exports JSON successfully",
			jsonUC:     &fakeJSONUseCase{json: []byte(`{"accounts":[],"categories":[],"transactions":[]}`)},
			wantStatus: http.StatusOK,
			wantHeader: "application/json",
		},
		{
			name:       "use case error returns 500",
			jsonUC:     &fakeJSONUseCase{err: errors.New("db error")},
			wantStatus: http.StatusInternalServerError,
			wantHeader: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := export.New(nil, tc.jsonUC)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/export/json", nil)
			rec := httptest.NewRecorder()

			h.HandleJSON(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantHeader != "" {
				assert.Equal(t, tc.wantHeader, rec.Header().Get("Content-Type"))
				assert.Contains(t, rec.Header().Get("Content-Disposition"), "attachment")
			}
		})
	}
}
