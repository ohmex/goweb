package services

import (
	"errors"
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

func (service *UserService) Register(e echo.Context, request *requests.RegisterRequest, domain *models.Domain) error {
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
		Domains:  []*models.Domain{domain},
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

func (service *UserService) GetUsersInDomain(users *[]*models.User, domain *models.Domain) {
	service.DB.
		Joins("JOIN domain_user ON domain_user.user_id = users.id").
		Where("domain_user.domain_id = ?", domain.ID).
		Find(users)
}

func (service *UserService) GetUserByUuidInDomain(user *models.User, uuid string, domain *models.Domain) {
	service.DB.
		Joins("JOIN domain_user ON domain_user.user_id = users.id").
		Where("domain_user.domain_id = ?", domain.ID).
		Where("uuid = ?", uuid).
		First(user)
}

func (service *UserService) DeleteUserByUuidInDomain(user *models.User, uuid string, domain *models.Domain) error {
	// delete the user only if User.uuid == uuid & User.domains contains selected domain
	// get the domains of the user that we want to delete
	// check if the seleted domain is contained in domains above
	service.DB.
		Joins("JOIN domain_user ON domain_user.user_id = users.id").
		Where("domain_user.domain_id = ?", domain.ID).
		Where("uuid = ?", uuid).
		Preload("Domains").
		First(user)
	if user.ID != 0 {
		for _, d := range user.Domains {
			if d.UUID.String() == domain.UUID.String() {
				service.DB.Delete(user)
				return nil
			}
		}
	}
	return errors.New("user not found")
}
