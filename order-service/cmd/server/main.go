package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/osamah22/nazim/order_service/internal/models"
	"github.com/osamah22/nazim/order_service/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	port := "8080"
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
	// Read connection string from env (set in docker-compose)
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		// fallback for local dev
		dsn = "host=localhost user=postgres password=postgres dbname=coffee port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// auto migrate your models
	if err := db.AutoMigrate(&models.Product{}, &models.Order{}, &models.LineItem{}); err != nil {
		log.Fatal("auto migration failed:", err)
	}

	return db
}
