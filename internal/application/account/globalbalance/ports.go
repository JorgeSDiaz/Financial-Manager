package globalbalance

import (
	"context"

	domainaccount "github.com/financial-manager/api/internal/domain/account"
)

// Repository is the narrow read port required by this use case.
type Repository interface {
	List(ctx context.Context) ([]domainaccount.Account, error)
}
