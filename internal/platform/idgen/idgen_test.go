package idgen_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/financial-manager/api/internal/platform/idgen"
)

func TestUUIDGenerator_NewID(t *testing.T) {
	t.Parallel()

	g := idgen.UUIDGenerator{}

	t.Run("returns a valid UUID v4 string", func(t *testing.T) {
		t.Parallel()

		id := g.NewID()

		parsed, err := uuid.Parse(id)
		assert.NoError(t, err)
		assert.Equal(t, uuid.Version(4), parsed.Version())
	})

	t.Run("returns a different ID on each call", func(t *testing.T) {
		t.Parallel()

		id1 := g.NewID()
		id2 := g.NewID()

		assert.NotEqual(t, id1, id2)
	})
}
