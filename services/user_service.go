package services

import (
	"goweb/api"
	"goweb/models"
	"goweb/requests"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{DB: db}
}

func (service *UserService) Register(e echo.Context, request *requests.RegisterRequest, tenant *models.Tenant) error {
	user := models.User{}

	service.GetUserByEmail(&user, request.Email)

	if user.ID != 0 {
		return api.WebResponse(e, http.StatusBadRequest, api.USER_EXISTS())
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return api.WebResponse(e, http.StatusInternalServerError, api.RESOURCE_CREATION_FAILED("Resource creation failed - error creating password"))
	}

	newUser := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(encryptedPassword),
		Tenants:  []*models.Tenant{tenant},
	}

	ok := service.DB.Create(&newUser).Error

	if ok != nil {
		return api.WebResponse(e, http.StatusInternalServerError, api.RESOURCE_CREATION_FAILED())
	}

	return api.WebResponse(e, http.StatusCreated, api.RESOURCE_CREATED("User created"))
}

func (service *UserService) GetUser(user *models.User, id int) {
	service.DB.First(user, id)
}

func (service *UserService) GetUserByUUID(user *models.User, uuid string) {
	service.DB.Where("uuid = ?", uuid).First(user)
}

func (service *UserService) GetUserByEmail(user *models.User, email string) {
	service.DB.Where("email = ?", email).First(user)
}

func (service *UserService) GetUsersByTenant(users *[]*models.User, tenant *models.Tenant) {
	service.DB.Joins("JOIN tenant_user ON tenant_user.user_id = users.id").
		Where("tenant_user.tenant_id = ?", tenant.ID).
		Preload("Tenants").
		Find(users)
}
