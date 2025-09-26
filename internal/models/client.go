package models

import (
	"time"

	"github.com/google/uuid"
)

type ClientType string

const (
	ClientTypePublic       ClientType = "public"
	ClientTypeConfidential ClientType = "confidential"
)

type Client struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primaryKey"`
	ClientID               string     `gorm:"type:varchar(128);uniqueIndex;not null"`
	Name                   string     `gorm:"type:varchar(255);not null"`
	Type                   ClientType `gorm:"type:varchar(32);not null"`
	SecretHash             *string    `gorm:"type:text"`
	RedirectURIs           []string   `gorm:"type:jsonb;not null;default:'[]';serializer:json"`
	PostLogoutRedirectURIs []string   `gorm:"type:jsonb;not null;default:'[]';serializer:json"`
	Scopes                 []string   `gorm:"type:jsonb;not null;default:'[]';serializer:json"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
