package dtos

import "github.com/osamah22/nazim/auth-service/internal/models"

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=30"`
	GivenName string `json:"given_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=30"`
	GivenName string `json:"given_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type ChangePasswordRequest struct {
	Password string `json:"password"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	GivenName string `json:"given_name"`
	Email     string `json:"email"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ToUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		GivenName: user.GivenName,
		Email:     user.Email,
	}
}
