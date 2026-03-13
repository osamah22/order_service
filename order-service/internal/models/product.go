package models

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name           string
	PriceInCents   int64
	AvailableStock int
}

// this function automatically runs when inserting new record to the database
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	// create uuid if not initialized
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func ValidateProduct(p *Product) error {
	if p.PriceInCents < 0 {
		return errors.New("price cannot be negative")
	}
	if p.AvailableStock < 0 {
		return errors.New("stock cannot be negative")
	}
	return nil
}
