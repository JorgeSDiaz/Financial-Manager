package delete

import "context"

// Repository is the narrow write port required by this use case.
type Repository interface {
	HasTransactions(ctx context.Context, id string) (bool, error)
	Delete(ctx context.Context, id string) error
}
