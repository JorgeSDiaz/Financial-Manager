package create_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/application/category/create"
	"github.com/financial-manager/api/internal/application/category/create/mocks"
	domaincategory "github.com/financial-manager/api/internal/domain/category"
)

func TestUseCase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   create.Input
		repo    *mocks.Repository
		idGen   *mocks.IDGenerator
		clock   *mocks.Clock
		wantErr error
		wantOut domaincategory.Category
	}{
		{
			name: "valid input creates custom category with fixed ID and is_system false",
			input: create.Input{
				Name:  "Comida",
				Type:  "expense",
				Color: "#FF5733",
				Icon:  "restaurant",
			},
			repo:    buildMockRepo(validCategory, nil),
			idGen:   buildMockIDGenerator(),
			clock:   buildMockClock(),
			wantOut: validCategory,
		},
		{
			name:    "empty name returns validation error",
			input:   create.Input{Type: "expense", Color: "red", Icon: "icon"},
			repo:    &mocks.Repository{},
			idGen:   &mocks.IDGenerator{},
			clock:   &mocks.Clock{},
			wantErr: errors.New("category name is required"),
		},
		{
			name:    "empty color returns validation error",
			input:   create.Input{Name: "X", Type: "expense", Icon: "icon"},
			repo:    &mocks.Repository{},
			idGen:   &mocks.IDGenerator{},
			clock:   &mocks.Clock{},
			wantErr: errors.New("category color is required"),
		},
		{
			name:    "empty icon returns validation error",
			input:   create.Input{Name: "X", Type: "expense", Color: "red"},
			repo:    &mocks.Repository{},
			idGen:   &mocks.IDGenerator{},
			clock:   &mocks.Clock{},
			wantErr: errors.New("category icon is required"),
		},
		{
			name:    "invalid type returns validation error",
			input:   create.Input{Name: "X", Type: "invalid", Color: "red", Icon: "icon"},
			repo:    &mocks.Repository{},
			idGen:   &mocks.IDGenerator{},
			clock:   &mocks.Clock{},
			wantErr: fmt.Errorf(`invalid category type "invalid": must be expense or income`),
		},
		{
			name:    "repository error is wrapped and propagated",
			input:   create.Input{Name: "X", Type: "expense", Color: "red", Icon: "icon"},
			repo:    buildMockRepo(errorCategory, errors.New("db unavailable")),
			idGen:   buildMockIDGenerator(),
			clock:   buildMockClock(),
			wantErr: fmt.Errorf("create category: %w", errors.New("db unavailable")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := create.New(tc.repo, tc.idGen, tc.clock)
			out, err := uc.Execute(context.Background(), tc.input)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantOut, out)
			tc.repo.AssertExpectations(t)
			tc.idGen.AssertExpectations(t)
			tc.clock.AssertExpectations(t)
		})
	}
}
