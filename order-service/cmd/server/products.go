package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/osamah22/nazim/order_service/internal/dtos"
	"github.com/osamah22/nazim/order_service/internal/models"
	"github.com/osamah22/nazim/order_service/internal/services"
)

type productHandler struct {
	svc *services.ProductService
}

// ListProducts godoc
// @Summary List all products
// @Description Get all available products
// @Tags products
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (h *productHandler) list(c *gin.Context) {
	products, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary Get a product
// @Description Get product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (h *productHandler) get(c *gin.Context) {
	queryId := c.Param("id")

	id, err := uuid.Parse(queryId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query id"})
		return
	}

	product, err := h.svc.Find(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary Create a product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body dtos.CreateProductRequest true "Product data"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products [post]
func (h *productHandler) create(c *gin.Context) {
	req := dtos.CreateProductRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &models.Product{
		Name:           req.Name,
		AvailableStock: req.AvailableStocks,
		PriceInCents:   int64(req.PriceInCents),
	}

	product, err := h.svc.AddProduct(c.Request.Context(), product)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Location", fmt.Sprintf("/products/%s", product.ID))
	c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body dtos.UpdateProductRequest true "Product data"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [put]
func (h *productHandler) update(c *gin.Context) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}
	req := dtos.UpdateProductRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := &models.Product{
		ID:             productID,
		Name:           req.Name,
		AvailableStock: req.AvailableStocks,
		PriceInCents:   int64(req.PriceInCents),
	}

	product, err = h.svc.UpdateProduct(c.Request.Context(), product)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Location", fmt.Sprintf("/products/%s", product.ID))
	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [delete]
func (h *productHandler) delete(c *gin.Context) {
	// Parse ID from path
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	// Call service
	err = h.svc.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Success → 204 No Content
	c.Status(http.StatusNoContent)
}
