// Package health contains the HTTP handler for the /health endpoint.
package health

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	domainhealth "github.com/financial-manager/api/internal/domain/health"
)

const (
	// timestampLayout is the ISO 8601 UTC format used in health responses.
	timestampLayout = "2006-01-02T15:04:05Z"
)

type (
	// checker is the contract the handler depends on to retrieve health status.
	checker interface {
		Execute(ctx context.Context) (domainhealth.Health, error)
	}

	// Handler handles HTTP requests for the health endpoint.
	Handler struct {
		checker checker
	}

	// response is the JSON shape returned by the health endpoint.
	response struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
		Version   string `json:"version"`
	}
)

// NewHandler creates a Handler with its required dependency.
func NewHandler(c checker) *Handler {
	return &Handler{checker: c}
}

// Check handles GET /health and returns the current application health status.
func (h *Handler) Check(w http.ResponseWriter, r *http.Request) {
	result, err := h.checker.Execute(r.Context())
	if err != nil {
		http.Error(w, `{"error":"service unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	resp := response{
		Status:    string(result.Status),
		Timestamp: result.Timestamp.Format(timestampLayout),
		Version:   result.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Headers already sent; log the failure — cannot change status code at this point.
		log.Printf("handlers/health: failed to encode health response: %v", err)
	}
}
