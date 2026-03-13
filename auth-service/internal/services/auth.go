package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/osamah22/nazim/auth-service/internal/dtos"
	"github.com/osamah22/nazim/auth-service/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		DB: db,
	}
}

func (s *AuthService) Register(req dtos.CreateUserRequest) (*models.User, error) {

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:       req.Username,
		GivenName:      req.GivenName,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) Login(email string, password string) (*dtos.TokenResponse, error) {

	var user models.User

	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := s.verifyPassword(user.HashedPassword, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.createRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dtos.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (s *AuthService) GetUser(username string) (*models.User, error) {

	var user models.User

	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) UpdateUser(username string, req dtos.UpdateUserRequest) (*models.User, error) {

	var user models.User

	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user.Username = req.Username
	user.GivenName = req.GivenName
	user.Email = req.Email
	user.HashedPassword = hashedPassword

	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) DeleteUser(username string) error {

	result := s.DB.Where("username = ?", username).Delete(&models.User{})

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return result.Error
}

func (s *AuthService) hashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (s *AuthService) verifyPassword(hash string, password string) error {

	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}
func (s *AuthService) generateRefreshToken() (string, error) {
	token := uuid.New().String()
	return token, nil
}

func (s *AuthService) createRefreshToken(userID uuid.UUID) (*models.RefreshToken, error) {
	tokenString, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	refreshToken := models.RefreshToken{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.DB.Create(&refreshToken).Error; err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (s *AuthService) generateAccessToken(userID uuid.UUID) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *AuthService) Logout(refreshToken string) error {

	return s.DB.Where("token = ?", refreshToken).
		Delete(&models.RefreshToken{}).Error
}
