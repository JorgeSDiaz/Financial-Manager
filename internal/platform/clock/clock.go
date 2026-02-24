// Package clock provides a wall clock implementation for production use.
package clock

import "time"

// WallClock is the production clock that returns the real current time.
type WallClock struct{}

// Now returns the current UTC time.
func (WallClock) Now() time.Time {
	return time.Now().UTC()
}
