package delete

import (
	"context"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

// Repository is the narrow port required by this use case.
type Repository interface {
	GetByID(ctx context.Context, id string) (domaincategory.Category, error)
	Delete(ctx context.Context, id string) error
	HasTransactions(ctx context.Context, id string) (bool, error)
}
