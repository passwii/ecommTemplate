package catalog

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestProductMarshalsEmptyListsAsArrays(t *testing.T) {
	product := Product{
		ID:             "grill-01",
		Slug:           "ridge-4-burner",
		Name:           "Ridge 4-Burner",
		Images:         []ProductImage{},
		Specifications: []ProductSpecification{},
		Applications:   []string{},
		Resources:      []ProductResource{},
		ContentStatus:  "sample",
	}

	encoded, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	body := string(encoded)
	for _, field := range []string{`"images":[]`, `"specifications":[]`, `"applications":[]`, `"resources":[]`} {
		if !strings.Contains(body, field) {
			t.Errorf("marshaled product missing %s; got %s", field, body)
		}
	}
}

func TestCompanyContentMarshalsEmptyListsAsArrays(t *testing.T) {
	content := CompanyContent{
		Credentials:   []CompanyCredential{},
		CaseStudies:   []CompanyCaseStudy{},
		FAQs:          []CompanyFAQ{},
		ContentStatus: "sample",
	}

	encoded, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	body := string(encoded)
	for _, field := range []string{`"credentials":[]`, `"caseStudies":[]`, `"faqs":[]`} {
		if !strings.Contains(body, field) {
			t.Errorf("marshaled company content missing %s; got %s", field, body)
		}
	}
}
