package dashboard_test

import (
	"context"

	appDashboard "github.com/financial-manager/api/internal/application/dashboard"
)

type fakeUseCase struct {
	out appDashboard.Output
	err error
}

func (f *fakeUseCase) Execute(_ context.Context) (appDashboard.Output, error) {
	return f.out, f.err
}
