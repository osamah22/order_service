package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username       string    `gorm:"uniqueIndex;not null"`
	GivenName      string    `gorm:"not null"`
	Email          string    `gorm:"uniqueIndex;not null"`
	HashedPassword string    `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
