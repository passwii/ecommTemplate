package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/passwii/ecommTemple/backend/internal/admin"
	"github.com/passwii/ecommTemple/backend/internal/inquiry"
)

type fakeAdminRepository struct {
	users    map[string]admin.User
	sessions map[string]admin.Session
}

func newFakeAdminRepository() *fakeAdminRepository {
	return &fakeAdminRepository{users: map[string]admin.User{}, sessions: map[string]admin.Session{}}
}

func (f *fakeAdminRepository) CreateUser(_ context.Context, user admin.User) error {
	if _, exists := f.users[user.Email]; exists {
		return admin.ErrEmailTaken
	}
	f.users[user.Email] = user
	return nil
}

func (f *fakeAdminRepository) FindUserByEmail(_ context.Context, email string) (admin.User, error) {
	if user, ok := f.users[email]; ok {
		return user, nil
	}
	return admin.User{}, admin.ErrNotFound
}

func (f *fakeAdminRepository) FindUserByID(_ context.Context, id string) (admin.User, error) {
	for _, user := range f.users {
		if user.ID == id {
			return user, nil
		}
	}
	return admin.User{}, admin.ErrNotFound
}

func (f *fakeAdminRepository) CreateSession(_ context.Context, session admin.Session) error {
	f.sessions[session.TokenHash] = session
	return nil
}

func (f *fakeAdminRepository) FindSession(_ context.Context, tokenHash string) (admin.Session, error) {
	if session, ok := f.sessions[tokenHash]; ok {
		return session, nil
	}
	return admin.Session{}, admin.ErrNotFound
}

func (f *fakeAdminRepository) TouchSession(context.Context, string, time.Time) error { return nil }

func (f *fakeAdminRepository) DeleteSession(_ context.Context, tokenHash string) error {
	delete(f.sessions, tokenHash)
	return nil
}

func (f *fakeAdminRepository) DeleteExpiredSessions(context.Context, time.Time) error { return nil }

type adminInquiryStore struct {
	value      inquiry.Inquiry
	activities []inquiry.Activity
	updated    inquiry.UpdateInput
	updateErr  error
}

func (s *adminInquiryStore) Create(context.Context, inquiry.Inquiry) (inquiry.Inquiry, error) {
	return inquiry.Inquiry{}, nil
}

func (s *adminInquiryStore) FindByRequestID(context.Context, string) (inquiry.Inquiry, error) {
	return inquiry.Inquiry{}, inquiry.ErrNotFound
}

func (s *adminInquiryStore) Get(context.Context, string) (inquiry.Inquiry, error) {
	if s.value.ID == "" {
		return inquiry.Inquiry{}, inquiry.ErrNotFound
	}
	return s.value, nil
}

func (s *adminInquiryStore) List(context.Context, inquiry.ListFilter) ([]inquiry.Inquiry, int, error) {
	if s.value.ID == "" {
		return []inquiry.Inquiry{}, 0, nil
	}
	return []inquiry.Inquiry{s.value}, 1, nil
}

func (s *adminInquiryStore) Update(_ context.Context, _ string, input inquiry.UpdateInput, _ string) (inquiry.Inquiry, error) {
	if s.updateErr != nil {
		return inquiry.Inquiry{}, s.updateErr
	}
	s.updated = input
	updated := s.value
	if input.Status != nil {
		updated.Status = *input.Status
	}
	if input.Owner != nil {
		updated.Owner = *input.Owner
	}
	s.value = updated
	return updated, nil
}

func (s *adminInquiryStore) ListActivities(context.Context, string, int) ([]inquiry.Activity, error) {
	return s.activities, nil
}

func (s *adminInquiryStore) AddActivity(_ context.Context, activity inquiry.Activity) (inquiry.Activity, error) {
	if s.value.ID == "" {
		return inquiry.Activity{}, inquiry.ErrNotFound
	}
	activity.ID = 1
	activity.CreatedAt = time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	s.activities = append(s.activities, activity)
	return activity, nil
}

type adminTestEnv struct {
	handler   http.Handler
	users     *fakeAdminRepository
	inquiries *adminInquiryStore
}

func newAdminTestEnv(t *testing.T) adminTestEnv {
	t.Helper()
	users := newFakeAdminRepository()
	service := admin.NewService(users)
	if _, err := service.CreateUser(context.Background(), "owner@example.com", "correct horse battery"); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	store := &adminInquiryStore{value: inquiry.Inquiry{
		ID:           "inq_1",
		ContactName:  "Alex Morgan",
		ContactEmail: "buyer@example.com",
		ProductSlug:  "ridge-4-burner",
		ProductModel: "Ridge 4-Burner",
		ProductName:  "Ridge 4-Burner",
		Status:       inquiry.StatusNew,
		CreatedAt:    time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC),
	}}

	handler := NewHandler(Deps{
		Admin:             service,
		AdminInquiryStore: store,
	})
	return adminTestEnv{handler: handler, users: users, inquiries: store}
}

func adminRequest(t *testing.T, handler http.Handler, method, target, body, cookie string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	request := httptest.NewRequest(method, target, reader)
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	if cookie != "" {
		request.Header.Set("Cookie", cookie)
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder
}

func loginCookie(t *testing.T, env adminTestEnv) string {
	t.Helper()
	recorder := adminRequest(t, env.handler, http.MethodPost, "/api/v1/admin/sessions",
		`{"email":"owner@example.com","password":"correct horse battery"}`, "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("login status = %d, want 200 (%s)", recorder.Code, recorder.Body.String())
	}
	cookies := recorder.Result().Cookies()
	for _, cookie := range cookies {
		if cookie.Name == AdminSessionCookie {
			if !cookie.HttpOnly {
				t.Error("session cookie must be HttpOnly")
			}
			if cookie.SameSite != http.SameSiteLaxMode {
				t.Error("session cookie must be SameSite=Lax")
			}
			return AdminSessionCookie + "=" + cookie.Value
		}
	}
	t.Fatal("login response did not set the session cookie")
	return ""
}

func TestAdminLoginRejectsBadCredentialsWithoutRevealingWhich(t *testing.T) {
	env := newAdminTestEnv(t)

	for _, body := range []string{
		`{"email":"owner@example.com","password":"wrong password here"}`,
		`{"email":"nobody@example.com","password":"correct horse battery"}`,
	} {
		recorder := adminRequest(t, env.handler, http.MethodPost, "/api/v1/admin/sessions", body, "")
		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want 401", recorder.Code)
		}
		var envelope struct {
			Error struct {
				Code string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &envelope); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if envelope.Error.Code != "invalid_credentials" {
			t.Errorf("code = %q, want invalid_credentials", envelope.Error.Code)
		}
	}
}

func TestAdminProtectedRoutesRequireASession(t *testing.T) {
	env := newAdminTestEnv(t)

	targets := []struct{ method, target string }{
		{http.MethodGet, "/api/v1/admin/inquiries"},
		{http.MethodGet, "/api/v1/admin/inquiries/inq_1"},
		{http.MethodPatch, "/api/v1/admin/inquiries/inq_1"},
		{http.MethodPost, "/api/v1/admin/inquiries/inq_1/activities"},
	}
	for _, target := range targets {
		recorder := adminRequest(t, env.handler, target.method, target.target, "", "")
		if recorder.Code != http.StatusUnauthorized {
			t.Errorf("%s %s status = %d, want 401", target.method, target.target, recorder.Code)
		}
	}
}

func TestAdminRejectsForgedSessionCookie(t *testing.T) {
	env := newAdminTestEnv(t)
	recorder := adminRequest(t, env.handler, http.MethodGet, "/api/v1/admin/inquiries", "", AdminSessionCookie+"=forged-token")
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", recorder.Code)
	}
}

func TestAdminRejectsExpiredSession(t *testing.T) {
	env := newAdminTestEnv(t)
	cookie := loginCookie(t, env)

	// Expire the stored session and try again.
	for hash, session := range env.users.sessions {
		session.ExpiresAt = time.Now().Add(-time.Hour)
		env.users.sessions[hash] = session
	}

	recorder := adminRequest(t, env.handler, http.MethodGet, "/api/v1/admin/inquiries", "", cookie)
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401 for an expired session", recorder.Code)
	}
}

func TestAdminListsAndReadsInquiries(t *testing.T) {
	env := newAdminTestEnv(t)
	cookie := loginCookie(t, env)

	list := adminRequest(t, env.handler, http.MethodGet, "/api/v1/admin/inquiries", "", cookie)
	if list.Code != http.StatusOK {
		t.Fatalf("list status = %d, want 200 (%s)", list.Code, list.Body.String())
	}
	var listBody struct {
		Data []map[string]any `json:"data"`
		Meta struct {
			Total int `json:"total"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(list.Body.Bytes(), &listBody); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(listBody.Data) != 1 || listBody.Meta.Total != 1 {
		t.Fatalf("list = %+v, want one inquiry", listBody)
	}

	detail := adminRequest(t, env.handler, http.MethodGet, "/api/v1/admin/inquiries/inq_1", "", cookie)
	if detail.Code != http.StatusOK {
		t.Fatalf("detail status = %d, want 200", detail.Code)
	}
	if !strings.Contains(detail.Body.String(), "ridge-4-burner") {
		t.Errorf("detail body missing product snapshot: %s", detail.Body.String())
	}
}

func TestAdminRejectsInvalidStatusFilter(t *testing.T) {
	env := newAdminTestEnv(t)
	cookie := loginCookie(t, env)

	recorder := adminRequest(t, env.handler, http.MethodGet, "/api/v1/admin/inquiries?status=bogus", "", cookie)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", recorder.Code)
	}
}

func TestAdminUpdatesStatusAndRejectsUnknownStatus(t *testing.T) {
	env := newAdminTestEnv(t)
	cookie := loginCookie(t, env)

	bad := adminRequest(t, env.handler, http.MethodPatch, "/api/v1/admin/inquiries/inq_1", `{"status":"bogus"}`, cookie)
	if bad.Code != http.StatusBadRequest {
		t.Fatalf("invalid status = %d, want 400", bad.Code)
	}

	good := adminRequest(t, env.handler, http.MethodPatch, "/api/v1/admin/inquiries/inq_1", `{"status":"quoted","nextAction":"Send pricing"}`, cookie)
	if good.Code != http.StatusOK {
		t.Fatalf("valid status = %d, want 200 (%s)", good.Code, good.Body.String())
	}
	if env.inquiries.updated.Status == nil || *env.inquiries.updated.Status != inquiry.StatusQuoted {
		t.Errorf("stored status = %v, want quoted", env.inquiries.updated.Status)
	}
}

func TestAdminRejectsClientSuppliedOperationalOverrides(t *testing.T) {
	env := newAdminTestEnv(t)
	cookie := loginCookie(t, env)

	// created_at is not an accepted field; strict decoding must reject it.
	recorder := adminRequest(t, env.handler, http.MethodPatch, "/api/v1/admin/inquiries/inq_1", `{"createdAt":"2020-01-01T00:00:00Z"}`, cookie)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 for unknown field", recorder.Code)
	}
}

func TestAdminAddsActivitiesAndValidatesThem(t *testing.T) {
	env := newAdminTestEnv(t)
	cookie := loginCookie(t, env)

	badType := adminRequest(t, env.handler, http.MethodPost, "/api/v1/admin/inquiries/inq_1/activities", `{"type":"created","body":"x"}`, cookie)
	if badType.Code != http.StatusBadRequest {
		t.Fatalf("invalid type = %d, want 400", badType.Code)
	}

	emptyBody := adminRequest(t, env.handler, http.MethodPost, "/api/v1/admin/inquiries/inq_1/activities", `{"type":"note","body":"  "}`, cookie)
	if emptyBody.Code != http.StatusBadRequest {
		t.Fatalf("empty body = %d, want 400", emptyBody.Code)
	}

	ok := adminRequest(t, env.handler, http.MethodPost, "/api/v1/admin/inquiries/inq_1/activities", `{"type":"note","body":"Called the buyer"}`, cookie)
	if ok.Code != http.StatusCreated {
		t.Fatalf("valid note = %d, want 201 (%s)", ok.Code, ok.Body.String())
	}
	if len(env.inquiries.activities) != 1 {
		t.Fatalf("activities = %d, want 1", len(env.inquiries.activities))
	}
	if env.inquiries.activities[0].Actor != "owner@example.com" {
		t.Errorf("actor = %q, want the signed-in admin", env.inquiries.activities[0].Actor)
	}
}

func TestAdminLogoutClearsTheSession(t *testing.T) {
	env := newAdminTestEnv(t)
	cookie := loginCookie(t, env)

	logout := adminRequest(t, env.handler, http.MethodDelete, "/api/v1/admin/sessions", "", cookie)
	if logout.Code != http.StatusOK {
		t.Fatalf("logout status = %d, want 200", logout.Code)
	}

	after := adminRequest(t, env.handler, http.MethodGet, "/api/v1/admin/inquiries", "", cookie)
	if after.Code != http.StatusUnauthorized {
		t.Fatalf("status after logout = %d, want 401", after.Code)
	}
}
