package inquiry

import (
	"context"
	"errors"
	"testing"
)

func seedInquiry(t *testing.T, repo *PostgresRepository, id, name, email string, status Status) {
	t.Helper()
	value := sampleInquiry(id, "")
	value.ContactName = name
	value.ContactEmail = email
	value.Status = status
	if _, err := repo.Create(context.Background(), value); err != nil {
		t.Fatalf("seed inquiry %s: %v", id, err)
	}
}

func TestAdminUpdateRecordsStatusChangeActivity(t *testing.T) {
	db := testDB(t)
	seedProduct(t, db)
	repo := NewPostgresRepository(db)
	seedInquiry(t, repo, "inq_1", "Alex", "alex@example.com", StatusNew)

	quoted := StatusQuoted
	updated, err := repo.Update(context.Background(), "inq_1", UpdateInput{Status: &quoted}, "owner@example.com")
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Status != StatusQuoted {
		t.Errorf("status = %q, want quoted", updated.Status)
	}

	activities, err := repo.ListActivities(context.Background(), "inq_1", 10)
	if err != nil {
		t.Fatalf("ListActivities() error = %v", err)
	}
	var statusChanges int
	for _, activity := range activities {
		if activity.Type == ActivityStatusChange {
			statusChanges++
			if activity.Actor != "owner@example.com" {
				t.Errorf("status change actor = %q", activity.Actor)
			}
		}
	}
	if statusChanges != 1 {
		t.Errorf("status change activities = %d, want 1", statusChanges)
	}
}

func TestAdminUpdateDoesNotDuplicateActivityWhenStatusIsUnchanged(t *testing.T) {
	db := testDB(t)
	seedProduct(t, db)
	repo := NewPostgresRepository(db)
	seedInquiry(t, repo, "inq_1", "Alex", "alex@example.com", StatusNew)

	same := StatusNew
	if _, err := repo.Update(context.Background(), "inq_1", UpdateInput{Status: &same}, "owner@example.com"); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	activities, err := repo.ListActivities(context.Background(), "inq_1", 10)
	if err != nil {
		t.Fatalf("ListActivities() error = %v", err)
	}
	for _, activity := range activities {
		if activity.Type == ActivityStatusChange {
			t.Fatal("unchanged status must not record a status-change activity")
		}
	}
}

func TestAdminUpdateReturnsNotFoundForUnknownInquiry(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	owner := "owner@example.com"
	if _, err := repo.Update(context.Background(), "inq_missing", UpdateInput{Owner: &owner}, "owner@example.com"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestAdminListFiltersAndCounts(t *testing.T) {
	db := testDB(t)
	seedProduct(t, db)
	repo := NewPostgresRepository(db)
	seedInquiry(t, repo, "inq_1", "Alex Morgan", "alex@example.com", StatusNew)
	seedInquiry(t, repo, "inq_2", "Bea Chen", "bea@example.com", StatusQuoted)
	seedInquiry(t, repo, "inq_3", "Carl Diaz", "carl@example.com", StatusNew)

	all, total, err := repo.List(context.Background(), ListFilter{Limit: 50})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(all) != 3 || total != 3 {
		t.Fatalf("List() = %d items, total %d, want 3/3", len(all), total)
	}

	byStatus, statusTotal, err := repo.List(context.Background(), ListFilter{Status: StatusNew, Limit: 50})
	if err != nil {
		t.Fatalf("List(status) error = %v", err)
	}
	if len(byStatus) != 2 || statusTotal != 2 {
		t.Fatalf("status filter = %d items, total %d, want 2/2", len(byStatus), statusTotal)
	}

	byQuery, queryTotal, err := repo.List(context.Background(), ListFilter{Query: "bea@example.com", Limit: 50})
	if err != nil {
		t.Fatalf("List(query) error = %v", err)
	}
	if len(byQuery) != 1 || queryTotal != 1 || byQuery[0].ContactName != "Bea Chen" {
		t.Fatalf("query filter = %+v (total %d), want Bea Chen", byQuery, queryTotal)
	}

	paged, _, err := repo.List(context.Background(), ListFilter{Limit: 2, Offset: 0})
	if err != nil {
		t.Fatalf("List(limit) error = %v", err)
	}
	if len(paged) != 2 {
		t.Fatalf("limit = %d items, want 2", len(paged))
	}
}

func TestAdminAddActivityRejectsUnknownInquiry(t *testing.T) {
	db := testDB(t)
	repo := NewPostgresRepository(db)

	_, err := repo.AddActivity(context.Background(), Activity{
		InquiryID: "inq_missing",
		Type:      ActivityNote,
		Body:      "Called the buyer",
		Actor:     "owner@example.com",
	})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("AddActivity() error = %v, want ErrNotFound", err)
	}
}

func TestAdminActivitiesAreNewestFirst(t *testing.T) {
	db := testDB(t)
	seedProduct(t, db)
	repo := NewPostgresRepository(db)
	seedInquiry(t, repo, "inq_1", "Alex", "alex@example.com", StatusNew)

	for _, body := range []string{"First note", "Second note"} {
		if _, err := repo.AddActivity(context.Background(), Activity{
			InquiryID: "inq_1",
			Type:      ActivityNote,
			Body:      body,
			Actor:     "owner@example.com",
		}); err != nil {
			t.Fatalf("AddActivity(%q) error = %v", body, err)
		}
	}

	activities, err := repo.ListActivities(context.Background(), "inq_1", 10)
	if err != nil {
		t.Fatalf("ListActivities() error = %v", err)
	}
	if len(activities) != 3 {
		t.Fatalf("activities = %d, want 3 (created + 2 notes)", len(activities))
	}
	if activities[0].Body != "Second note" {
		t.Errorf("first activity = %q, want the newest note", activities[0].Body)
	}
}
