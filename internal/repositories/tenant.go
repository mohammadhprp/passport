package repositories

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/mohammadhprp/passport/internal/models"
	"gorm.io/gorm"
)

// TenantRepository exposes persistence helpers for tenant entities.
type TenantRepository struct {
	db *gorm.DB
}

// NewTenantRepository builds a repository bound to the provided database handle.
func NewTenantRepository(db *gorm.DB) *TenantRepository {
	if err := db.AutoMigrate(&models.Tenant{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	return &TenantRepository{db: db}
}

// Create persists a new tenant.
func (r *TenantRepository) Create(ctx context.Context, tenant *models.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

// List returns all tenants.
func (r *TenantRepository) List(ctx context.Context) ([]models.Tenant, error) {
	var tenants []models.Tenant
	if err := r.db.WithContext(ctx).Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

// GetByID returns the tenant matching the given id.
func (r *TenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := r.db.WithContext(ctx).First(&tenant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// Update saves the provided tenant state.
func (r *TenantRepository) Update(ctx context.Context, tenant *models.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

// Delete removes the tenant with the given id.
func (r *TenantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Tenant{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
