package inquiry

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/passwii/ecommTemple/backend/internal/catalog"
)

// ProductLookup resolves a published product by its stable id.
type ProductLookup interface {
	GetPublishedProductByID(ctx context.Context, id string) (catalog.Product, error)
}

// CreateResult is the outcome of an inquiry submission.
type CreateResult struct {
	Inquiry Inquiry
	// Created is false when an idempotent replay returned an existing inquiry.
	Created bool
}

// Service validates and persists visitor inquiries.
type Service struct {
	repo     Repository
	products ProductLookup
	now      func() time.Time
	newID    func() (string, error)
}

// NewService builds an inquiry service.
func NewService(repo Repository, products ProductLookup) *Service {
	return &Service{
		repo:     repo,
		products: products,
		now:      time.Now,
		newID:    generateID,
	}
}

// Create validates and stores a visitor inquiry. It returns field errors for
// invalid input and ErrProductNotAvailable when a submitted product is unknown
// or unpublished.
func (s *Service) Create(ctx context.Context, input CreateInput) (CreateResult, map[string]string, error) {
	normalized, fieldErrors := Validate(input)
	if fieldErrors != nil {
		return CreateResult{}, fieldErrors, nil
	}

	if normalized.RequestID != "" {
		existing, err := s.repo.FindByRequestID(ctx, normalized.RequestID)
		if err == nil {
			return CreateResult{Inquiry: existing, Created: false}, nil, nil
		}
		if !errors.Is(err, ErrNotFound) {
			return CreateResult{}, nil, fmt.Errorf("check idempotency key: %w", err)
		}
	}

	value := Inquiry{
		ContactName:          normalized.ContactName,
		ContactEmail:         normalized.ContactEmail,
		ContactCompany:       normalized.ContactCompany,
		ContactPhone:         normalized.ContactPhone,
		DestinationCountry:   normalized.DestinationCountry,
		Quantity:             normalized.Quantity,
		CustomizationRequest: normalized.CustomizationRequest,
		PrivacyNoticeVersion: normalized.PrivacyNoticeVersion,
		PrivacyConsentedAt:   s.now().UTC(),
		SourcePath:           normalized.SourcePath,
		RequestID:            normalized.RequestID,
		Status:               StatusNew,
	}

	if normalized.ProductID != "" {
		product, err := s.products.GetPublishedProductByID(ctx, normalized.ProductID)
		if errors.Is(err, catalog.ErrNotFound) {
			return CreateResult{}, nil, ErrProductNotAvailable
		}
		if err != nil {
			return CreateResult{}, nil, fmt.Errorf("resolve product: %w", err)
		}
		snapshot := snapshotFor(product)
		value.ProductID = &snapshot.ID
		value.ProductSlug = snapshot.Slug
		value.ProductModel = snapshot.Model
		value.ProductName = snapshot.Name
	}

	id, err := s.newID()
	if err != nil {
		return CreateResult{}, nil, fmt.Errorf("generate inquiry id: %w", err)
	}
	value.ID = id

	stored, err := s.repo.Create(ctx, value)
	if err != nil {
		// A concurrent request with the same idempotency key won the race.
		if normalized.RequestID != "" {
			if existing, lookupErr := s.repo.FindByRequestID(ctx, normalized.RequestID); lookupErr == nil {
				return CreateResult{Inquiry: existing, Created: false}, nil, nil
			}
		}
		return CreateResult{}, nil, fmt.Errorf("store inquiry: %w", err)
	}
	return CreateResult{Inquiry: stored, Created: true}, nil, nil
}

func snapshotFor(product catalog.Product) ProductSnapshot {
	model := product.Model
	if model == "" {
		model = product.Name
	}
	return ProductSnapshot{
		ID:    product.ID,
		Slug:  product.Slug,
		Model: model,
		Name:  product.Name,
	}
}

func generateID() (string, error) {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return "inq_" + hex.EncodeToString(buffer), nil
}
