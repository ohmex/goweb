package services

import (
	"errors"
	"goweb/api"
	"goweb/models"
	"goweb/requests"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RoleService struct {
	DB *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{DB: db}
}

func (service *RoleService) GetRolesInDomain(roles *[]*models.Role, domain *models.Domain) error {
	return service.DB.
		Where("domain_id in (-1, ?)", domain.ID).
		Find(roles).Error
}

func (service *RoleService) Create(e echo.Context, request *requests.RoleRequest, domain *models.Domain) error {
	var role models.Role

	service.DB.
		Where("name = ? AND domain_id = ?", request.Name, domain.ID).
		First(role)

	if role.ID != 0 {
		return api.WebResponse(e, http.StatusBadRequest, api.RESOURCE_EXISTS("Role already exists"))
	}

	role.Name = request.Name
	role.DomainID = domain.ID

	ok := service.DB.Create(&role).Error

	if ok != nil {
		return api.WebResponse(e, http.StatusInternalServerError, api.RESOURCE_CREATION_FAILED("Error creating role"))
	}

	return api.WebResponse(e, http.StatusCreated, api.RESOURCE_CREATED("Role created"))
}

func (service *RoleService) GetRoleByUuidInDomain(role *models.Role, uuid string, domain *models.Domain) error {
	return service.DB.
		Where("domain_id = ? AND uuid = ?", domain.ID, uuid).
		First(role).Error
}

func (service *RoleService) UpdateRole(role *models.Role) error {
	return service.DB.
		Save(role).Error
}

func (service *RoleService) DeleteRoleByUuidInDomain(role *models.Role, uuid string, domain *models.Domain) error {
	records := service.DB.
		Where("domain_id = ? AND uuid = ?", domain.ID, uuid).
		Delete(role)

	if records.RowsAffected > 0 {
		return nil
	}

	return errors.New("Failed")
}
