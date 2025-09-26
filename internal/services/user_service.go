package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"

	"github.com/mohammadhprp/passport/internal/models"
	"github.com/mohammadhprp/passport/internal/repositories"
)

const (
	argonTime    uint32 = 3
	argonMemory  uint32 = 64 * 1024
	argonThreads uint8  = 2
	argonKeyLen  uint32 = 32
	argonSaltLen        = 16
)

var (
	ErrInvalidUserStatus = errors.New("invalid user status")
	ErrInvalidPassword   = errors.New("password must not be empty")
)

type CreateUserParams struct {
	Email         string
	Password      string
	MFAFactors    []string
	Status        models.UserStatus
	EmailVerified bool
}

type ListUsersFilter struct {
	Offset int
	Limit  int
}

type UserService interface {
	CreateUser(ctx context.Context, params CreateUserParams) (*models.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListUsers(ctx context.Context, filter ListUsersFilter) ([]models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, params CreateUserParams) (*models.User, error) {
	if params.Password == "" {
		return nil, ErrInvalidPassword
	}

	status := params.Status
	if status == "" {
		status = models.UserStatusPending
	}
	if !isValidStatus(status) {
		return nil, ErrInvalidUserStatus
	}

	hash, err := hashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:            uuid.New(),
		Email:         params.Email,
		EmailVerified: params.EmailVerified,
		PasswordHash:  hash,
		MFAFactors:    cloneStringSlice(params.MFAFactors),
		Status:        status,
	}

	if user.MFAFactors == nil {
		user.MFAFactors = []string{}
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, filter ListUsersFilter) ([]models.User, error) {
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

func isValidStatus(status models.UserStatus) bool {
	switch status {
	case models.UserStatusPending, models.UserStatusActive, models.UserStatusDisabled:
		return true
	default:
		return false
	}
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", argonMemory, argonTime, argonThreads, b64Salt, b64Hash)
	return encoded, nil
}

func cloneStringSlice(input []string) []string {
	if len(input) == 0 {
		return nil
	}
	out := make([]string, len(input))
	copy(out, input)
	return out
}
