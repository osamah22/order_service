package main

import (
	"github.com/gin-gonic/gin"
	"github.com/osamah22/nazim/auth-service/internal/services"
)

type authHandler struct {
	authService *services.AuthService
}

func (h *authHandler) register(c *gin.Context) {

}

func (h *authHandler) login(c *gin.Context) {

}

func (h *authHandler) refreshToken(c *gin.Context) {

}
