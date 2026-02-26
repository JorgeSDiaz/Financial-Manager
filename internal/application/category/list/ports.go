package list

import (
	"context"

	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

// Repository is the narrow read port required by this use case.
type Repository interface {
	List(ctx context.Context, categoryType *domaincategory.Type) ([]domaincategory.Category, error)
}
