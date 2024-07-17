package services

import (
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

func (service *RoleService) GetRolesInDomain(roles *[]*models.Role, domain *models.Domain) {
	service.DB.
		Where("domain_id in (-1, ?)", domain.ID).
		Find(roles)
}

func (service *RoleService) Create(e echo.Context, request *requests.RoleRequest, domain *models.Domain) error {
	role := models.Role{}
	role.Name = request.Name

	service.GetRoleByDomain(&role, int(domain.ID))

	if role.ID != 0 {
		return api.WebResponse(e, http.StatusBadRequest, api.RESOURCE_EXISTS())
	}

	return nil
}

func (service *RoleService) GetRoleByDomain(role *models.Role, id int) {
	service.DB.Where("id = ?", id).First(role)
}
