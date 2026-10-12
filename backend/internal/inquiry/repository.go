package inquiry

import (
	"context"
	"errors"
)

// ErrNotFound is returned when an inquiry does not exist.
var ErrNotFound = errors.New("inquiry: not found")

// ErrProductNotAvailable is returned when a submitted product is unknown or not
// published.
var ErrProductNotAvailable = errors.New("inquiry: product not available")

// Repository persists inquiries and their activity history.
type Repository interface {
	// Create stores the inquiry and its initial activity atomically.
	Create(ctx context.Context, value Inquiry) (Inquiry, error)
	// FindByRequestID returns an existing inquiry for an idempotency key.
	FindByRequestID(ctx context.Context, requestID string) (Inquiry, error)
	// Get returns an inquiry by id.
	Get(ctx context.Context, id string) (Inquiry, error)
}
