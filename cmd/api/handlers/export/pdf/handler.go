// Package pdf handles POST /api/v1/export/pdf.
package pdf

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/financial-manager/api/internal/application/pdfexport"
)

type useCase interface {
	Execute(ctx context.Context, in pdfexport.Input) ([]byte, error)
}

// Handler handles POST /api/v1/export/pdf.
type Handler struct {
	uc useCase
}

// New creates a Handler with its required use case dependency.
func New(uc useCase) *Handler {
	return &Handler{uc: uc}
}

// Handle processes POST /api/v1/export/pdf.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var input pdfexport.Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		// If no body or invalid JSON, use current month as default
		now := time.Now()
		input.Month = now.Format("2006-01")
	}

	// Validate month format
	if input.Month == "" {
		now := time.Now()
		input.Month = now.Format("2006-01")
	}

	pdfData, err := h.uc.Execute(r.Context(), input)
	if err != nil {
		// Check if it's a validation error (invalid month format)
		if err.Error() == "invalid month format: parsing time" {
			http.Error(w, `{"error":"invalid month format, expected YYYY-MM"}`, http.StatusBadRequest)
			return
		}
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	filename := "report_" + input.Month + ".pdf"
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.WriteHeader(http.StatusOK)
	w.Write(pdfData)
}
