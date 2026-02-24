package clock_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/platform/clock"
)

func TestWallClock_Now(t *testing.T) {
	t.Parallel()

	c := clock.WallClock{}

	t.Run("returns a time close to the current UTC time", func(t *testing.T) {
		t.Parallel()

		before := time.Now().UTC()
		now := c.Now()
		after := time.Now().UTC()

		assert.True(t, !now.Before(before) && !now.After(after))
	})

	t.Run("returns time in UTC", func(t *testing.T) {
		t.Parallel()

		now := c.Now()

		assert.Equal(t, time.UTC, now.Location())
	})
}
