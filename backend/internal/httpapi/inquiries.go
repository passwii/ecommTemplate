package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/passwii/ecommTemple/backend/internal/inquiry"
)

const (
	maxInquiryBodyBytes = 64 * 1024
	inquiryPrivacyCode  = "inquiry-privacy-v1"
)

type inquiryContactPayload struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
	Phone   string `json:"phone"`
}

type inquiryPrivacyPayload struct {
	Accepted      bool   `json:"accepted"`
	NoticeVersion string `json:"noticeVersion"`
}

type inquiryPayload struct {
	Contact              inquiryContactPayload `json:"contact"`
	DestinationCountry   string                `json:"destinationCountry"`
	ProductID            string                `json:"productId"`
	Quantity             *int                  `json:"quantity"`
	CustomizationRequest string                `json:"customizationRequest"`
	PrivacyConsent       inquiryPrivacyPayload `json:"privacyConsent"`
	SourcePath           string                `json:"sourcePath"`
	RequestID            string                `json:"requestId"`
	// ReferenceCode is a honeypot field: real visitors never fill it. The name is
	// deliberately non-semantic so browser autofill does not trip it.
	ReferenceCode string `json:"referenceCode"`
}

type inquiryHandlers struct {
	service          *inquiry.Service
	limiter          *inquiry.RateLimiter
	trustProxyHeader bool
}

func registerInquiryRoutes(mux *http.ServeMux, service *inquiry.Service, limiter *inquiry.RateLimiter, trustProxyHeader bool) {
	handlers := &inquiryHandlers{service: service, limiter: limiter, trustProxyHeader: trustProxyHeader}
	mux.HandleFunc("POST /api/v1/inquiries", handlers.create)
}

func (h *inquiryHandlers) create(w http.ResponseWriter, r *http.Request) {
	if !isJSONContentType(r.Header.Get("Content-Type")) {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Send the request as JSON.")
		return
	}

	if allowed, _ := h.limiter.Allow(clientKey(r, h.trustProxyHeader)); !allowed {
		slog.Warn("inquiry submission rate limited")
		writeError(w, http.StatusTooManyRequests, "rate_limited", "Too many submissions. Please try again later.")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxInquiryBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var payload inquiryPayload
	if err := decoder.Decode(&payload); err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			writeError(w, http.StatusRequestEntityTooLarge, "validation_error", "The submission is too large.")
			return
		}
		writeError(w, http.StatusBadRequest, "validation_error", "The submission could not be read.")
		return
	}
	if err := decoder.Decode(new(struct{})); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "validation_error", "The submission could not be read.")
		return
	}

	input := inquiry.CreateInput{
		ContactName:          payload.Contact.Name,
		ContactEmail:         payload.Contact.Email,
		ContactCompany:       payload.Contact.Company,
		ContactPhone:         payload.Contact.Phone,
		DestinationCountry:   payload.DestinationCountry,
		ProductID:            payload.ProductID,
		Quantity:             payload.Quantity,
		CustomizationRequest: payload.CustomizationRequest,
		PrivacyAccepted:      payload.PrivacyConsent.Accepted,
		PrivacyNoticeVersion: payload.PrivacyConsent.NoticeVersion,
		SourcePath:           payload.SourcePath,
		RequestID:            payload.RequestID,
		Honeypot:             payload.ReferenceCode,
	}

	result, fieldErrors, err := h.service.Create(r.Context(), input)
	if fieldErrors != nil {
		writeFieldErrors(w, http.StatusBadRequest, "validation_error", "Please correct the highlighted fields.", fieldErrors)
		return
	}
	if errors.Is(err, inquiry.ErrProductNotAvailable) {
		writeError(w, http.StatusNotFound, "not_found", "The requested product is not available.")
		return
	}
	if err != nil {
		slog.Error("inquiry submission failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "The request could not be completed.")
		return
	}

	status := http.StatusCreated
	if !result.Created {
		status = http.StatusOK
	}
	writeData(w, status, map[string]any{
		"id":            result.Inquiry.ID,
		"status":        "received",
		"receivedAt":    result.Inquiry.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		"privacyNotice": inquiryPrivacyCode,
	})
}

func isJSONContentType(value string) bool {
	if value == "" {
		return false
	}
	mediaType := strings.TrimSpace(strings.Split(value, ";")[0])
	return strings.EqualFold(mediaType, "application/json")
}

// clientKey identifies the submitting client for rate limiting.
//
// By default it uses the transport peer address and deliberately ignores
// forwarding headers, which a client could spoof. When the API runs behind a
// trusted reverse proxy that always overwrites X-Forwarded-For, enabling
// trustProxyHeader makes the limiter key on the real visitor instead of the
// proxy, so one noisy visitor cannot exhaust everyone's quota.
func clientKey(r *http.Request, trustProxyHeader bool) string {
	if trustProxyHeader {
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			first := strings.TrimSpace(strings.Split(forwarded, ",")[0])
			if first != "" {
				return first
			}
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
