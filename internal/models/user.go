package models

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusPending  UserStatus = "pending"
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
)

type User struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Email         string     `gorm:"type:varchar(320);uniqueIndex;not null"`
	EmailVerified bool       `gorm:"not null;default:false"`
	PasswordHash  string     `gorm:"type:text;not null"`
	MFAFactors    []string   `gorm:"type:jsonb;not null;default:'[]';serializer:json"`
	Status        UserStatus `gorm:"type:varchar(32);not null;default:'pending'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
