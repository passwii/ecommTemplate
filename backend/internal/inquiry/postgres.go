package inquiry

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// PostgresRepository is the SQL-backed implementation of Repository.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository builds an inquiry repository over the given database.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

const inquiryColumns = `
	id, COALESCE(request_id, ''), contact_name, contact_email, contact_company, contact_phone,
	destination_country, product_id, product_slug, product_model, product_name,
	quantity, customization_request, privacy_notice_version, privacy_consented_at, source_path,
	status, owner, last_contacted_at, next_action, next_action_due_at, loss_reason, created_at, updated_at`

// Create stores the inquiry and its initial activity in one transaction.
func (r *PostgresRepository) Create(ctx context.Context, value Inquiry) (Inquiry, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Inquiry{}, fmt.Errorf("begin inquiry transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var requestID any
	if value.RequestID != "" {
		requestID = value.RequestID
	}

	row := tx.QueryRowContext(ctx, `
		INSERT INTO inquiries (
			id, request_id, contact_name, contact_email, contact_company, contact_phone,
			destination_country, product_id, product_slug, product_model, product_name,
			quantity, customization_request, privacy_notice_version, privacy_consented_at, source_path,
			status, owner
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16,
			$17, $18
		)
		RETURNING `+inquiryColumns,
		value.ID, requestID, value.ContactName, value.ContactEmail, value.ContactCompany, value.ContactPhone,
		value.DestinationCountry, value.ProductID, value.ProductSlug, value.ProductModel, value.ProductName,
		value.Quantity, value.CustomizationRequest, value.PrivacyNoticeVersion, value.PrivacyConsentedAt, value.SourcePath,
		string(value.Status), value.Owner)

	stored, err := scanInquiry(row)
	if err != nil {
		return Inquiry{}, fmt.Errorf("insert inquiry: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO inquiry_activities (inquiry_id, activity_type, body, actor)
		VALUES ($1, $2, $3, $4)`,
		stored.ID, string(ActivityCreated), "Inquiry received", "system"); err != nil {
		return Inquiry{}, fmt.Errorf("insert inquiry activity: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return Inquiry{}, fmt.Errorf("commit inquiry transaction: %w", err)
	}
	return stored, nil
}

// FindByRequestID returns an existing inquiry for an idempotency key.
func (r *PostgresRepository) FindByRequestID(ctx context.Context, requestID string) (Inquiry, error) {
	if requestID == "" {
		return Inquiry{}, ErrNotFound
	}
	row := r.db.QueryRowContext(ctx, `SELECT `+inquiryColumns+` FROM inquiries WHERE request_id = $1`, requestID)
	stored, err := scanInquiry(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Inquiry{}, ErrNotFound
	}
	if err != nil {
		return Inquiry{}, fmt.Errorf("find inquiry by request id: %w", err)
	}
	return stored, nil
}

// Get returns an inquiry by id.
func (r *PostgresRepository) Get(ctx context.Context, id string) (Inquiry, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+inquiryColumns+` FROM inquiries WHERE id = $1`, id)
	stored, err := scanInquiry(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Inquiry{}, ErrNotFound
	}
	if err != nil {
		return Inquiry{}, fmt.Errorf("get inquiry: %w", err)
	}
	return stored, nil
}

func scanInquiry(row interface{ Scan(dest ...any) error }) (Inquiry, error) {
	var (
		value           Inquiry
		productID       sql.NullString
		quantity        sql.NullInt64
		lastContactedAt sql.NullTime
		nextActionDueAt sql.NullTime
		status          string
	)

	err := row.Scan(
		&value.ID, &value.RequestID, &value.ContactName, &value.ContactEmail, &value.ContactCompany, &value.ContactPhone,
		&value.DestinationCountry, &productID, &value.ProductSlug, &value.ProductModel, &value.ProductName,
		&quantity, &value.CustomizationRequest, &value.PrivacyNoticeVersion, &value.PrivacyConsentedAt, &value.SourcePath,
		&status, &value.Owner, &lastContactedAt, &value.NextAction, &nextActionDueAt, &value.LossReason, &value.CreatedAt, &value.UpdatedAt,
	)
	if err != nil {
		return Inquiry{}, err
	}

	if productID.Valid {
		id := productID.String
		value.ProductID = &id
	}
	if quantity.Valid {
		amount := int(quantity.Int64)
		value.Quantity = &amount
	}
	if lastContactedAt.Valid {
		at := lastContactedAt.Time
		value.LastContactedAt = &at
	}
	if nextActionDueAt.Valid {
		at := nextActionDueAt.Time
		value.NextActionDueAt = &at
	}
	value.Status = Status(status)
	return value, nil
}
