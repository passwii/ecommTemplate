package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/passwii/ecommTemple/backend/internal/admin"
	"github.com/passwii/ecommTemple/backend/internal/inquiry"
)

// AdminSessionCookie is the session cookie name used by the admin UI.
const AdminSessionCookie = "storefront_admin_session"

const maxAdminBodyBytes = 32 * 1024

type adminHandlers struct {
	service      *admin.Service
	inquiries    inquiry.Repository
	secureCookie bool
}

func registerAdminRoutes(mux *http.ServeMux, service *admin.Service, inquiries inquiry.Repository, secureCookie bool) {
	handlers := &adminHandlers{service: service, inquiries: inquiries, secureCookie: secureCookie}

	mux.HandleFunc("POST /api/v1/admin/sessions", handlers.login)
	mux.HandleFunc("DELETE /api/v1/admin/sessions", handlers.logout)
	mux.HandleFunc("GET /api/v1/admin/session", handlers.session)
	mux.HandleFunc("GET /api/v1/admin/inquiries", handlers.requireAuth(handlers.listInquiries))
	mux.HandleFunc("GET /api/v1/admin/inquiries/{id}", handlers.requireAuth(handlers.getInquiry))
	mux.HandleFunc("PATCH /api/v1/admin/inquiries/{id}", handlers.requireAuth(handlers.updateInquiry))
	mux.HandleFunc("POST /api/v1/admin/inquiries/{id}/activities", handlers.requireAuth(handlers.addActivity))
}

// requireAuth rejects unauthenticated requests before the wrapped handler runs.
// Server-side enforcement is the only authorization boundary; hiding UI routes
// is not a security control.
func (h *adminHandlers) requireAuth(next func(http.ResponseWriter, *http.Request, admin.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := h.service.Authenticate(r.Context(), h.sessionToken(r))
		if err != nil {
			if !errors.Is(err, admin.ErrInvalidCredentials) {
				slog.Error("admin session lookup failed", "error", err)
			}
			writeError(w, http.StatusUnauthorized, "unauthorized", "Sign in to continue.")
			return
		}
		next(w, r, user)
	}
}

func (h *adminHandlers) sessionToken(r *http.Request) string {
	cookie, err := r.Cookie(AdminSessionCookie)
	if err != nil {
		return ""
	}
	return cookie.Value
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *adminHandlers) login(w http.ResponseWriter, r *http.Request) {
	if !isJSONContentType(r.Header.Get("Content-Type")) {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Send the request as JSON.")
		return
	}

	var payload loginPayload
	if !decodeStrictJSON(w, r, &payload, maxAdminBodyBytes) {
		return
	}

	token, user, err := h.service.Login(r.Context(), payload.Email, payload.Password)
	if errors.Is(err, admin.ErrInvalidCredentials) {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "The email or password is incorrect.")
		return
	}
	if err != nil {
		slog.Error("admin login failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "The request could not be completed.")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     AdminSessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.secureCookie,
		MaxAge:   int(admin.DefaultSessionTTL.Seconds()),
	})
	writeData(w, http.StatusOK, map[string]any{"email": user.Email})
}

func (h *adminHandlers) logout(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Logout(r.Context(), h.sessionToken(r)); err != nil {
		slog.Error("admin logout failed", "error", err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     AdminSessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.secureCookie,
		MaxAge:   -1,
	})
	writeData(w, http.StatusOK, map[string]any{"signedOut": true})
}

func (h *adminHandlers) session(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.Authenticate(r.Context(), h.sessionToken(r))
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Sign in to continue.")
		return
	}
	writeData(w, http.StatusOK, map[string]any{"email": user.Email})
}

func (h *adminHandlers) listInquiries(w http.ResponseWriter, r *http.Request, _ admin.User) {
	query := r.URL.Query()

	filter := inquiry.ListFilter{Query: strings.TrimSpace(query.Get("q"))}
	if raw := query.Get("status"); raw != "" {
		status := inquiry.Status(raw)
		if !inquiry.ValidStatus(status) {
			writeFieldErrors(w, http.StatusBadRequest, "invalid_query", "Invalid query parameters.", map[string]string{
				"status": "Unknown inquiry status.",
			})
			return
		}
		filter.Status = status
	}
	if raw := query.Get("limit"); raw != "" {
		limit, err := strconv.Atoi(raw)
		if err != nil || limit < 1 || limit > 200 {
			writeFieldErrors(w, http.StatusBadRequest, "invalid_query", "Invalid query parameters.", map[string]string{
				"limit": "Use a whole number between 1 and 200.",
			})
			return
		}
		filter.Limit = limit
	}
	if raw := query.Get("offset"); raw != "" {
		offset, err := strconv.Atoi(raw)
		if err != nil || offset < 0 {
			writeFieldErrors(w, http.StatusBadRequest, "invalid_query", "Invalid query parameters.", map[string]string{
				"offset": "Use a whole number of zero or more.",
			})
			return
		}
		filter.Offset = offset
	}

	values, total, err := h.inquiries.List(r.Context(), filter)
	if err != nil {
		slog.Error("admin inquiry list failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "The request could not be completed.")
		return
	}

	items := make([]map[string]any, 0, len(values))
	for _, value := range values {
		items = append(items, inquiryToJSON(value))
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"data": items,
		"meta": map[string]any{"total": total, "limit": filter.Limit, "offset": filter.Offset},
	})
}

func (h *adminHandlers) getInquiry(w http.ResponseWriter, r *http.Request, _ admin.User) {
	id := r.PathValue("id")
	value, err := h.inquiries.Get(r.Context(), id)
	if errors.Is(err, inquiry.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", "The requested inquiry was not found.")
		return
	}
	if err != nil {
		slog.Error("admin inquiry read failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "The request could not be completed.")
		return
	}

	activities, err := h.inquiries.ListActivities(r.Context(), id, 100)
	if err != nil {
		slog.Error("admin inquiry activity read failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "The request could not be completed.")
		return
	}

	payload := inquiryToJSON(value)
	payload["activities"] = activitiesToJSON(activities)
	writeData(w, http.StatusOK, payload)
}

type updateInquiryPayload struct {
	Status          *string `json:"status"`
	Owner           *string `json:"owner"`
	NextAction      *string `json:"nextAction"`
	NextActionDueAt *string `json:"nextActionDueAt"`
	LossReason      *string `json:"lossReason"`
	LastContactedAt *string `json:"lastContactedAt"`
}

func (h *adminHandlers) updateInquiry(w http.ResponseWriter, r *http.Request, user admin.User) {
	if !isJSONContentType(r.Header.Get("Content-Type")) {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Send the request as JSON.")
		return
	}

	var payload updateInquiryPayload
	if !decodeStrictJSON(w, r, &payload, maxAdminBodyBytes) {
		return
	}

	input := inquiry.UpdateInput{}
	if payload.Status != nil {
		status := inquiry.Status(*payload.Status)
		if !inquiry.ValidStatus(status) {
			writeFieldErrors(w, http.StatusBadRequest, "validation_error", "Please correct the highlighted fields.", map[string]string{
				"status": "Unknown inquiry status.",
			})
			return
		}
		input.Status = &status
	}
	if payload.Owner != nil {
		value := strings.TrimSpace(*payload.Owner)
		input.Owner = &value
	}
	if payload.NextAction != nil {
		value := strings.TrimSpace(*payload.NextAction)
		input.NextAction = &value
	}
	if payload.LossReason != nil {
		value := strings.TrimSpace(*payload.LossReason)
		input.LossReason = &value
	}
	if payload.NextActionDueAt != nil {
		parsed, ok := parseOptionalTimestamp(w, "nextActionDueAt", *payload.NextActionDueAt)
		if !ok {
			return
		}
		input.NextActionDueAt = parsed
	}
	if payload.LastContactedAt != nil {
		parsed, ok := parseOptionalTimestamp(w, "lastContactedAt", *payload.LastContactedAt)
		if !ok {
			return
		}
		input.LastContactedAt = parsed
	}

	updated, err := h.inquiries.Update(r.Context(), r.PathValue("id"), input, user.Email)
	if errors.Is(err, inquiry.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", "The requested inquiry was not found.")
		return
	}
	if err != nil {
		slog.Error("admin inquiry update failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "The request could not be completed.")
		return
	}
	writeData(w, http.StatusOK, inquiryToJSON(updated))
}

type activityPayload struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

func (h *adminHandlers) addActivity(w http.ResponseWriter, r *http.Request, user admin.User) {
	if !isJSONContentType(r.Header.Get("Content-Type")) {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Send the request as JSON.")
		return
	}

	var payload activityPayload
	if !decodeStrictJSON(w, r, &payload, maxAdminBodyBytes) {
		return
	}

	kind := inquiry.ActivityType(payload.Type)
	if kind != inquiry.ActivityNote && kind != inquiry.ActivityContact {
		writeFieldErrors(w, http.StatusBadRequest, "validation_error", "Please correct the highlighted fields.", map[string]string{
			"type": "Use note or contact.",
		})
		return
	}
	body := strings.TrimSpace(payload.Body)
	if body == "" {
		writeFieldErrors(w, http.StatusBadRequest, "validation_error", "Please correct the highlighted fields.", map[string]string{
			"body": "Enter a note.",
		})
		return
	}

	activity, err := h.inquiries.AddActivity(r.Context(), inquiry.Activity{
		InquiryID: r.PathValue("id"),
		Type:      kind,
		Body:      body,
		Actor:     user.Email,
	})
	if errors.Is(err, inquiry.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not_found", "The requested inquiry was not found.")
		return
	}
	if err != nil {
		slog.Error("admin inquiry activity create failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "The request could not be completed.")
		return
	}
	writeData(w, http.StatusCreated, activityToJSON(activity))
}

func parseOptionalTimestamp(w http.ResponseWriter, field, value string) (*time.Time, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, true
	}
	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		writeFieldErrors(w, http.StatusBadRequest, "validation_error", "Please correct the highlighted fields.", map[string]string{
			field: "Use an RFC 3339 timestamp.",
		})
		return nil, false
	}
	utc := parsed.UTC()
	return &utc, true
}

func decodeStrictJSON(w http.ResponseWriter, r *http.Request, target any, maxBytes int64) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			writeError(w, http.StatusRequestEntityTooLarge, "validation_error", "The request is too large.")
			return false
		}
		writeError(w, http.StatusBadRequest, "validation_error", "The request could not be read.")
		return false
	}
	if err := decoder.Decode(new(struct{})); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "validation_error", "The request could not be read.")
		return false
	}
	return true
}

func inquiryToJSON(value inquiry.Inquiry) map[string]any {
	product := map[string]any{
		"slug":  value.ProductSlug,
		"model": value.ProductModel,
		"name":  value.ProductName,
	}
	if value.ProductID != nil {
		product["id"] = *value.ProductID
	}

	return map[string]any{
		"id":     value.ID,
		"status": string(value.Status),
		"contact": map[string]any{
			"name":    value.ContactName,
			"email":   value.ContactEmail,
			"company": value.ContactCompany,
			"phone":   value.ContactPhone,
		},
		"destinationCountry":   value.DestinationCountry,
		"product":              product,
		"quantity":             value.Quantity,
		"customizationRequest": value.CustomizationRequest,
		"sourcePath":           value.SourcePath,
		"privacyNoticeVersion": value.PrivacyNoticeVersion,
		"privacyConsentedAt":   formatTimestamp(value.PrivacyConsentedAt),
		"owner":                value.Owner,
		"lastContactedAt":      formatOptionalTimestamp(value.LastContactedAt),
		"nextAction":           value.NextAction,
		"nextActionDueAt":      formatOptionalTimestamp(value.NextActionDueAt),
		"lossReason":           value.LossReason,
		"createdAt":            formatTimestamp(value.CreatedAt),
		"updatedAt":            formatTimestamp(value.UpdatedAt),
	}
}

func activitiesToJSON(activities []inquiry.Activity) []map[string]any {
	items := make([]map[string]any, 0, len(activities))
	for _, activity := range activities {
		items = append(items, activityToJSON(activity))
	}
	return items
}

func activityToJSON(activity inquiry.Activity) map[string]any {
	return map[string]any{
		"id":        activity.ID,
		"type":      string(activity.Type),
		"body":      activity.Body,
		"actor":     activity.Actor,
		"createdAt": formatTimestamp(activity.CreatedAt),
	}
}

func formatTimestamp(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

func formatOptionalTimestamp(value *time.Time) any {
	if value == nil {
		return nil
	}
	return formatTimestamp(*value)
}
