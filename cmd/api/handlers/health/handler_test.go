package health_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/financial-manager/api/cmd/api/handlers/health"
	"github.com/financial-manager/api/cmd/api/handlers/health/mocks"
	domainhealth "github.com/financial-manager/api/internal/domain/health"
)

func TestHandler_Check(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		buildChecker  func(t *testing.T) *mocks.Checker
		wantStatus    int
		wantCT        string
		wantBody      map[string]string
		wantTimestamp string
	}{
		{
			name: "returns 200 OK with StatusUp",
			buildChecker: func(t *testing.T) *mocks.Checker {
				return buildMockChecker(t, newHealthResult(domainhealth.StatusUp, "1.0.0"))
			},
			wantStatus: http.StatusOK,
			wantCT:     "application/json",
			wantBody: map[string]string{
				"status":  "up",
				"version": "1.0.0",
			},
			wantTimestamp: "2026-02-23T10:00:00Z",
		},
		{
			name: "returns 200 OK with StatusDown",
			buildChecker: func(t *testing.T) *mocks.Checker {
				return buildMockChecker(t, newHealthResult(domainhealth.StatusDown, "2.0.0"))
			},
			wantStatus: http.StatusOK,
			wantCT:     "application/json",
			wantBody: map[string]string{
				"status":  "down",
				"version": "2.0.0",
			},
			wantTimestamp: "2026-02-23T10:00:00Z",
		},
		{
			name: "returns 503 when checker returns error",
			buildChecker: func(t *testing.T) *mocks.Checker {
				return buildMockCheckerWithError(t, errors.New("dependency down"))
			},
			wantStatus: http.StatusServiceUnavailable,
			wantCT:     "text/plain; charset=utf-8",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			checker := tc.buildChecker(t)
			handler := health.NewHandler(checker)
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			rec := httptest.NewRecorder()

			handler.Check(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Contains(t, rec.Header().Get("Content-Type"), tc.wantCT)

			if tc.wantBody != nil {
				var body map[string]string
				require.NoError(t, json.NewDecoder(rec.Body).Decode(&body))

				for k, v := range tc.wantBody {
					assert.Equal(t, v, body[k])
				}

				assert.Equal(t, tc.wantTimestamp, body["timestamp"])
			}

			checker.AssertExpectations(t)
		})
	}
}

func TestHandler_Check_EncodingError_LogsAndDoesNotPanic(t *testing.T) {
	t.Parallel()

	checker := buildMockChecker(t, newHealthResult(domainhealth.StatusUp, "1.0.0"))
	handler := health.NewHandler(checker)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := &brokenWriter{ResponseRecorder: *httptest.NewRecorder()}

	// Must not panic when the writer fails mid-encode.
	assert.NotPanics(t, func() {
		handler.Check(rec, req)
	})

	checker.AssertExpectations(t)
}
