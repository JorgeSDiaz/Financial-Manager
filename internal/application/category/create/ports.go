package create

import (
	"context"
	"time"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

// Repository is the narrow write port required by this use case.
type Repository interface {
	Create(ctx context.Context, category domaincategory.Category) error
}

// IDGenerator is the port for unique identifier generation.
type IDGenerator interface {
	NewID() string
}

// Clock is the port for current-time access.
type Clock interface {
	Now() time.Time
}
