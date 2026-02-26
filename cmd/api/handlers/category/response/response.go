// Package response provides shared HTTP response types and helpers for
// the category handler sub-packages.
package response

import (
	"encoding/json"
	"log"
	"net/http"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

const timestampLayout = "2006-01-02T15:04:05Z"

// CategoryResponse is the JSON representation of a category returned by all endpoints.
type CategoryResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Color     string `json:"color"`
	Icon      string `json:"icon"`
	IsSystem  bool   `json:"is_system"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Error is the JSON response body for error cases.
type Error struct {
	Error string `json:"error"`
}

// ToCategory converts a domain category into its HTTP response representation.
func ToCategory(c domaincategory.Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		Name:      c.Name,
		Type:      string(c.Type),
		Color:     c.Color,
		Icon:      c.Icon,
		IsSystem:  c.IsSystem,
		IsActive:  c.IsActive,
		CreatedAt: c.CreatedAt.UTC().Format(timestampLayout),
		UpdatedAt: c.UpdatedAt.UTC().Format(timestampLayout),
	}
}

// WriteJSON encodes v as JSON and writes it with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("handlers/category: failed to encode response: %v", err)
	}
}

// WriteError writes a JSON error response with the given status code and message.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, Error{Error: msg})
}
