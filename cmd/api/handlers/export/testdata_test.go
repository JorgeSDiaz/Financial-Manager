package export_test

import (
	"context"

	"github.com/financial-manager/api/internal/application/export"
)

type fakeCSVUseCase struct {
	csv string
	err error
}

func (f *fakeCSVUseCase) ExportCSV(_ context.Context, _ export.CSVFilters) (string, error) {
	return f.csv, f.err
}

type fakeJSONUseCase struct {
	json []byte
	err  error
}

func (f *fakeJSONUseCase) ExportJSON(_ context.Context) ([]byte, error) {
	return f.json, f.err
}
