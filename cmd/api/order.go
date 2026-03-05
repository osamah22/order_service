package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/osamah22/order_service/internal/dtos"
	"github.com/osamah22/order_service/internal/models"
	"github.com/osamah22/order_service/internal/services"
)

type orderHandler struct {
	orderSvc   *services.OrderService
	productSvc *services.ProductService
}

// ListOrders godoc
// @Summary List all orders
// @Tags orders
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} map[string]string
// @Router /orders [get]
func (h *orderHandler) list(c *gin.Context) {
	orders, err := h.orderSvc.ListOrders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var response []dtos.OrderResponse
	for _, o := range orders {
		response = append(response, dtos.ToOrderResponse(&o))
	}

	c.JSON(http.StatusOK, response)
}

// GetOrder godoc
// @Summary Get an order by ID
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [get]
func (h *orderHandler) get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := h.orderSvc.GetOrder(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dtos.ToOrderResponse(order)
	c.JSON(http.StatusOK, response)
}

// CreateOrder godoc
// @Summary Create a new order
// @Tags orders
// @Accept json
// @Produce json
// @Param order body dtos.CreateOrderRequest true "Order data"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Router /orders [post]
func (h *orderHandler) create(c *gin.Context) {
	var req dtos.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert request to models
	order := &models.Order{}
	for _, item := range req.Items {
		productID, err := uuid.Parse(item.ProductID)
		if err != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"error": fmt.Sprintf("invalid product ID: %s", item.ProductID)})
			return
		}

		product, err := h.productSvc.Find(c.Request.Context(), productID)

		if isNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		order.Items = append(order.Items, models.LineItem{
			ProductID:    productID,
			PriceInCents: product.PriceInCents,
			Quantity:     item.Quantity,
		})
	}

	order, err := h.orderSvc.CreateOrder(c.Request.Context(), order)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Location", fmt.Sprintf("/orders/%s", order.ID))
	response := dtos.ToOrderResponse(order)
	c.JSON(http.StatusCreated, response)
}

// DeleteOrder godoc
// @Summary Delete an order by ID
// @Tags orders
// @Param id path string true "Order ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [delete]
func (h *orderHandler) delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	err = h.orderSvc.DeleteOrder(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
