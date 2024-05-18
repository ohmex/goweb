package services

import (
	"goweb/models"
	"goweb/requests"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func NewUserService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (userService *Service) Register(request *requests.RegisterRequest) error {
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(encryptedPassword),
	}

	return userService.DB.Create(&user).Error
}
