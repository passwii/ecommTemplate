package inquiry

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const (
	defaultListLimit     = 50
	maxListLimit         = 200
	defaultActivityLimit = 100
)

// List returns inquiries matching the filter plus the total number of matches.
func (r *PostgresRepository) List(ctx context.Context, filter ListFilter) ([]Inquiry, int, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = defaultListLimit
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	status := string(filter.Status)
	query := strings.TrimSpace(filter.Query)
	// A single search term matches contact name, email, company or product name.
	pattern := "%" + query + "%"

	var total int
	if err := r.db.QueryRowContext(ctx, `
		SELECT count(*)
		FROM inquiries
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR contact_name ILIKE $3 OR contact_email ILIKE $3
		       OR contact_company ILIKE $3 OR product_name ILIKE $3 OR product_model ILIKE $3)`,
		status, query, pattern).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count inquiries: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT `+inquiryColumns+`
		FROM inquiries
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR contact_name ILIKE $3 OR contact_email ILIKE $3
		       OR contact_company ILIKE $3 OR product_name ILIKE $3 OR product_model ILIKE $3)
		ORDER BY created_at DESC, id
		LIMIT $4 OFFSET $5`, status, query, pattern, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list inquiries: %w", err)
	}
	defer rows.Close()

	values := []Inquiry{}
	for rows.Next() {
		value, err := scanInquiry(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan inquiry: %w", err)
		}
		values = append(values, value)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("list inquiries: %w", err)
	}
	return values, total, nil
}

// Update applies operational changes and records a status-change activity when
// the status actually changes.
func (r *PostgresRepository) Update(ctx context.Context, id string, input UpdateInput, actor string) (Inquiry, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Inquiry{}, fmt.Errorf("begin inquiry update: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var previousStatus string
	if err := tx.QueryRowContext(ctx, `SELECT status FROM inquiries WHERE id = $1 FOR UPDATE`, id).Scan(&previousStatus); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Inquiry{}, ErrNotFound
		}
		return Inquiry{}, fmt.Errorf("lock inquiry: %w", err)
	}

	row := tx.QueryRowContext(ctx, `
		UPDATE inquiries SET
			status = COALESCE($2, status),
			owner = COALESCE($3, owner),
			next_action = COALESCE($4, next_action),
			next_action_due_at = COALESCE($5, next_action_due_at),
			loss_reason = COALESCE($6, loss_reason),
			last_contacted_at = COALESCE($7, last_contacted_at),
			updated_at = now()
		WHERE id = $1
		RETURNING `+inquiryColumns,
		id, nullableStatus(input.Status), input.Owner, input.NextAction, input.NextActionDueAt, input.LossReason, input.LastContactedAt)

	updated, err := scanInquiry(row)
	if err != nil {
		return Inquiry{}, fmt.Errorf("update inquiry: %w", err)
	}

	if input.Status != nil && string(*input.Status) != previousStatus {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO inquiry_activities (inquiry_id, activity_type, body, actor)
			VALUES ($1, $2, $3, $4)`,
			id, string(ActivityStatusChange), "Status changed from "+previousStatus+" to "+string(*input.Status), actor); err != nil {
			return Inquiry{}, fmt.Errorf("record status change: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return Inquiry{}, fmt.Errorf("commit inquiry update: %w", err)
	}
	return updated, nil
}

func nullableStatus(status *Status) any {
	if status == nil {
		return nil
	}
	return string(*status)
}

// ListActivities returns an inquiry's history, newest first.
func (r *PostgresRepository) ListActivities(ctx context.Context, inquiryID string, limit int) ([]Activity, error) {
	if limit <= 0 || limit > maxListLimit {
		limit = defaultActivityLimit
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, inquiry_id, activity_type, body, actor, created_at
		FROM inquiry_activities
		WHERE inquiry_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2`, inquiryID, limit)
	if err != nil {
		return nil, fmt.Errorf("list inquiry activities: %w", err)
	}
	defer rows.Close()

	activities := []Activity{}
	for rows.Next() {
		var (
			activity Activity
			kind     string
		)
		if err := rows.Scan(&activity.ID, &activity.InquiryID, &kind, &activity.Body, &activity.Actor, &activity.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan inquiry activity: %w", err)
		}
		activity.Type = ActivityType(kind)
		activities = append(activities, activity)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list inquiry activities: %w", err)
	}
	return activities, nil
}

// AddActivity appends a note or contact activity to an existing inquiry.
func (r *PostgresRepository) AddActivity(ctx context.Context, activity Activity) (Activity, error) {
	var exists bool
	if err := r.db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM inquiries WHERE id = $1)`, activity.InquiryID).Scan(&exists); err != nil {
		return Activity{}, fmt.Errorf("check inquiry: %w", err)
	}
	if !exists {
		return Activity{}, ErrNotFound
	}

	if err := r.db.QueryRowContext(ctx, `
		INSERT INTO inquiry_activities (inquiry_id, activity_type, body, actor)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`,
		activity.InquiryID, string(activity.Type), activity.Body, activity.Actor).
		Scan(&activity.ID, &activity.CreatedAt); err != nil {
		return Activity{}, fmt.Errorf("add inquiry activity: %w", err)
	}
	return activity, nil
}
