package inquiry

import (
	"context"
	"errors"
	"time"
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
	// List returns inquiries matching the filter plus the total match count.
	List(ctx context.Context, filter ListFilter) ([]Inquiry, int, error)
	// Update applies operational changes and records a status-change activity
	// when the status changes.
	Update(ctx context.Context, id string, input UpdateInput, actor string) (Inquiry, error)
	// ListActivities returns an inquiry's history, newest first.
	ListActivities(ctx context.Context, inquiryID string, limit int) ([]Activity, error)
	// AddActivity appends a note or contact activity.
	AddActivity(ctx context.Context, activity Activity) (Activity, error)
}

// ListFilter narrows an administrative inquiry list query.
type ListFilter struct {
	Status Status
	Query  string
	Limit  int
	Offset int
}

// UpdateInput carries operational changes. Nil fields are left unchanged.
type UpdateInput struct {
	Status          *Status
	Owner           *string
	NextAction      *string
	NextActionDueAt *time.Time
	LossReason      *string
	LastContactedAt *time.Time
}
