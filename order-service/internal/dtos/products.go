package dtos

type CreateProductRequest struct {
	Name            string `json:"name" binding:"required"`
	AvailableStocks int    `json:"available_stocks" binding:"required,gte=0"`
	PriceInCents    int    `json:"price_in_cents" binding:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name            string `json:"name" binding:"required"`
	AvailableStocks int    `json:"available_stocks" binding:"required,gte=0"`
	PriceInCents    int    `json:"price_in_cents" binding:"required,gte=0"`
}
