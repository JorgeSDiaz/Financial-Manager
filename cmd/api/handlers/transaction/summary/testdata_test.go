package summary_test

import (
	"context"

	appsummary "github.com/financial-manager/api/internal/application/transaction/summary"
)

type fakeUseCase struct {
	out appsummary.Summary
	err error
}

func (f *fakeUseCase) Execute(_ context.Context, _ appsummary.Input) (appsummary.Summary, error) {
	return f.out, f.err
}
