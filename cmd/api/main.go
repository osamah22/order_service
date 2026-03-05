package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/osamah22/order_service/internal/models"
	"github.com/osamah22/order_service/internal/services"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	port := "7070"
	db := setupDatabase()
	productService := services.ProductService{
		DB: db,
	}
	orderService := services.OrderService{
		DB: db,
	}

	router := gin.New()
	addRoutes(router,
		productHandler{svc: &productService},
		orderHandler{orderSvc: &orderService, productSvc: &productService})

	fmt.Printf("starting service on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func setupDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.Product{}, &models.Order{}, &models.LineItem{})
	return db
}
