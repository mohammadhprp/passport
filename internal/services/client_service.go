package services

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"github.com/mohammadhprp/passport/internal/models"
	"github.com/mohammadhprp/passport/internal/repositories"
	"github.com/mohammadhprp/passport/internal/utils"
)

const clientSecretEntropyBytes = 32

var (
	ErrClientNameRequired  = errors.New("client name must not be empty")
	ErrInvalidClientType   = errors.New("invalid client type")
	ErrMissingRedirectURIs = errors.New("at least one redirect uri is required")
	ErrInvalidRedirectURI  = errors.New("redirect uri must be absolute")
	ErrClientIsPublic      = errors.New("client is public and does not have a secret")
)

type CreateClientParams struct {
	Name                   string
	Type                   models.ClientType
	RedirectURIs           []string
	PostLogoutRedirectURIs []string
	Scopes                 []string
}

type CreateClientResult struct {
	Client      *models.Client
	PlainSecret *string
}

type RotateClientSecretResult struct {
	Client      *models.Client
	PlainSecret string
}

type ListClientsFilter struct {
	Offset int
	Limit  int
}

type ClientService interface {
	CreateClient(ctx context.Context, params CreateClientParams) (*CreateClientResult, error)
	GetClient(ctx context.Context, id uuid.UUID) (*models.Client, error)
	GetClientByClientID(ctx context.Context, clientID string) (*models.Client, error)
	ListClients(ctx context.Context, filter ListClientsFilter) ([]models.Client, error)
	RotateClientSecret(ctx context.Context, id uuid.UUID) (*RotateClientSecretResult, error)
}

type clientService struct {
	repo repositories.ClientRepository
}

func NewClientService(repo repositories.ClientRepository) ClientService {
	return &clientService{repo: repo}
}

func (s *clientService) CreateClient(ctx context.Context, params CreateClientParams) (*CreateClientResult, error) {
	if strings.TrimSpace(params.Name) == "" {
		return nil, ErrClientNameRequired
	}

	clientType := params.Type
	if clientType == "" {
		clientType = models.ClientTypeConfidential
	}
	if !isValidClientType(clientType) {
		return nil, ErrInvalidClientType
	}

	redirectURIs, err := normalizeURIList(params.RedirectURIs, true)
	if err != nil {
		return nil, err
	}

	postLogoutURIs, err := normalizeURIList(params.PostLogoutRedirectURIs, false)
	if err != nil {
		return nil, err
	}

	var secretHash *string
	var plainSecret *string

	if clientType == models.ClientTypeConfidential {
		secretValue, err := utils.GenerateRandomToken(clientSecretEntropyBytes)
		if err != nil {
			return nil, err
		}
		hash, err := utils.HashSensitiveValue(secretValue)
		if err != nil {
			return nil, err
		}

		secretHash = &hash
		plainSecret = &secretValue
	}

	client := &models.Client{
		ID:                     uuid.New(),
		ClientID:               uuid.NewString(),
		Name:                   strings.TrimSpace(params.Name),
		Type:                   clientType,
		SecretHash:             secretHash,
		RedirectURIs:           redirectURIs,
		PostLogoutRedirectURIs: postLogoutURIs,
		Scopes:                 sanitizeScopes(params.Scopes),
	}

	if client.RedirectURIs == nil {
		client.RedirectURIs = []string{}
	}
	if client.PostLogoutRedirectURIs == nil {
		client.PostLogoutRedirectURIs = []string{}
	}
	if client.Scopes == nil {
		client.Scopes = []string{}
	}

	if err := s.repo.Create(ctx, client); err != nil {
		return nil, err
	}

	return &CreateClientResult{Client: client, PlainSecret: plainSecret}, nil
}

func (s *clientService) GetClient(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *clientService) GetClientByClientID(ctx context.Context, clientID string) (*models.Client, error) {
	return s.repo.GetByClientID(ctx, clientID)
}

func (s *clientService) ListClients(ctx context.Context, filter ListClientsFilter) ([]models.Client, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	return s.repo.List(ctx, offset, limit)
}

func (s *clientService) RotateClientSecret(ctx context.Context, id uuid.UUID) (*RotateClientSecretResult, error) {
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if client.Type != models.ClientTypeConfidential {
		return nil, ErrClientIsPublic
	}

	secretValue, err := utils.GenerateRandomToken(clientSecretEntropyBytes)
	if err != nil {
		return nil, err
	}

	hash, err := utils.HashSensitiveValue(secretValue)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpdateSecret(ctx, id, &hash); err != nil {
		return nil, err
	}

	client.SecretHash = &hash

	return &RotateClientSecretResult{Client: client, PlainSecret: secretValue}, nil
}

func isValidClientType(clientType models.ClientType) bool {
	switch clientType {
	case models.ClientTypePublic, models.ClientTypeConfidential:
		return true
	default:
		return false
	}
}

func normalizeURIList(values []string, required bool) ([]string, error) {
	if len(values) == 0 {
		if required {
			return nil, ErrMissingRedirectURIs
		}
		return nil, nil
	}

	seen := make(map[string]struct{})
	normalized := make([]string, 0, len(values))

	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}

		parsed, err := url.Parse(trimmed)
		if err != nil || !parsed.IsAbs() {
			return nil, ErrInvalidRedirectURI
		}

		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	if required && len(normalized) == 0 {
		return nil, ErrMissingRedirectURIs
	}

	if len(normalized) == 0 {
		return nil, nil
	}

	return normalized, nil
}

func sanitizeScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return nil
	}

	seen := make(map[string]struct{})
	clean := make([]string, 0, len(scopes))

	for _, scope := range scopes {
		trimmed := strings.TrimSpace(scope)
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		clean = append(clean, trimmed)
	}

	if len(clean) == 0 {
		return nil
	}

	return clean
}
