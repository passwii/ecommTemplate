package inquiry

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Field length bounds for visitor-submitted input.
const (
	maxNameLength       = 200
	maxEmailLength      = 320
	maxCompanyLength    = 200
	maxPhoneLength      = 50
	maxCountryLength    = 2
	maxCustomizationLen = 4000
	maxSourcePathLength = 300
	maxNoticeVersionLen = 100
	maxRequestIDLength  = 100
	maxQuantity         = 1_000_000
	minQuantity         = 1
)

var countryCodePattern = regexp.MustCompile(`^[A-Za-z]{2}$`)

// CreateInput is the visitor-submitted inquiry payload. Operational fields
// (status, owner, activities) are never accepted from visitors.
type CreateInput struct {
	ContactName          string
	ContactEmail         string
	ContactCompany       string
	ContactPhone         string
	DestinationCountry   string
	ProductID            string
	Quantity             *int
	CustomizationRequest string
	PrivacyAccepted      bool
	PrivacyNoticeVersion string
	SourcePath           string
	RequestID            string
	Honeypot             string
}

// Validate normalizes and checks visitor input. It returns the normalized input
// and a map of field errors. The map is nil when the input is valid.
func Validate(input CreateInput) (CreateInput, map[string]string) {
	errors := map[string]string{}

	// Honeypot: a real visitor never fills this hidden field.
	if strings.TrimSpace(input.Honeypot) != "" {
		errors["honeypot"] = "This submission was rejected."
	}

	input.ContactName = strings.TrimSpace(input.ContactName)
	switch {
	case input.ContactName == "":
		errors["contactName"] = "Enter your name."
	case utf8.RuneCountInString(input.ContactName) > maxNameLength:
		errors["contactName"] = "This value is too long."
	}

	input.ContactEmail = strings.TrimSpace(input.ContactEmail)
	switch {
	case input.ContactEmail == "":
		errors["contactEmail"] = "Enter your email address."
	case len(input.ContactEmail) > maxEmailLength:
		errors["contactEmail"] = "This value is too long."
	case !isValidEmail(input.ContactEmail):
		errors["contactEmail"] = "Enter a valid email address."
	}

	input.ContactCompany = strings.TrimSpace(input.ContactCompany)
	if utf8.RuneCountInString(input.ContactCompany) > maxCompanyLength {
		errors["contactCompany"] = "This value is too long."
	}

	input.ContactPhone = strings.TrimSpace(input.ContactPhone)
	if utf8.RuneCountInString(input.ContactPhone) > maxPhoneLength {
		errors["contactPhone"] = "This value is too long."
	}

	input.DestinationCountry = strings.ToUpper(strings.TrimSpace(input.DestinationCountry))
	if input.DestinationCountry != "" && !countryCodePattern.MatchString(input.DestinationCountry) {
		errors["destinationCountry"] = "Use a two-letter country code."
	}

	input.ProductID = strings.TrimSpace(input.ProductID)

	if input.Quantity != nil {
		if *input.Quantity < minQuantity || *input.Quantity > maxQuantity {
			errors["quantity"] = "Enter a quantity greater than zero."
		}
	}

	input.CustomizationRequest = strings.TrimSpace(input.CustomizationRequest)
	if utf8.RuneCountInString(input.CustomizationRequest) > maxCustomizationLen {
		errors["customizationRequest"] = "This value is too long."
	}

	input.PrivacyNoticeVersion = strings.TrimSpace(input.PrivacyNoticeVersion)
	switch {
	case !input.PrivacyAccepted:
		errors["privacyConsent"] = "Please accept the privacy notice."
	case input.PrivacyNoticeVersion == "":
		errors["privacyConsent"] = "The privacy notice version is required."
	case utf8.RuneCountInString(input.PrivacyNoticeVersion) > maxNoticeVersionLen:
		errors["privacyConsent"] = "The privacy notice version is invalid."
	}

	input.SourcePath = strings.TrimSpace(input.SourcePath)
	switch {
	case utf8.RuneCountInString(input.SourcePath) > maxSourcePathLength:
		errors["sourcePath"] = "This value is too long."
	case input.SourcePath != "" && !isSafeSourcePath(input.SourcePath):
		errors["sourcePath"] = "Use a path within this storefront."
	}

	input.RequestID = strings.TrimSpace(input.RequestID)
	if utf8.RuneCountInString(input.RequestID) > maxRequestIDLength {
		errors["requestId"] = "This value is too long."
	}

	if len(errors) > 0 {
		return input, errors
	}
	return input, nil
}

func isValidEmail(value string) bool {
	if strings.ContainsAny(value, " \t\r\n") {
		return false
	}
	address, err := mail.ParseAddress(value)
	if err != nil {
		return false
	}
	return address.Address == value && strings.Contains(value, "@")
}

func isSafeSourcePath(value string) bool {
	if !strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") || strings.Contains(value, "\\") {
		return false
	}
	for _, r := range value {
		if r < 0x20 {
			return false
		}
	}
	return true
}
