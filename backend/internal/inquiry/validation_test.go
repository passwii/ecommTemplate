package inquiry

import (
	"strings"
	"testing"
)

func validInput() CreateInput {
	quantity := 120
	return CreateInput{
		ContactName:          "Alex Morgan",
		ContactEmail:         "buyer@example.com",
		ContactCompany:       "Example Importers",
		ContactPhone:         "+1 555 0100",
		DestinationCountry:   "us",
		ProductID:            "grill-01",
		Quantity:             &quantity,
		CustomizationRequest: "Private-label packaging options",
		PrivacyAccepted:      true,
		PrivacyNoticeVersion: "inquiry-privacy-v1",
		SourcePath:           "/products/ridge-4-burner",
	}
}

func TestValidateAcceptsValidInputAndNormalizes(t *testing.T) {
	normalized, fieldErrors := Validate(validInput())
	if fieldErrors != nil {
		t.Fatalf("fieldErrors = %v, want nil", fieldErrors)
	}
	if normalized.DestinationCountry != "US" {
		t.Errorf("DestinationCountry = %q, want US", normalized.DestinationCountry)
	}
	if normalized.ContactName != "Alex Morgan" {
		t.Errorf("ContactName = %q", normalized.ContactName)
	}
}

func TestValidateTrimsSurroundingWhitespace(t *testing.T) {
	input := validInput()
	input.ContactName = "  Alex Morgan  "
	input.ContactEmail = " buyer@example.com "
	normalized, fieldErrors := Validate(input)
	if fieldErrors != nil {
		t.Fatalf("fieldErrors = %v, want nil", fieldErrors)
	}
	if normalized.ContactName != "Alex Morgan" || normalized.ContactEmail != "buyer@example.com" {
		t.Errorf("not trimmed: %+v", normalized)
	}
}

func TestValidateRejectsInvalidInput(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*CreateInput)
		field  string
	}{
		{"missing name", func(in *CreateInput) { in.ContactName = "   " }, "contactName"},
		{"missing email", func(in *CreateInput) { in.ContactEmail = "" }, "contactEmail"},
		{"invalid email", func(in *CreateInput) { in.ContactEmail = "not-an-email" }, "contactEmail"},
		{"email with space", func(in *CreateInput) { in.ContactEmail = "a b@example.com" }, "contactEmail"},
		{"overlong name", func(in *CreateInput) { in.ContactName = strings.Repeat("a", 201) }, "contactName"},
		{"zero quantity", func(in *CreateInput) { zero := 0; in.Quantity = &zero }, "quantity"},
		{"negative quantity", func(in *CreateInput) { negative := -5; in.Quantity = &negative }, "quantity"},
		{"huge quantity", func(in *CreateInput) { huge := 2_000_000; in.Quantity = &huge }, "quantity"},
		{"privacy not accepted", func(in *CreateInput) { in.PrivacyAccepted = false }, "privacyConsent"},
		{"missing notice version", func(in *CreateInput) { in.PrivacyNoticeVersion = "" }, "privacyConsent"},
		{"bad country code", func(in *CreateInput) { in.DestinationCountry = "United States" }, "destinationCountry"},
		{"source path not absolute", func(in *CreateInput) { in.SourcePath = "products/x" }, "sourcePath"},
		{"source path protocol relative", func(in *CreateInput) { in.SourcePath = "//evil.example.com" }, "sourcePath"},
		{"source path with backslash", func(in *CreateInput) { in.SourcePath = `/a\b` }, "sourcePath"},
		{"source path control char", func(in *CreateInput) { in.SourcePath = "/a\u0000b" }, "sourcePath"},
		{"overlong customization", func(in *CreateInput) { in.CustomizationRequest = strings.Repeat("x", 4001) }, "customizationRequest"},
		{"overlong request id", func(in *CreateInput) { in.RequestID = strings.Repeat("r", 101) }, "requestId"},
		{"honeypot filled", func(in *CreateInput) { in.Honeypot = "http://spam.example.com" }, "honeypot"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			input := validInput()
			tc.mutate(&input)
			_, fieldErrors := Validate(input)
			if fieldErrors == nil {
				t.Fatalf("fieldErrors = nil, want error on %s", tc.field)
			}
			if _, ok := fieldErrors[tc.field]; !ok {
				t.Fatalf("fieldErrors = %v, want key %s", fieldErrors, tc.field)
			}
		})
	}
}

func TestValidateAllowsOptionalFieldsToBeEmpty(t *testing.T) {
	input := CreateInput{
		ContactName:          "Alex",
		ContactEmail:         "alex@example.com",
		PrivacyAccepted:      true,
		PrivacyNoticeVersion: "v1",
	}
	normalized, fieldErrors := Validate(input)
	if fieldErrors != nil {
		t.Fatalf("fieldErrors = %v, want nil", fieldErrors)
	}
	if normalized.Quantity != nil {
		t.Errorf("Quantity = %v, want nil", normalized.Quantity)
	}
	if normalized.ProductID != "" {
		t.Errorf("ProductID = %q, want empty", normalized.ProductID)
	}
}
