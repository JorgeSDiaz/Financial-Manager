package health_test

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/cmd/api/handlers/health/mocks"
	domainhealth "github.com/financial-manager/api/internal/domain/health"
)

// buildMockChecker creates a mocks.Checker pre-configured to return result once.
func buildMockChecker(t *testing.T, result domainhealth.Health) *mocks.Checker {
	t.Helper()
	m := &mocks.Checker{}
	m.On("Execute", mock.Anything).Return(result, nil).Once()
	return m
}

// buildMockCheckerWithError creates a mocks.Checker pre-configured to return an error once.
func buildMockCheckerWithError(t *testing.T, err error) *mocks.Checker {
	t.Helper()
	m := &mocks.Checker{}
	m.On("Execute", mock.Anything).Return(domainhealth.Health{}, err).Once()
	return m
}

// newHealthResult builds a domain Health value with a fixed UTC timestamp for deterministic assertions.
func newHealthResult(status domainhealth.Status, version string) domainhealth.Health {
	return domainhealth.Health{
		Status:    status,
		Timestamp: time.Date(2026, 2, 23, 10, 0, 0, 0, time.UTC),
		Version:   version,
	}
}

// brokenWriter simulates an http.ResponseWriter that always fails on Write.
type brokenWriter struct {
	httptest.ResponseRecorder
}

func (b *brokenWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write error")
}
