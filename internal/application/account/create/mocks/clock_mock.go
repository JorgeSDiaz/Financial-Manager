package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// Clock is a testify mock for the create.Clock interface.
type Clock struct {
	mock.Mock
}

// Now mocks Clock.Now.
func (m *Clock) Now() time.Time {
	return m.Called().Get(0).(time.Time)
}
