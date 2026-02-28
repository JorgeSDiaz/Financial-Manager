package update

import (
	"context"
	"time"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

// Repository is the narrow port required by this use case.
type Repository interface {
	GetByID(ctx context.Context, id string) (domaincategory.Category, error)
	Update(ctx context.Context, category domaincategory.Category) error
}

// Clock is the port for current-time access.
type Clock interface {
	Now() time.Time
}
