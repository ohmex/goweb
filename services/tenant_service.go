package services

import (
	"goweb/models"

	"gorm.io/gorm"
)

type TenantService struct {
	DB *gorm.DB
}

func NewTenantService(db *gorm.DB) *TenantService {
	return &TenantService{DB: db}
}

func (service *TenantService) GetTenantByUUID(tenant *models.Tenant, id string) {
	service.DB.Where("uuid = ?", id).Find(tenant)
}
