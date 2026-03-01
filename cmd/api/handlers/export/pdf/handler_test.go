// Package pdf_test contains tests for the PDF export handler.
package pdf_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/cmd/api/handlers/export/pdf"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		uc         *fakeUseCase
		body       string
		wantStatus int
		wantHeader string
	}{
		{
			name:       "generates PDF successfully",
			uc:         &fakeUseCase{pdf: []byte("%PDF-1.4 test")},
			body:       `{"month":"2026-02"}`,
			wantStatus: http.StatusOK,
			wantHeader: "application/pdf",
		},
		{
			name:       "empty body uses current month",
			uc:         &fakeUseCase{pdf: []byte("%PDF-1.4 test")},
			body:       "",
			wantStatus: http.StatusOK,
			wantHeader: "application/pdf",
		},
		{
			name:       "invalid month returns 400",
			uc:         &fakeUseCase{err: errors.New("invalid month format: parsing time")},
			body:       `{"month":"invalid"}`,
			wantStatus: http.StatusBadRequest,
			wantHeader: "",
		},
		{
			name:       "use case error returns 500",
			uc:         &fakeUseCase{err: errors.New("db error")},
			body:       `{"month":"2026-02"}`,
			wantStatus: http.StatusInternalServerError,
			wantHeader: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := pdf.New(tc.uc)

			var req *http.Request
			if tc.body != "" {
				req = httptest.NewRequest(http.MethodPost, "/api/v1/export/pdf", strings.NewReader(tc.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(http.MethodPost, "/api/v1/export/pdf", nil)
			}
			rec := httptest.NewRecorder()

			h.Handle(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantHeader != "" {
				assert.Equal(t, tc.wantHeader, rec.Header().Get("Content-Type"))
				assert.Contains(t, rec.Header().Get("Content-Disposition"), "attachment")
			}
		})
	}
}
