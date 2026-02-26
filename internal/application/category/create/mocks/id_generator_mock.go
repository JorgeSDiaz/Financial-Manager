package mocks

import (
	"github.com/stretchr/testify/mock"
)

// IDGenerator is a testify mock for the create.IDGenerator interface.
type IDGenerator struct {
	mock.Mock
}

// NewID mocks IDGenerator.NewID.
func (m *IDGenerator) NewID() string {
	return m.Called().String(0)
}
