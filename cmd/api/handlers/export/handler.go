// Package export handles GET /api/v1/export/csv and GET /api/v1/export/json.
package export

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/financial-manager/api/internal/application/export"
)

type csvUseCase interface {
	ExportCSV(ctx context.Context, filters export.CSVFilters) (string, error)
}

type jsonUseCase interface {
	ExportJSON(ctx context.Context) ([]byte, error)
}

// Handler handles export endpoints.
type Handler struct {
	csvUC  csvUseCase
	jsonUC jsonUseCase
}

// New creates a Handler with its required use case dependencies.
func New(csvUC csvUseCase, jsonUC jsonUseCase) *Handler {
	return &Handler{csvUC: csvUC, jsonUC: jsonUC}
}

// HandleCSV processes GET /api/v1/export/csv.
func (h *Handler) HandleCSV(w http.ResponseWriter, r *http.Request) {
	filters := export.CSVFilters{
		StartDate: r.URL.Query().Get("date_from"),
		EndDate:   r.URL.Query().Get("date_to"),
		Type:      r.URL.Query().Get("type"),
	}

	csvData, err := h.csvUC.ExportCSV(r.Context(), filters)
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("transactions_%s.csv", time.Now().Format("2006-01"))
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(csvData))
}

// HandleJSON processes GET /api/v1/export/json.
func (h *Handler) HandleJSON(w http.ResponseWriter, r *http.Request) {
	jsonData, err := h.jsonUC.ExportJSON(r.Context())
	if err != nil {
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("backup_%s.json", time.Now().Format("2006-01-02"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
