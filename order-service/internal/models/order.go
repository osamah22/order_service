package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusCompleted OrderStatus = "completed"
	StatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Items     []LineItem `gorm:"foreignKey:OrderID"`
	Total     int64
	Status    OrderStatus `gorm:"default:'pending'"`
	CreatedAt time.Time
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	o.CreatedAt = time.Now()
	return nil
}

type LineItem struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrderID      uuid.UUID `gorm:"index:idx_order_product,unique"`
	ProductID    uuid.UUID `gorm:"index:idx_order_product,unique"`
	Quantity     int
	PriceInCents int64
}

func (li *LineItem) BeforeCreate(tx *gorm.DB) (err error) {
	if li.ID == uuid.Nil {
		li.ID = uuid.New()
	}
	if li.Quantity < 1 {
		return errors.New("quantity must be greater than zero")
	}
	return nil
}
