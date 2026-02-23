// Package health contains the Health entity, value objects, and constants.
package health

import "time"

type (
	// Status represents the current operational state of the application.
	Status string

	// Health represents the health state of the application at a point in time.
	Health struct {
		Status    Status
		Timestamp time.Time
		Version   string
	}
)

const (
	// StatusUp indicates the service is operating normally.
	StatusUp Status = "up"
	// StatusDown indicates the service is not operational.
	StatusDown Status = "down"
)
