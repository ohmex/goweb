package services

import (
	"goweb/models"
	"goweb/requests"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (userService *UserService) Register(request *requests.RegisterRequest) error {
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

func (service *UserService) GetUser(user *models.User, id int) {
	service.DB.First(user, id)
}

func (service *UserService) GetUserByEmail(user *models.User, email string) {
	service.DB.Where("email = ?", email).Find(user)
}

func (service *UserService) GetUsersByTenant(users *[]*models.User, tenant *models.Tenant) {
	service.DB.Joins("JOIN tenant_user ON tenant_user.user_id = users.id").
		Where("tenant_user.tenant_id = ?", tenant.ID).
		Preload("Tenants").
		Find(users)
}
