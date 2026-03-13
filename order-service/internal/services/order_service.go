package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/osamah22/nazim/order_service/internal/models"
	"gorm.io/gorm"
)

var ErrOrderNotFound = errors.New("order was not found")

type OrderService struct {
	DB *gorm.DB
}

func (svc *OrderService) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	tx := svc.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if len(order.Items) == 0 {
		tx.Rollback()
		return nil, errors.New("order must have at least one item")
	}

	var total int64
	for i := range order.Items {
		item := &order.Items[i]

		if item.Quantity < 1 {
			tx.Rollback()
			return nil, errors.New("quantity must be at least 1")
		}

		if item.PriceInCents < 0 {
			tx.Rollback()
			return nil, errors.New("price cannot be negative")
		}

		total += int64(item.Quantity) * item.PriceInCents
	}

	order.Total = total

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return order, tx.Commit().Error
}

// GetOrder retrieves an order by ID with items
func (svc *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	tx := svc.DB.WithContext(ctx).Preload("Items").First(&order, "id = ?", id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrOrderNotFound
		}
		return nil, tx.Error
	}
	return &order, nil
}

// ListOrders returns all orders with items
func (svc *OrderService) ListOrders(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	tx := svc.DB.WithContext(ctx).Preload("Items").Find(&orders)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return orders, nil
}

func (svc *OrderService) UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus) (*models.Order, error) {

	var order models.Order
	if err := svc.DB.WithContext(ctx).First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if order.Status == models.StatusCompleted || order.Status == models.StatusCancelled {
		return nil, errors.New("order already finalized")
	}

	order.Status = status

	if err := svc.DB.WithContext(ctx).Save(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// DeleteOrder removes an order by ID
func (svc *OrderService) DeleteOrder(ctx context.Context, id uuid.UUID) error {
	tx := svc.DB.WithContext(ctx).Delete(&models.Order{}, "id = ?", id)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return ErrOrderNotFound
	}
	return nil
}
