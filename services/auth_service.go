package services

import (
	"errors"
	"golang-api/models"
	"golang-api/utils"

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

func (as *AuthService) Register(user *models.User) (string, error) {
	if err := user.HashPassword(user.Password); err != nil {
		return "", errors.New("error hashing password")
	}

	if err := as.DB.Create(user).Error; err != nil {
		return "", errors.New("error creating user")
	}

	token, err := utils.GenerateToken(user.ID)

	if err != nil {
		return "", errors.New("error generateing token")
	}

	return token, nil
}

func (as *AuthService) Login(loginReq *models.LoginRequest) (string, error) {
	var user models.User

	if err := as.DB.Where("email = ?", loginReq.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid email or password")
		}
		return "", nil
	}

	if err := user.CheckPassword(loginReq.Password); err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID)

	if err != nil {
		return "", errors.New("error generating token")
	}

	return token, nil
}
