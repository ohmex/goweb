package services

import (
	"goweb/models"

	"gorm.io/gorm"
)

type DomainService struct {
	DB *gorm.DB
}

func NewDomainService(db *gorm.DB) *DomainService {
	return &DomainService{DB: db}
}

func (service *DomainService) GetDomainByUUID(domain *models.Domain, id string) error {
	return service.DB.
		Where("uuid = ?", id).
		First(domain).Error
}
