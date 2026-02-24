package get

import (
	"context"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// Repository is the narrow read port required by this use case.
type Repository interface {
	GetByID(ctx context.Context, id string) (domainaccount.Account, error)
}
