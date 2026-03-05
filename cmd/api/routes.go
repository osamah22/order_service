package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/osamah22/order_service/docs" // <-- important: generated docs package
)

// @title Product API
// @version 1.0
// @description Simple product service
// @host localhost:7070
// @BasePath /
func addRoutes(router *gin.Engine,
	productHandler productHandler,
	orderHandler orderHandler) {

	// register product routes
	products := router.Group("/products")
	products.GET("/", productHandler.list)
	products.GET("/:id", productHandler.get)
	products.POST("/", productHandler.create)
	products.PUT("/:id", productHandler.update)
	products.DELETE("/:id", productHandler.delete)

	// register order routes
	orders := router.Group("/orders")
	orders.GET("/", orderHandler.list)
	orders.GET("/:id", orderHandler.get)
	orders.POST("/", orderHandler.create)
	orders.DELETE("/:id", orderHandler.delete)

	// setup swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
