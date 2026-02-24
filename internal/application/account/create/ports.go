package create

import (
	"context"
	"time"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// Repository is the narrow write port required by this use case.
type Repository interface {
	Create(ctx context.Context, account domainaccount.Account) error
}

// IDGenerator is the port for unique identifier generation.
type IDGenerator interface {
	NewID() string
}

// Clock is the port for current-time access.
type Clock interface {
	Now() time.Time
}
