package pdf_test

import (
	"context"

	"github.com/financial-manager/api/internal/application/pdfexport"
)

type fakeUseCase struct {
	pdf []byte
	err error
}

func (f *fakeUseCase) Execute(_ context.Context, _ pdfexport.Input) ([]byte, error) {
	return f.pdf, f.err
}
