package oglcore

import (
	"context"
)

// Module defines the contract for all isolated domains in the monolith
type Module interface {
	// StartWorkers runs background tasks (Outbox Relays, Event Listeners).
	// It should block until the context is canceled or a fatal error occurs.
	Start(ctx context.Context) error

	// GetName returns the module name for login purpose.
	GetName() string
}
