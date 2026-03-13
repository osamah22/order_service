package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/osamah22/nazim/auth-service/docs" // <-- important: generated docs package
)

// @title Product API
// @version 1.0
// @description Simple product service
// @host localhost:7070
// @BasePath /
func addRoutes(router *gin.Engine) {

	// authentication routes
	auth := router.Group("/auth")

	// setup swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
