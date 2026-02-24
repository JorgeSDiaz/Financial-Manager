package update

import (
	"context"
	"time"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// Repository is the narrow read-write port required by this use case.
type Repository interface {
	GetByID(ctx context.Context, id string) (domainaccount.Account, error)
	Update(ctx context.Context, account domainaccount.Account) error
}

// Clock is the port for current-time access.
type Clock interface {
	Now() time.Time
}
