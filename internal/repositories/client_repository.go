package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mohammadhprp/passport/internal/models"
)

var (
	ErrClientNotFound      = errors.New("client not found")
	ErrClientAlreadyExists = errors.New("client already exists")
)

type ClientRepository interface {
	Create(ctx context.Context, client *models.Client) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Client, error)
	GetByClientID(ctx context.Context, clientID string) (*models.Client, error)
	List(ctx context.Context, offset, limit int) ([]models.Client, error)
	UpdateSecret(ctx context.Context, id uuid.UUID, secretHash *string) error
}

type clientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(ctx context.Context, client *models.Client) error {
	err := r.db.WithContext(ctx).Create(client).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrClientAlreadyExists
		}
		return err
	}
	return nil
}

func (r *clientRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	var client models.Client
	err := r.db.WithContext(ctx).First(&client, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, err
	}
	return &client, nil
}

func (r *clientRepository) GetByClientID(ctx context.Context, clientID string) (*models.Client, error) {
	var client models.Client
	err := r.db.WithContext(ctx).First(&client, "client_id = ?", clientID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, err
	}
	return &client, nil
}

func (r *clientRepository) List(ctx context.Context, offset, limit int) ([]models.Client, error) {
	var clients []models.Client

	query := r.db.WithContext(ctx).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&clients).Error; err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *clientRepository) UpdateSecret(ctx context.Context, id uuid.UUID, secretHash *string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Client{}).
		Where("id = ?", id).
		Update("secret_hash", secretHash)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrClientNotFound
	}

	return nil
}
