// Package response provides shared HTTP response types and helpers for
// the account handler sub-packages.
package response

import (
	"encoding/json"
	"log"
	"net/http"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

const timestampLayout = "2006-01-02T15:04:05Z"

// Account is the JSON representation of an account returned by all endpoints.
type Account struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	InitialBalance float64 `json:"initial_balance"`
	CurrentBalance float64 `json:"current_balance"`
	Currency       string  `json:"currency"`
	Color          string  `json:"color"`
	Icon           string  `json:"icon"`
	IsActive       bool    `json:"is_active"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

// Error is the JSON response body for error cases.
type Error struct {
	Error string `json:"error"`
}

// ToAccount converts a domain account into its HTTP response representation.
func ToAccount(a domainaccount.Account) Account {
	return Account{
		ID:             a.ID,
		Name:           a.Name,
		Type:           string(a.Type),
		InitialBalance: a.InitialBalance,
		CurrentBalance: a.CurrentBalance,
		Currency:       a.Currency,
		Color:          a.Color,
		Icon:           a.Icon,
		IsActive:       a.IsActive,
		CreatedAt:      a.CreatedAt.UTC().Format(timestampLayout),
		UpdatedAt:      a.UpdatedAt.UTC().Format(timestampLayout),
	}
}

// WriteJSON encodes v as JSON and writes it with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("handlers/account: failed to encode response: %v", err)
	}
}

// WriteError writes a JSON error response with the given status code and message.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, Error{Error: msg})
}
