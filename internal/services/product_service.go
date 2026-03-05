package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/osamah22/order_service/internal/models"
	"gorm.io/gorm"
)

var ErrProductNotFound = errors.New("product was not found")

type ProductService struct {
	DB *gorm.DB
}

func (svc *ProductService) ListAll(
	ctx context.Context,
) ([]models.Product, error) {
	var products []models.Product

	if err := svc.DB.WithContext(ctx).
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (svc *ProductService) Find(
	ctx context.Context,
	id uuid.UUID,
) (*models.Product, error) {
	var product models.Product
	tx := svc.DB.WithContext(ctx).
		Where("id = ?", id).
		Find(&product)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, ErrProductNotFound
	}
	return &product, nil
}
func (svc *ProductService) AddProduct(
	ctx context.Context,
	product *models.Product,
) (*models.Product, error) {
	tx := svc.DB.WithContext(ctx).Create(product)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return product, nil
}

func (svc *ProductService) UpdateProduct(
	ctx context.Context,
	product *models.Product,
) (*models.Product, error) {
	result := svc.DB.WithContext(ctx).
		Model(&models.Product{}).
		Where("id = ?", product.ID).
		Updates(product)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrProductNotFound
	}

	return product, nil
}

func (svc *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	result := svc.DB.WithContext(ctx).
		Delete(&models.Product{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrProductNotFound
	}

	return nil
}
