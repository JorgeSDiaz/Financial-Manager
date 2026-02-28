package delete_test

import (
	"context"
)

type fakeUseCase struct {
	err error
}

func (f *fakeUseCase) Execute(_ context.Context, id string) error {
	return f.err
}
