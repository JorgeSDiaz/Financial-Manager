// Package response provides shared HTTP response types and helpers for
// the transaction handler sub-packages.
package response

import (
	"encoding/json"
	"log"
	"net/http"

	domaintransaction "github.com/financial-manager/api/internal/domain/transaction"
)

const timestampLayout = "2006-01-02T15:04:05Z"
const dateLayout = "2006-01-02"

// Transaction is the JSON representation of a transaction returned by all endpoints.
type Transaction struct {
	ID          string  `json:"id"`
	AccountID   string  `json:"account_id"`
	CategoryID  string  `json:"category_id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Summary is the JSON response for the transaction summary endpoint.
type Summary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}

// Error is the JSON response body for error cases.
type Error struct {
	Error string `json:"error"`
}

// ToTransaction converts a domain transaction into its HTTP response representation.
func ToTransaction(t domaintransaction.Transaction) Transaction {
	return Transaction{
		ID:          t.ID,
		AccountID:   t.AccountID,
		CategoryID:  t.CategoryID,
		Type:        string(t.Type),
		Amount:      t.Amount,
		Description: t.Description,
		Date:        t.Date.Format(dateLayout),
		IsActive:    t.IsActive,
		CreatedAt:   t.CreatedAt.UTC().Format(timestampLayout),
		UpdatedAt:   t.UpdatedAt.UTC().Format(timestampLayout),
	}
}

// WriteJSON encodes v as JSON and writes it with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("handlers/transaction: failed to encode response: %v", err)
	}
}

// WriteError writes a JSON error response with the given status code and message.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, Error{Error: msg})
}
