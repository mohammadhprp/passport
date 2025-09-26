package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mohammadhprp/passport/internal/models"
	"github.com/mohammadhprp/passport/internal/repositories"
	"gorm.io/gorm"
)

// ErrTenantNotFound describes a missing tenant lookup.
var ErrTenantNotFound = errors.New("tenant not found")

// TenantService orchestrates business logic for tenant workflows.
type TenantService struct {
	repo *repositories.TenantRepository
}

// NewTenantService builds a tenant service using the provided repository.
func NewTenantService(repo *repositories.TenantRepository) *TenantService {
	return &TenantService{repo: repo}
}

// CreateTenant persists a new tenant with the given name.
func (s *TenantService) CreateTenant(ctx context.Context, name string) (*models.Tenant, error) {
	tenant := &models.Tenant{Name: name}
	if err := s.repo.Create(ctx, tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

// ListTenants returns every tenant.
func (s *TenantService) ListTenants(ctx context.Context) ([]models.Tenant, error) {
	return s.repo.List(ctx)
}

// GetTenant fetches a tenant by id.
func (s *TenantService) GetTenant(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}
	return tenant, nil
}

// UpdateTenant renames an existing tenant.
func (s *TenantService) UpdateTenant(ctx context.Context, id uuid.UUID, name string) (*models.Tenant, error) {
	tenant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, err
	}

	tenant.Name = name

	if err := s.repo.Update(ctx, tenant); err != nil {
		return nil, err
	}

	return tenant, nil
}

// DeleteTenant removes a tenant.
func (s *TenantService) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTenantNotFound
		}
		return err
	}
	return nil
}
