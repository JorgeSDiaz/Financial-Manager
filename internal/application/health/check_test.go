package health_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/financial-manager/api/internal/application/health"
	domainhealth "github.com/financial-manager/api/internal/domain/health"
)

func TestCheckUseCase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		wantStatus domainhealth.Status
		wantErr    bool
	}{
		{
			name:       "returns StatusUp with no error",
			wantStatus: domainhealth.StatusUp,
			wantErr:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			uc := health.NewCheckUseCase()
			before := time.Now().UTC()
			got, err := uc.Execute(context.Background())
			after := time.Now().UTC()

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantStatus, got.Status)
			assert.NotEmpty(t, got.Version)
			assert.True(t,
				!got.Timestamp.Before(before) && !got.Timestamp.After(after),
				"Timestamp %v must be between %v and %v", got.Timestamp, before, after,
			)
		})
	}
}
