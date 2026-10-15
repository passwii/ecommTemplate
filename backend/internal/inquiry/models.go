package inquiry

import "time"

// Status is the operational state of an inquiry. Visitors can never set it.
type Status string

// Operational inquiry statuses.
const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusQuoted     Status = "quoted"
	StatusWon        Status = "won"
	StatusLost       Status = "lost"
	StatusOnHold     Status = "on_hold"
)

// ValidStatus reports whether the given status is a known operational status.
func ValidStatus(status Status) bool {
	switch status {
	case StatusNew, StatusInProgress, StatusQuoted, StatusWon, StatusLost, StatusOnHold:
		return true
	default:
		return false
	}
}

// ActivityType classifies an entry in an inquiry's history.
type ActivityType string

// Activity types recorded against an inquiry.
const (
	ActivityCreated      ActivityType = "created"
	ActivityNote         ActivityType = "note"
	ActivityContact      ActivityType = "contact"
	ActivityStatusChange ActivityType = "status_change"
)

// Inquiry is a persisted buyer inquiry.
type Inquiry struct {
	ID                   string
	RequestID            string
	ContactName          string
	ContactEmail         string
	ContactCompany       string
	ContactPhone         string
	DestinationCountry   string
	ProductID            *string
	ProductSlug          string
	ProductModel         string
	ProductName          string
	Quantity             *int
	CustomizationRequest string
	PrivacyNoticeVersion string
	PrivacyConsentedAt   time.Time
	SourcePath           string
	Status               Status
	Owner                string
	LastContactedAt      *time.Time
	NextAction           string
	NextActionDueAt      *time.Time
	LossReason           string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// Activity is one entry in an inquiry's history.
type Activity struct {
	ID        int64
	InquiryID string
	Type      ActivityType
	Body      string
	Actor     string
	CreatedAt time.Time
}

// ProductSnapshot is the submission-time product identity stored with an
// inquiry so later product edits do not rewrite history.
type ProductSnapshot struct {
	ID    string
	Slug  string
	Model string
	Name  string
}
